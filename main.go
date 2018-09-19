
package main

import (
    "fmt"
    "strings"
    "flag"
    "os"    
    "os/signal"
    "runtime"
    log "./log"
    render "./render"
    conf "./conf"
)


const AUTHOR = `
   _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _        
  |_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|  
  |   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  . \_  \_/ |  |  
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
    CONF    Command = "conf"    
    INFO    Command = "info"
    HELP    Command = "help"
    TEST    Command = "test"
)
var cmds = []Command{CONF,PIPE,TEST}



var (
    textPort       uint     = 0xfcd
    confPort       uint     = 0xfcc
    connectHost    string   = "localhost"
    connectTimeout float64  = 5.0
    listenHost     string   = "0.0.0.0"
    daemonize      bool     = false
)




func main() {
    quiet, verbose, debug := false, false, false
    
    
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

    if render.RENDERER_AVAILABLE {
        cmds = append(cmds, READ)
        cmds = append(cmds, RECV)
    }
    
    for _,cmd := range cmds {
        flags[cmd] = flag.NewFlagSet(string(cmd), flag.ExitOnError)
    }

    for _,cmd := range []Command{PIPE} {
        flags[cmd].UintVar(&textPort, "textport", textPort, "connect to `port` for text" )
    }
    
    for _,cmd := range []Command{PIPE,CONF} {
        flags[cmd].UintVar(&confPort, "confport", confPort, "connect to `port` for config" )
        flags[cmd].StringVar(&connectHost, "host", connectHost, "connect to `host`" )
        flags[cmd].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout after `seconds`") 
    }

    if render.RENDERER_AVAILABLE {
        flags[RECV].UintVar(&confPort, "confport", confPort, "listen on `port` for config" )
        flags[RECV].UintVar(&textPort, "textport", textPort, "listen on `port` for text" )
        flags[RECV].StringVar(&listenHost, "host", listenHost, "listen on `host`" )
        flags[RECV].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
    }
    
    


    all := []*flag.FlagSet{flag.CommandLine}
//    for _,cmd := range cmds {
//        all = append(all,flags[cmd])
//    }
    for _,flagSet := range all {
        flagSet.BoolVar(&verbose,"v", verbose, "show info messages")
        flagSet.BoolVar(&debug,  "d", debug,   "show debug messages")
        flagSet.BoolVar(&debug,  "q", debug,   "show warning messages only")
    }
    
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
    var renderer *render.Renderer
    var tester *Tester
    
    
    
    cmd := Command(flag.Args()[0])
    switch ( cmd ) {

        case READ:
            if !render.RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[READ].Usage = func() { ShowCommandHelp(READ,flags) }
            flags[READ].Parse( flag.Args()[1:] )
            renderer = render.NewRenderer()
            scanner = NewScanner()

        case RECV:
            if !render.RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[RECV].Usage = func() { ShowCommandHelp(RECV,flags) }
            flags[RECV].Parse( flag.Args()[1:] )
            server = NewServer(listenHost,confPort,textPort)
            renderer = render.NewRenderer()

        case PIPE:
            flags[PIPE].Usage = func() { ShowCommandHelp(PIPE,flags) }
            flags[PIPE].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
            
        case CONF:
            flags[CONF].Usage = func() { ShowCommandHelp(CONF,flags) }
            flags[CONF].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
        case INFO:
            ShowVersion()
            os.Exit(-2)

        case HELP:
            ShowHelp()
            os.Exit(-2)
            
        
        case TEST:
            flags[TEST].Usage = func() {ShowCommandHelp(TEST,flags) }
            flags[TEST].Parse( flag.Args()[1:] )
            tester = NewTester()

        default:
            ShowHelp()
            os.Exit(-2)
    }
    
    
    
    var config *conf.Config = conf.NewConfig(conf.DEFAULT_MODE)
    var modeflags *flag.FlagSet = config.FlagSet()
    args := flags[cmd].Args()
    
    if len(args) < 1 {
        if cmd == CONF { 
            ShowCommandHelp(CONF,flags)
            os.Exit(-2)
        } else {
            
                modeflags.Usage = func() { ShowModeHelp(config.Mode,cmd,modeflags) }
                modeflags.Parse( args[0:] )
                
        }
         
            
    } else {
        mode := conf.Mode(args[0])
        switch (mode) {
            
            case conf.GRID:
                config = conf.NewConfig(conf.GRID)
                modeflags = config.FlagSet()
                modeflags.Usage = func() { ShowModeHelp(conf.GRID,cmd,modeflags) }
                modeflags.Parse( args[1:] )


                        
            default:
                ShowHelp()
                os.Exit(-2)    
        }
    }
        
    
    
    
    
    

    switch ( cmd ) {

        case READ:
            log.Info(AUTHOR)
            if renderer == nil { log.PANIC("renderer not available") }
            if scanner == nil { log.PANIC("scanner not available") }
            texts := make(chan conf.Text)
            go scanner.ScanText(texts)
            go renderer.ReadText(texts)
            runtime.LockOSThread()
            renderer.Init() 
            renderer.Configure(config)
            renderer.Render()

        case RECV:
            log.Info(AUTHOR)
            if server == nil { log.PANIC("server not available") }
            if renderer == nil { log.PANIC("renderer not available") }
            texts := make(chan conf.Text)
            confs := make(chan conf.Config)
            go server.ListenText(texts)
            go server.ListenConf(confs)
            go renderer.ReadText(texts)
            go renderer.ReadConf(confs)
            runtime.LockOSThread()
            renderer.Init() 
            renderer.Configure(config)
            renderer.Render()
                    
        case PIPE:
            if client == nil { log.PANIC("client not available") }
            if config != nil {
                client.SendConf(config)
            }
            client.ScanAndSendText()
            
        case CONF:
            if client == nil { log.PANIC("client not available") }
            if config == nil { log.PANIC("config not available") }
            client.SendConf(config)


        case TEST:
            if tester == nil { log.PANIC("tester not available") }
            str := "FEEDFACE.COM"
            if modeflags.NArg() > 0 {
                str = strings.Join(modeflags.Args()," ")
            }
            tester.Configure(config)
            tester.testCharMap()
            tester.testTextTex(str)
            
        default:
            log.PANIC("inconsistent command")
    }
        
        
    
}


func ShowModeHelp(mode conf.Mode, cmd Command, flagset *flag.FlagSet) {
    switches := ""
    flags := ""
    flagset.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) == 1 { switches += " [ -"+f.Name+name+" ]" }
        if len(f.Name) >  1 { flags += " [ -"+f.Name+name+" ]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")    
    fmt.Fprintf(os.Stderr,"\n")
    fmt.Fprintf(os.Stderr,"  %s %s %s%s%s\n",BUILD_NAME,cmd,mode,switches,flags)
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagset.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}



func ShowCommandHelp(cmd Command, flagSetMap map[Command]*flag.FlagSet) {
    switches := ""
    flags := ""
    flagSetMap[cmd].VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) == 1 { switches += " [ -"+f.Name+name+" ]" }
        if len(f.Name) >  1 { flags += " [ -"+f.Name+name+" ]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s%s%s\n",BUILD_NAME,cmd,switches,flags)
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagSetMap[cmd].PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}


func ShowHelp() {
    flags := ""
    flag.CommandLine.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) >=  1 { flags +=    " [-"+f.Name+name+"]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s   ",BUILD_NAME,flags)
    for _,cmd := range cmds {
        fmt.Fprintf(os.Stderr,"%s|",cmd)
    }
    fmt.Fprintf(os.Stderr,"    ")
    for _,m := range conf.Modes {
        fmt.Fprintf(os.Stderr,"%s|",m)
    }
    fmt.Fprintf(os.Stderr,"\n")
    fmt.Fprintf(os.Stderr,"\nCommands:\n")
    if render.RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"  %6s    # %s\n",READ,"pipe stdin to display")
        fmt.Fprintf(os.Stderr,"  %6s    # %s\n",RECV,"receive text and display")
    }
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",PIPE,"pipe stdin to remote facade")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",CONF,"control remote facade")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",INFO,"show facade info")
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",conf.GRID,"a grid")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",conf.CLOUD,"a cloud")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",conf.SCROLL,"a scroller")
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}
    

func ShowVersion() {
    fmt.Printf(AUTHOR)
    fmt.Fprintf(os.Stderr,"\n%s version %s for %s, built %s\n",BUILD_NAME,BUILD_VERSION,BUILD_PLATFORM,BUILD_DATE)
}    
    






