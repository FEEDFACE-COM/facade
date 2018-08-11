
package main

import (
    "fmt"
    "strings"
    "flag"
    "os"    
    "os/signal"
//    "bufio"
    log "./log"
    render "./render"
)


const AUTHOR = 

// "FACADE by FEEDFACE.COM"

`   _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _        
  |_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|  
  |   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  o \_  \_/ |  |  
`

//`   _   _   _   _   _   _
//  |_  |_| /   |_| | \ |_
//  |   | | \_  | | |_/ |_  by FEEDFACE.COM
//
//`


var BUILD_NAME, BUILD_VERSION, BUILD_PLATFORM, BUILD_DATE string

type Mode string
const (
    Read    Mode = "read"
    Listen  Mode = "listen"
    Send    Mode = "send"
    Pipe    Mode = "pipe"
    Conf    Mode = "conf"    
    Version Mode = "info"
    Help    Mode = "help"
    Test    Mode = "test"
)
var modes = []Mode{Send,Conf}



var (
    textPort       uint     = 0xfcd
    confPort       uint     = 0xfcdc
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

    flags := make(map[Mode] *flag.FlagSet)

    if render.RENDERER_AVAILABLE {
        modes = append(modes, Read)
        modes = append(modes, Listen)
    }
    
    for _,mode := range modes {
        flags[mode] = flag.NewFlagSet(string(mode), flag.ExitOnError)
        flags[mode].Usage = func() { ShowModeHelp(mode,flags) }
    }


    flags[Send].UintVar(&textPort, "textport", textPort, "connect to `port` for text" )

    for _,mode := range []Mode{Send,Conf} {
        flags[mode].UintVar(&confPort, "confport", confPort, "connect to `port` for config" )
        flags[mode].StringVar(&connectHost, "host", connectHost, "connect to `host`" )
        flags[mode].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout after `seconds`") 
    }

    if render.RENDERER_AVAILABLE {
        flags[Listen].UintVar(&confPort, "confport", confPort, "listen on `port` for config" )
        flags[Listen].UintVar(&textPort, "textport", textPort, "listen on `port` for text" )
        flags[Listen].StringVar(&listenHost, "host", listenHost, "listen on `host`" )
        flags[Listen].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
    }

    all := []*flag.FlagSet{flag.CommandLine}
    for _,mode := range modes {
        all = append(all,flags[mode])
    }
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
    
    var client *Client
    var server *Server
    var renderer *render.Renderer
    var scanner *Scanner
    
    mode := Mode(flag.Args()[0])
    switch ( mode ) {

        case Read:
            if !render.RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[Read].Parse( flag.Args()[1:] )
            renderer = render.NewRenderer()
            scanner = NewScanner()

        case Listen:
            if !render.RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[Read].Parse( flag.Args()[1:] )
            server = NewServer(listenHost,confPort,textPort)
            renderer = render.NewRenderer()

            
        case Send:
            flags[Send].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            client.text = strings.Join(flags[Send].Args()[0:], "  ")
            
        case Conf:
            flags[Conf].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
        case Test:
            log.Fatal("TEST TEST TEST")

        case Version:
            ShowVersion()
            os.Exit(-2)

        case Help:
            ShowHelp()
            os.Exit(-2)

        default:
            ShowHelp()
            os.Exit(-2)
    }
    
    if debug { 
        log.SetVerbosity(log.DEBUG) 
    } else if verbose { 
        log.SetVerbosity(log.INFO)
    } else if quiet { 
        log.SetVerbosity(log.WARNING) 
    }
    
    
    log.Info(AUTHOR)


    switch ( mode ) {

        case Read:
            if renderer == nil { log.PANIC("renderer not available") }
            if scanner == nil { log.PANIC("scanner not available") }
            log.Debug("initialize graphics")
            err := renderer.Init() 
            if err != nil {
                log.Fatal("could not initialize graphics: %s",err)    
            }
            texts := make(chan render.Text)
            go scanner.ScanText(texts)
            go renderer.ReadText(texts)
            renderer.Render()

        case Listen:
            if server == nil { log.PANIC("server not available") }
            if renderer == nil { log.PANIC("renderer not available") }
            server.Serve()
                    
        case Send:
            if client == nil { log.PANIC("client not available") }
            client.SendText()
            
        case Conf:
            if client == nil { log.PANIC("client not available") }
            client.SendConf()
            
        default:
            log.PANIC("inconsistent mode")
    }
        
        
    
}




func ShowModeHelp(mode Mode, flagSetMap map[Mode]*flag.FlagSet) {
    switches := ""
    flags := ""
    flagSetMap[mode].VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) == 1 { switches += " [ -"+f.Name+name+" ]" }
        if len(f.Name) >  1 { flags += " [ -"+f.Name+name+" ]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s%s%s\n",BUILD_NAME,mode,switches,flags)
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagSetMap[mode].PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}


func ShowHelp() {
    flags := ""
    flag.CommandLine.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) >=  1 { flags +=    " [ -"+f.Name+name+" ]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s   ",BUILD_NAME,flags)
    for _,mode := range modes {
        fmt.Fprintf(os.Stderr,"%s | ",mode)
    }
    fmt.Fprintf(os.Stderr,"%s\n",Version)
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    if render.RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"  %s    # %s\n",Read,"receive text and display")
    }
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Send,"send text to display")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Conf,"control display style")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Version,"show version info")
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}
    

func ShowVersion() {
    fmt.Printf(AUTHOR)
    fmt.Fprintf(os.Stderr,"\n%s, version %s for %s, built %s\n",BUILD_NAME,BUILD_VERSION,BUILD_PLATFORM,BUILD_DATE)
}    
    






