
package main

import (
    "fmt"
    "strings"
    "flag"
//    "bufio"
//    "io"
    "os"    
    "os/signal"
    "runtime"
    log "./log"
    facade "./facade"
    gfx "./gfx"
//    proto "./facade/proto"
)


const DEBUG_CLOCK    = false
const DEBUG_CONFIG   = true
const DEBUG_MODE     = true
const DEBUG_GRID     = false
const DEBUG_DIAG     = false
const DEBUG_MEMORY   = false
const DEBUG_MESSAGES = true
const DEBUG_BUFFER   = true


const FRAME_RATE = 60.0
const BUFFER_SIZE = 80





const AUTHOR = `
   _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _        
  |_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|  
  |   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |  
`


var BUILD_NAME string     = "fcd"
var BUILD_VERSION string  = "0.0.0"
var BUILD_PLATFORM string = "os-arch"
var BUILD_DATE string     = "1970-01-01"

type Command string
const (
    READ    Command = "read"
    RECV    Command = "recv"
    PIPE    Command = "pipe"
    EXEC    Command = "exec"
    CONF    Command = "conf"    
    INFO    Command = "info"
    HELP    Command = "help"
    TEST    Command = "test"
)



var (
    textPort       uint     = 0xfcd
    port           uint     = 0xfcc
    host         string     = ""
    connectTimeout float64  = 5.0
    readTimeout    float64  = 0.0
)


const DEFAULT_LISTEN_HOST = "0.0.0.0"
const DEFAULT_CONNECT_HOST = "localhost"



