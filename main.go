
package main

import (
    "fmt"
    "strings"
    "flag"
    "os"    
    "os/signal"
    "runtime"
    log "./log"
    facade "./facade"
)


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
    CONF    Command = "conf"    
    INFO    Command = "info"
    HELP    Command = "help"
    TEST    Command = "test"
)
var cmds = []Command{CONF,PIPE,TEST}



var (
    textPort       uint     = 0xfcd
    confPort       uint     = 0xfcc
    connectHost    string   = "fcd.hq.feedface.com"
    connectTimeout float64  = 5.0
    listenHost     string   = "0.0.0.0"
    daemonize      bool     = false
)




func main() {
    quiet, verbose, debug := false, false, false
    directory := "~/.fcd/"
    
    
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, os.Interrupt)
    go func() {
        sig := <-signals
        log.Notice("%s",sig)
        os.Exit(0)
    }()


    log.SetVerbosity(log.NOTICE)
    
    flag.Usage = ShowHelpGeneral

    flags := make(map[Command] *flag.FlagSet)

    if RENDERER_AVAILABLE {
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

    if RENDERER_AVAILABLE {
        flags[RECV].UintVar(&confPort, "confport", confPort, "listen on `port` for config" )
        flags[RECV].UintVar(&textPort, "textport", textPort, "listen on `port` for text" )
        flags[RECV].StringVar(&listenHost, "host", listenHost, "listen on `host`" )
        flags[RECV].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
    }
    
    if RENDERER_AVAILABLE {
        flag.CommandLine.StringVar(&directory,  "D", directory,   "shader/font/texture")
    }    


    flag.CommandLine.BoolVar(&verbose,"v", verbose, "show info messages")
    flag.CommandLine.BoolVar(&debug,  "d", debug,   "show debug messages")
    flag.CommandLine.BoolVar(&quiet,  "q", quiet,   "show warnings only")
        
    
    
    flag.Parse()
    if flag.NArg() < 1 { 
        ShowHelpGeneral(); 
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
    
    
    
    cmd := Command(flag.Args()[0])
    switch ( cmd ) {

        case READ:
            if !RENDERER_AVAILABLE {
                ShowHelpGeneral()
                os.Exit(-2)    
            }
            flags[READ].Usage = func() { ShowHelpCommand(READ,flags) }
            flags[READ].Parse( flag.Args()[1:] )
            renderer = NewRenderer(directory)
            scanner = NewScanner()

        case RECV:
            if !RENDERER_AVAILABLE {
                ShowHelpGeneral()
                os.Exit(-2)    
            }
            flags[RECV].Usage = func() { ShowHelpCommand(RECV,flags) }
            flags[RECV].Parse( flag.Args()[1:] )
            server = NewServer(listenHost,confPort,textPort)
            renderer = NewRenderer(directory)

        case PIPE:
            flags[PIPE].Usage = func() { ShowHelpCommand(PIPE,flags) }
            flags[PIPE].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
            
        case CONF:
            flags[CONF].Usage = func() { ShowHelpCommand(CONF,flags) }
            flags[CONF].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
        case INFO:
            ShowVersion()
            os.Exit(-2)

        case HELP:
            ShowHelpGeneral()
            os.Exit(-2)
            
        
        case TEST:
            flags[TEST].Usage = func() {ShowHelpCommand(TEST,flags) }
            flags[TEST].Parse( flag.Args()[1:] )
            tester = NewTester()

        default:
            ShowHelpGeneral()
            os.Exit(-2)
    }
    
    
    
    var config *facade.Config = facade.NewConfig( facade.DEFAULT_MODE )
    var modeflags *flag.FlagSet = config.FlagSet()
    args := flags[cmd].Args()
    
    if len(args) < 1 {
        if cmd == CONF { 
            ShowHelpCommand(CONF,flags)
            os.Exit(-2)
        } else {
            
                modeflags.Usage = func() { ShowModeHelp(facade.Mode(config.Mode),cmd,modeflags) }
                modeflags.Parse( args[0:] )
                
        }
         
            
    } else {
        mode := facade.Mode(args[0])
        switch facade.Mode(mode) {
            
            case facade.GRID:
                config = facade.NewConfig(mode)
                modeflags = config.FlagSet()
                modeflags.Usage = func() { ShowModeHelp(facade.Mode(mode),cmd,modeflags) }
                modeflags.Parse( args[1:] )

            case facade.LINES:
                config = facade.NewConfig(mode)
                modeflags = config.FlagSet()
                modeflags.Usage = func() { ShowModeHelp(facade.Mode(mode),cmd,modeflags) }
                modeflags.Parse( args[1:] )

            case facade.TEST:
                config = facade.NewConfig(mode)
                modeflags = config.FlagSet()
                modeflags.Usage = func() { ShowModeHelp(facade.Mode(mode),cmd,modeflags) }
                modeflags.Parse( args[1:] )


                        
            default:
                ShowHelpGeneral()
                os.Exit(-2)    
        }
    }
        
    
    
    
    
    

    switch ( cmd ) {

        case READ:
            log.Info(AUTHOR)
            if renderer == nil { log.PANIC("renderer not available") }
            if scanner == nil { log.PANIC("scanner not available") }
            rawTexts := make(chan facade.RawText)
            texts := make(chan string)
            go scanner.ScanText(rawTexts)
            go renderer.ProcessText(rawTexts,texts)
            runtime.LockOSThread()
            renderer.Init(config) 
            renderer.Render(nil, texts)

        case RECV:
            log.Info(AUTHOR)
            if server == nil { log.PANIC("server not available") }
            if renderer == nil { log.PANIC("renderer not available") }
            rawConfs := make(chan facade.Config)
            rawTexts := make(chan facade.RawText)
            confs := make(chan facade.Config)
            texts := make(chan string)
            go server.ListenConf(rawConfs)
            go server.ListenText(rawTexts)
            go renderer.ProcessConf(rawConfs,confs)
            go renderer.ProcessText(rawTexts,texts)
            runtime.LockOSThread()
            renderer.Init(config) 
            renderer.Render(confs, texts)
                    
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
            tester.Test(str)
            
        default:
            log.PANIC("inconsistent command")
    }
        
        
    
}


func ShowModeHelp(mode facade.Mode, cmd Command, flagset *flag.FlagSet) {
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



func ShowHelpCommand(cmd Command, flagSetMap map[Command]*flag.FlagSet) {
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

func ShowHelpGeneral() {
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
    for _,m := range facade.Modes {
        fmt.Fprintf(os.Stderr,"%s|",m)
    }
    fmt.Fprintf(os.Stderr,"\n")
    fmt.Fprintf(os.Stderr,"\nCommands:\n")
    if RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"  %6s    # %s\n",READ,"read stdin and display")
        fmt.Fprintf(os.Stderr,"  %6s    # %s\n",RECV,"receive and display ")
    }
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",PIPE,"pipe stdin to remote")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",CONF,"configure remote")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",INFO,"show version info")
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr,"  %6s    # %s\n",facade.GRID,"character grid")
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
    
    
}
    

func ShowVersion() {
    fmt.Printf(AUTHOR)
    fmt.Fprintf(os.Stderr,"\n%s version %s for %s, built %s\n",BUILD_NAME,BUILD_VERSION,BUILD_PLATFORM,BUILD_DATE)
}    
    






