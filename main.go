
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
    "time"
    log "./log"
    facade "./facade"
    gfx "./gfx"
//    proto "./facade/proto"
)


const DEBUG_CLOCK    = false
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
    confPort       uint     = 0xfcc
    connectHost    string   = "fcd.hq.feedface.com"
    connectTimeout float64  = 5.0
    readTimeout    float64  = 0.0
    listenHost     string   = "0.0.0.0"
//    daemonize      bool     = false
)




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
        flags[cmd].UintVar(&confPort, "cp", confPort, "connect to `port` for config" )
        flags[cmd].StringVar(&connectHost, "h", connectHost, "connect to `host`" )
        flags[cmd].Float64Var(&connectTimeout, "t", connectTimeout, "timeout connect after `seconds`") 
    }

    if flags[RECV] != nil {
        flags[RECV].UintVar(&confPort, "cp", confPort, "listen on `port` for config" )
        flags[RECV].UintVar(&textPort, "tp", textPort, "listen on `port` for text" )
        flags[RECV].StringVar(&listenHost, "h", listenHost, "listen on `host`" )
//        flags[RECV].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
        flags[RECV].Float64Var(&readTimeout, "t", readTimeout, "timeout read after `seconds`") 
    }

    if flags[TEST] != nil {
        flags[TEST].UintVar(&confPort, "cp", confPort, "listen on `port` for config" )
        flags[TEST].UintVar(&textPort, "tp", textPort, "listen on `port` for text" )
        flags[TEST].StringVar(&listenHost, "h", listenHost, "listen on `host`" )
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
        case READ:
            renderer = NewRenderer(directory)
            scanner = NewScanner()
            
        case RECV:
            server = NewServer(listenHost,confPort,textPort,readTimeout)
            renderer = NewRenderer(directory)

        case PIPE:
            client = NewClient(connectHost,confPort,textPort,connectTimeout)

        case CONF:
            client = NewClient(connectHost,confPort,textPort,connectTimeout)

        case EXEC:
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            executor = NewExecutor(client)

        case TEST:
            scanner = NewScanner()
            server = NewServer(listenHost,confPort,textPort,readTimeout)
            tester = NewTester(directory)

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


    args := flags[cmd].Args()
    var mode string
    var modeFlags *flag.FlagSet
    var config *facade.Config = &facade.Config{}
    
    
    
//    // parse mode, if given
    if len(args) > 0 {
        mode = strings.ToLower( args[0] )
        
        switch strings.ToUpper(mode) {

            case facade.Mode_GRID.String():
                config.CheckMode = true
                config.Mode = facade.Mode_GRID
                config.Grid = &facade.GridConfig{}
                modeFlags = flag.NewFlagSet(mode, flag.ExitOnError)
                modeFlags.Usage = func() { ShowHelpMode(mode,cmd,modeFlags) }
        }

        args = args[1:]

    }
    
    

    if config.Grid != nil {
        config.Grid.AddFlags( modeFlags )
        modeFlags.Parse( args )
        config.Grid.VisitFlags( modeFlags )
        log.Debug("grid %s parsed",config.Desc())
    }

        
        
    
//    config = ParseGridFlags(modeFlags, vars)
    

//    
//        
//    state = facade.NewState(mode)
//    var modeFlags = flag.NewFlagSet(string(mode), flag.ExitOnError)    
//    
//    state.AddFlags( modeFlags )
//    modeFlags.Usage = func() { ShowHelpMode(mode,cmd,modeFlags) }
//    modeFlags.Parse( args[0:] )
//
//    
//    
//    config := state.CheckFlags(modeFlags)


//    if cmd == EXEC {
//        var ok bool
//        var grid facade.GridConfig
//        if grid,ok = config.Grid(); !ok {                           //REM, make grid!!
//            grid = facade.GridConfig{}
//            config.SetGrid(grid)
////            log.PANIC("exec without grid")
//        } 
//        var cols,rows = uint(40), uint(12)
//        
//        if c,ok := grid.Width(); ok  { cols = c }
//        if r,ok := grid.Height(); ok { rows = r }
//        grid.SetWidth(cols)
//        grid.SetHeight(rows)
//        executor.SetSize(cols,rows)
//               
//        args := modeFlags.Args()
//        
//        if len(args) <= 0 {
//            ShowHelpMode(facade.GRID,EXEC,modeFlags)
//            os.Exit(-2)             
//        }
//        executor.SetPath(args[0])
//        executor.SetArgs(args[1:])
//        
//    }
    
    
//    var config *facade.Config = &facade.Config{}
    
    var err error
    switch ( cmd ) {

        case READ:
            log.Info(AUTHOR)
            if renderer == nil { log.PANIC("renderer not available") }
            if scanner == nil { log.PANIC("scanner not available") }
            bufChan := make(chan facade.BufferItem)
            go scanner.ScanText(bufChan)
            runtime.LockOSThread()
            renderer.Init(config) 
            
            //start processing only after init!
            go renderer.ProcessBufferItems(bufChan)
            
            err = renderer.Render(nil)
            

        case RECV:
            log.Info(AUTHOR)
            if server == nil { log.PANIC("server not available") }
            if renderer == nil { log.PANIC("renderer not available") }
            rawConfs := make(chan facade.Config)
            confs := make(chan facade.Config)
            bufChan := make(chan facade.BufferItem)

            go server.ListenConf(rawConfs)
            go server.ListenText(bufChan)


            runtime.LockOSThread()
            renderer.Init(config) 

            //start processing only after init!
            go renderer.ProcessRawConfs(rawConfs,confs)
            go renderer.ProcessBufferItems(bufChan)
            
            err = renderer.Render(confs)
                    
        case PIPE:
            if client == nil { log.PANIC("client not available") }
            if config != nil {
                client.SendConf(config)
            }
            err = client.OpenText()
            client.ScanAndSendText()
            client.CloseText()
            
        case CONF:
            if client == nil { log.PANIC("client not available") }
            if config == nil { log.PANIC("config not available") }
            err = client.SendConf(config)

        case EXEC:
            if client == nil { log.PANIC("client not available") }
            if executor == nil { log.PANIC("executor not available") }
            for config != nil { 
                err = client.SendConf(config)
                if err == nil {
                    log.Debug("sent config %s",config.Desc())
                    break
                }
                time.Sleep( time.Duration( 200 * time.Millisecond ) )
            }
            for {
                err = client.OpenText()
                if err == nil {
                    log.Debug("connected text.")
                    break
                }
                time.Sleep( time.Duration( 200 * time.Millisecond ) )
            }
            err = executor.Execute()
            

        case TEST:
            log.Info(AUTHOR)
            if scanner == nil { log.PANIC("scanner not available") }
            if server == nil { log.PANIC("server not available") }
            if tester == nil { log.PANIC("tester not available") }
            rawConfs := make(chan facade.Config)
            confs := make(chan facade.Config)
            bufChan := make(chan facade.BufferItem)

            

            go server.ListenConf(rawConfs)
            go server.ListenText(bufChan)
            go scanner.ScanText(bufChan)

            runtime.LockOSThread()
            tester.Init(config) 
            tester.Configure(config)
            
            //start processing only after init!
            go tester.ProcessRawConfs(rawConfs,confs)
            go tester.ProcessBufferItems(bufChan)

            
            err = tester.Test(confs)
            

        default:
            log.PANIC("inconsistent command")
    }
        

    if err != nil {
        log.Error("could not %s: %s",cmd,err)
        os.Exit(-1)
    }
    
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
    