func main() {
    quiet, verbose, debug := false, false, false
    directory := facade.DEFAULT_DIRECTORY
    
    
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    go func() {
        sig := <-signals
        log.Notice("%s",sig)
        os.Exit(0)
    }()


    log.SetVerbosity(log.NOTICE)
    
    flag.Usage = ShowHelp

    flags := make(map[Command] *flag.FlagSet)

    var commands = []Command{CONF,PIPE,EXEC,INFO,TEST}
    if RENDERER_AVAILABLE {
        commands = append(commands, READ)
        commands = append(commands, RECV)
    }
    
    for _,cmd := range commands {
        flags[cmd] = flag.NewFlagSet(string(cmd), flag.ExitOnError)
    }

    for _,cmd := range []Command{PIPE,EXEC} {
        flags[cmd].UintVar(&textPort, "tp", textPort, "connect to `port` for text" )
    }
    
    for _,cmd := range []Command{PIPE,CONF,EXEC} {
        flags[cmd].UintVar(&port, "p", port, "connect to `port` for config" )
        flags[cmd].StringVar(&host, "h", DEFAULT_CONNECT_HOST, "connect to `host`" )
        flags[cmd].Float64Var(&connectTimeout, "t", connectTimeout, "timeout connect after `seconds`") 
    }

    if flags[RECV] != nil {
        flags[RECV].UintVar(&port, "cp", port, "listen on `port` for config" )
        flags[RECV].UintVar(&textPort, "tp", textPort, "listen on `port` for text" )
        flags[RECV].StringVar(&host, "h", DEFAULT_LISTEN_HOST, "listen on `host`" )
        flags[RECV].Float64Var(&readTimeout, "t", readTimeout, "timeout read after `seconds`") 
    }

    if flags[TEST] != nil {
        flags[TEST].UintVar(&port, "cp", port, "listen on `port` for config" )
        flags[TEST].UintVar(&textPort, "tp", textPort, "listen on `port` for text" )
        flags[TEST].StringVar(&host, "h", DEFAULT_LISTEN_HOST, "listen on `host`" )
        flags[TEST].Float64Var(&readTimeout, "t", readTimeout, "timeout read after `seconds`") 
    }
    
    {
        flag.CommandLine.StringVar(&directory,  "D", directory,   "directory")
    }    


    flag.CommandLine.BoolVar(&verbose,"v", verbose, "show info messages")
    flag.CommandLine.BoolVar(&debug,  "d", debug,   "show debug messages")
    flag.CommandLine.BoolVar(&quiet,  "q", quiet,   "show warnings only")
        
    
    
    flag.Parse()
    if flag.NArg() < 1 { 
        ShowHelp(); 
        os.Exit(-2) 
    }
    if debug { 
        log.SetVerbosity(log.DEBUG) 
    } else if verbose { 
        log.SetVerbosity(log.INFO)
    } else if quiet { 
        log.SetVerbosity(log.WARNING) 
    }
    
    
    var client *Client
    var server *Server
    var scanner *Scanner
    var renderer *Renderer
    var tester *Tester
    var executor *Executor
    var path string
    
    
    
    cmd := Command(flag.Args()[0])

    switch (cmd) {
        case READ, RECV:
            if !RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            fallthrough
        case PIPE, CONF, EXEC, TEST:
            flags[cmd].Usage = func() { ShowHelpCommand(cmd,flags) }
            flags[cmd].Parse( flag.Args()[1:] )
    }
    

    switch (cmd) {

        case READ,RECV,PIPE,CONF,EXEC,TEST:
            break
            
        case INFO:
            ShowVersion()
            ShowAssets()
            os.Exit(-2)

        case HELP:
            ShowHelp()
            os.Exit(-2)
            
        default:
            ShowHelp()
            os.Exit(-2)
        
    }

    var args []string
    var mode string
    var modeFlags *flag.FlagSet

    var config *facade.Config = &facade.Config{}
    config.Font = &facade.FontConfig{}
    config.Camera = &facade.CameraConfig{}
    config.Mask = &facade.MaskConfig{}
    

    // parse mode, if given
    args = flags[cmd].Args()
    if len(args) > 0 {
        mode = strings.ToLower( args[0] )
        
        switch strings.ToUpper(mode) {

            case facade.Mode_GRID.String():
                config.SetMode = true
                config.Mode = facade.Mode_GRID
                config.Grid = &facade.GridConfig{}
                
                
            case facade.Mode_DRAFT.String():
                config.SetMode = true
        }
        args = args[1:]
    }


    modeFlags = flag.NewFlagSet(mode, flag.ExitOnError)
    modeFlags.Usage = func() { ShowHelpMode(mode,cmd,modeFlags) }
    
    
    config.AddFlags( modeFlags )
    modeFlags.Parse( args )
    config.VisitFlags( modeFlags )


    if cmd == EXEC {
        
        args = modeFlags.Args()
        
        if len(args) <= 0 { // no command given
            ShowHelpMode(facade.Mode_GRID.String(),EXEC,modeFlags)
            os.Exit(-2)             
        }
        
        path = args[0]
        args = args[1:]
        

    }        
        
        
    
    
    
    var err error
    switch ( cmd ) {

        case READ:
            log.Info(AUTHOR)
            scanner = NewScanner()
            renderer = NewRenderer(directory)
            texts := make(chan facade.BufferItem)
            go scanner.ScanText(texts)
            runtime.LockOSThread()
            renderer.Init(config)  
            go renderer.ProcessBufferItems(texts)
            err = renderer.Render(nil)
            

        case RECV:
            log.Info(AUTHOR)
            server = NewServer(host,port,textPort,readTimeout)
            renderer = NewRenderer(directory)
            confs := make(chan facade.Config)
            texts := make(chan facade.BufferItem)
            go server.Listen(confs,texts)
            go server.ListenText(texts)
            runtime.LockOSThread()
            renderer.Init(config) 
            go renderer.ProcessBufferItems(texts)
            err = renderer.Render(confs)


                    
        case PIPE:
            client = NewClient(host,port,connectTimeout)
            if err=client.Dial(); err!=nil { log.Error("fail to dial: %s",err) }
            defer client.Close()
            if config != nil {
                if client.SendConf(config); err!=nil { log.Error("fail to send conf: %s",err) }
            }
            if err=client.OpenTextStream(); err!=nil { log.Error("fail to open stream: %s",err) }
            defer client.CloseTextStream()
            if err=client.ScanAndSendText(); err!=nil { log.Error("fail to scan and send: %s",err) }
            
        case CONF:
            client = NewClient(host,port,connectTimeout)
            if err=client.Dial(); err!=nil { log.Error("fail to dial: %s",err) }
            defer client.Close()
            if err=client.SendConf(config); err!=nil { log.Error("fail to send conf: %s",err) }

        case EXEC:

            var cols,rows = uint64(40), uint64(12)
            if config.GetGrid() == nil {
                config.Grid = &facade.GridConfig{}
            }
            
            if config.GetGrid().GetSetWidth() {
                cols = config.GetGrid().GetWidth()
            }
            if config.GetGrid().GetSetWidth() {
                rows = config.GetGrid().GetHeight()
            }
            
            config.Grid.Width = cols
            config.Grid.SetWidth = true

            config.Grid.Height = rows
            config.Grid.SetHeight = true

            config.Grid.Terminal = true
            config.Grid.SetTerminal = true
            
            client = NewClient(host,port,connectTimeout)
            executor = NewExecutor( client, uint(cols), uint(rows), path, args )



            if err=client.Dial(); err!=nil { log.Error("fail to dial: %s",err) }
            defer client.Close()
            if config != nil {
                if client.SendConf(config); err!=nil { log.Error("fail to send conf: %s",err) }
            }
            if err=client.OpenTextStream(); err!=nil { log.Error("fail to open stream: %s",err) }
            defer client.CloseTextStream()




            err = executor.Execute()
            

        case TEST:
            log.Info(AUTHOR)
            scanner = NewScanner()
            server = NewServer(host,port,textPort,readTimeout)
            tester = NewTester(directory)
            confs := make(chan facade.Config)
            texts := make(chan facade.BufferItem)

            

            go server.Listen(confs,texts)
            go server.ListenText(texts)
            go scanner.ScanText(texts)

            runtime.LockOSThread()
            tester.Init(config) 
            tester.Configure(config)
            
            //start processing only after init!
//            go tester.ProcessRawConfs(rawConfs,confs)
            go tester.ProcessBufferItems(texts)
            err = tester.Test(confs)
            

        default:
            log.PANIC("unexpected command %s",cmd)
    }
        

    if err != nil {
        log.Error("could not %s: %s",cmd,err)
        os.Exit(-1)
    }
    
    
    os.Exit(0)
}


