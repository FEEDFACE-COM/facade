
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


const DEBUG_CLOCK    = true
const DEBUG_CONFIG   = true
const DEBUG_FONT     = true
const DEBUG_MODE     = true
const DEBUG_GRID     = false
const DEBUG_DIAG     = false
const DEBUG_MEMORY   = false
const DEBUG_BUFFER   = false


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
    
    for _,cmd := range []Command{PIPE,CONF,EXEC,INFO} {
        flags[cmd].UintVar(&port, "p", port, "connect to `port`" )
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




    var config *facade.Config = &facade.Config{}
    config.Font = &facade.FontConfig{}
    config.Camera = &facade.CameraConfig{}
    config.Mask = &facade.MaskConfig{}

    var args []string
    var mode string
    var modeFlags *flag.FlagSet
    

    

    switch (cmd) {

        case READ,RECV,PIPE,CONF,EXEC,TEST:
            // parse mode, if given
            args = flags[cmd].Args()
        
            if cmd != INFO && len(args) > 0 {
                mode = strings.ToLower( args[0] )
                
                switch strings.ToUpper(mode) {
        
                    case facade.Mode_GRID.String():
                        config.SetMode = true
                        config.Mode = facade.Mode_GRID
                        config.Grid = &facade.GridConfig{}
                        
                        
                    case facade.Mode_DRAFT.String():
                        config.SetMode = true
                        config.Mode = facade.Mode_GRID
                }
                args = args[1:]
                
                
                
                modeFlags = flag.NewFlagSet(mode, flag.ExitOnError)
                modeFlags.Usage = func() { ShowHelpMode(mode,cmd,modeFlags) }
                
            
                config.AddFlags( modeFlags )
                modeFlags.Parse( args )
                config.VisitFlags( modeFlags )
                
            }




            
        case INFO:
            // no args, print local info
            if len(flag.Args()) <= 1 { 
                ShowVersion()
                if log.InfoLogging() {
                    ShowAssets()
                }
                fmt.Fprintf(os.Stderr,"\n\n")
                os.Exit(0)
            } else { 
                
                // query remote host
                flags[INFO].Usage = func() { ShowHelpCommand(INFO,flags) }
                flags[INFO].Parse( flag.Args()[1:] )

            }
            
        
        

        case HELP:
            ShowHelp()
            os.Exit(-2)
            
        default:
            ShowHelp()
            os.Exit(-2)
        
    }



    

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
    confs := make(chan facade.Config)
    texts := make(chan facade.TextSeq)
    quers := make(chan (chan string))
    switch ( cmd ) {
        

        case READ:
            log.Info(AUTHOR)
            scanner = NewScanner()
            renderer = NewRenderer(directory)
            go scanner.ScanText(texts)
            runtime.LockOSThread()
            renderer.Init(config)  
            go renderer.ProcessTextSeqs(texts)
            err = renderer.Render(nil)
            

        case RECV:
            log.Info(AUTHOR)
            server = NewServer(host,port,textPort,readTimeout)
            renderer = NewRenderer(directory)
            go server.Listen(confs,texts,quers)
            go server.ListenText(texts)
            runtime.LockOSThread()
            renderer.Init(config) 
            go renderer.ProcessTextSeqs(texts)
            go renderer.ProcessQueries(quers)
            err = renderer.Render(confs)


                    
        case PIPE:
            client = NewClient(host,port,connectTimeout)
            if err=client.Dial(); err!=nil { log.Error("fail to dial: %s",err) }
            defer client.Close()
            if config != nil {
                log.Debug("configure %s",config.Desc())
                if client.SendConf(config); err!=nil { log.Error("fail to send conf: %s",err) }
            }
            if err=client.OpenTextStream(); err!=nil { log.Error("fail to open stream: %s",err) }
//            defer client.CloseTextStream()
            if err=client.ScanAndSendText(); err!=nil { log.Error("fail to scan and send: %s",err) }
            

        case CONF:
            if config == nil {
                ShowHelpMode(mode,cmd,modeFlags)
                os.Exit(-1)
            }
            log.Debug("configure %s",config.Desc())
            client = NewClient(host,port,connectTimeout)
            if err=client.Dial(); err!=nil { log.Error("fail to dial: %s",err) }
            defer client.Close()
            if err=client.SendConf(config); err!=nil { log.Error("fail to conf: %s",err) }


        case INFO:
            client = NewClient(host,port,connectTimeout)
            if err=client.Dial(); err!=nil { log.Error("fail to dial: %s",err) }
            defer client.Close()
            info, err := client.QueryInfo()
            if err != nil {
                log.Error("fail to query: %s",err)
            } else {
                log.Notice("%s",info)
            }




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
                log.Debug("configure %s",config.Desc())
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

            go server.Listen(confs,texts,quers)
            go server.ListenText(texts)
            go scanner.ScanText(texts)

            runtime.LockOSThread()
            tester.Init(config) 
            tester.Configure(config)
            
            //start processing only after init!
            go tester.ProcessTextSeqs(texts)
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
    mode = strings.ToLower(mode)
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
    if cmd == INFO {
        fmt.Fprintf(os.Stderr,"  %s %s [%s] %s\n",BUILD_NAME,cmd,switches,flags)
    } else {
        fmt.Fprintf(os.Stderr,"  %s %s [%s] %s\n",BUILD_NAME,cmd,switches,flags)
        ShowModes()
    }

    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagSetMap[cmd].PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}

func ShowCommands() {
    fmt.Fprintf(os.Stderr,"\nCommands:\n")
    if RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"%6s     %s\n",READ,"read text from stdin and render")
        fmt.Fprintf(os.Stderr,"%6s     %s\n",RECV,"receive text from remote and render ")
    }
    fmt.Fprintf(os.Stderr,"%6s     %s\n",PIPE,"read text from stdin and send to remote server")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",CONF,"change configuration of remote server")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",EXEC,"execute command and send stdout,stderr to server")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",INFO,"show available shaders and fonts of remote server ")
}


func ShowModes() {
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",strings.ToLower(facade.Mode_GRID.String()),"character grid")
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

func ShowAssets() { fmt.Fprintf(os.Stderr,InfoAssets()) }
func InfoAssets() string {
    ret := ""
    shaders := gfx.ListShaderNames()
    fonts := gfx.ListFontNames()
//    ret += "\nShaders:\n"
//    for _,s := range shaders {
//        ret += fmt.Sprintf("  %s\n",s)
//    }

    for _,prefix := range []string{"grid/",} {
        for _,suffix := range []string{".vert",".frag"} {
            ret += "\n" + strings.TrimSuffix(prefix,"/") + " "
            ret += "-" + strings.TrimPrefix(suffix,".") + "=  "
            for _,shader := range shaders {
                if strings.HasPrefix(shader,prefix) && strings.HasSuffix(shader,suffix) {
                    ret += strings.TrimSuffix(strings.TrimPrefix(shader,prefix),suffix)
                    ret += "  "
                }
                    
            }
        }
    }
    
    
    ret += "\n-mask=  "
    for _,shader := range shaders {
        if strings.HasPrefix(shader,"mask/") && strings.HasSuffix(shader,"frag") {
            ret += strings.TrimSuffix(strings.TrimPrefix(shader,"mask/"),".frag")
            ret += "  "
        }
    }


    ret += "\n-font=  "
    for _,font := range fonts {
        ret += font
        ret += "  "
    }
    ret += "\n"
    return ret
}

func ShowVersion() { fmt.Fprintf(os.Stderr,InfoVersion()) }
func InfoVersion() string {
    ret := ""
    ret += AUTHOR
    ret += fmt.Sprintf("\n%s version %s for %s, built %s\n",BUILD_NAME,BUILD_VERSION,BUILD_PLATFORM,BUILD_DATE)
    return ret
}    
    