func ShowHelpMode(mode string, cmd Command, flagset *flag.FlagSet) {
    switches := "-"
    flags := ""
    flagset.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "=" }
        if len(f.Name) == 1 && name == "" { switches += f.Name }
        if len(f.Name) >  1 || name != "" { flags += " [-"+f.Name+name+"]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s %s [%s]%s\n",BUILD_NAME,cmd,mode,switches,flags)
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagset.VisitAll( func( f *flag.Flag) {
        tmp,_ := flag.UnquoteUsage(f)
        typ := ""
        if tmp != "" {
            typ = fmt.Sprintf("(%s)",tmp)
        }
        fmt.Fprintf(os.Stderr,"  -%-8s %-24s %-8s\n",f.Name,f.Usage,typ)
    })
    fmt.Fprintf(os.Stderr,"\n")
}



func ShowHelpCommand(cmd Command, flagSetMap map[Command]*flag.FlagSet) {
    modes := []string{}
    for _,m := range( []facade.Mode{ facade.Mode_GRID } ) { 
        modes = append(modes, string(m) )
    }
    switches := "-"
    flags := ""
    flagSetMap[cmd].VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) == 1 && name == "" { switches += f.Name }
        if len(f.Name) >  1 || name != "" { flags += " [-"+f.Name+name+"]" }
    })

    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s [%s] %s    %s\n",BUILD_NAME,cmd,switches,flags,strings.Join(modes," | "))
    ShowModes()

    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagSetMap[cmd].PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}

func ShowCommands() {
    fmt.Fprintf(os.Stderr,"\nCommands:\n")
    if RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"%6s     %s\n",READ,"read text from stdin and render")
        fmt.Fprintf(os.Stderr,"%6s     %s\n",RECV,"receive text from network and render ")
    }
    fmt.Fprintf(os.Stderr,"%6s     %s\n",PIPE,"read text from stdin and send to server")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",CONF,"change configuration of server")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",EXEC,"execute command and send stdout,stderr to server")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",INFO,"show available shaders and fonts ")
}


func ShowModes() {
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",facade.Mode_GRID,"character grid")
}

func ShowHelp() {
    flags := ""
    flag.CommandLine.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) >=  1 { flags +=    " [-"+f.Name+name+"]" }
    })
    cmds := []string{}
        if RENDERER_AVAILABLE {
        for _,c := range( []Command{READ,RECV}){
            cmds = append(cmds, string(c) )
        }
    }
    for _,c := range( []Command{PIPE,CONF,EXEC,INFO}){
        cmds = append(cmds, string(c) )
    }

    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s      %s\n",BUILD_NAME,flags,strings.Join(cmds," | "))
    ShowCommands()
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
    
    
}

func ShowAssets() {
    shaders := gfx.ListShaderNames()
    fonts := gfx.ListFontNames()
    fmt.Fprintf(os.Stderr,"\nShaders:\n")
    for _,s := range shaders {
        fmt.Fprintf(os.Stderr,"  %s\n",s)
    }
    fmt.Fprintf(os.Stderr,"\nFonts:\n")
    for _,f := range fonts {
        fmt.Fprintf(os.Stderr,"  %s\n",f)
    }
    fmt.Fprintf(os.Stderr,"\n")
}


func ShowVersion() {
    fmt.Printf(AUTHOR)
    fmt.Fprintf(os.Stderr,"\n%s version %s for %s, built %s\n",BUILD_NAME,BUILD_VERSION,BUILD_PLATFORM,BUILD_DATE)
}    
    






