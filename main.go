
package main

import (
    "fmt"
    "strings"
    "flag"
    "os"    
    "os/signal"
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

type Command string
const (
    Read    Command = "read"
    Listen  Command = "recv"
    Send    Command = "send"
    Pipe    Command = "pipe"
    Conf    Command = "conf"    
    Version Command = "info"
    Help    Command = "help"
    Test    Command = "test"
)
var cmds = []Command{Send,Conf,Pipe}



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

    flags := make(map[Command] *flag.FlagSet)

    if render.RENDERER_AVAILABLE {
        cmds = append(cmds, Read)
        cmds = append(cmds, Listen)
    }
    
    for _,cmd := range cmds {
        flags[cmd] = flag.NewFlagSet(string(cmd), flag.ExitOnError)
        flags[cmd].Usage = func() { ShowCommandHelp(cmd,flags) }
    }


    for _,cmd := range []Command{Send,Pipe} {
        flags[cmd].UintVar(&textPort, "textport", textPort, "connect to `port` for text" )
    }
    
    for _,cmd := range []Command{Send,Pipe,Conf} {
        flags[cmd].UintVar(&confPort, "confport", confPort, "connect to `port` for config" )
        flags[cmd].StringVar(&connectHost, "host", connectHost, "connect to `host`" )
        flags[cmd].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout after `seconds`") 
    }

    if render.RENDERER_AVAILABLE {
        flags[Listen].UintVar(&confPort, "confport", confPort, "listen on `port` for config" )
        flags[Listen].UintVar(&textPort, "textport", textPort, "listen on `port` for text" )
        flags[Listen].StringVar(&listenHost, "host", listenHost, "listen on `host`" )
        flags[Listen].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
    }

    all := []*flag.FlagSet{flag.CommandLine}
    for _,cmd := range cmds {
        all = append(all,flags[cmd])
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
    var text string
    var conf string
    
    cmd := Command(flag.Args()[0])
    switch ( cmd ) {

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
            if len(flags[Send].Args()) < 1 {
                ShowHelp()
                os.Exit(-1)                    
            }
            text = strings.Join(flags[Send].Args()[0:], " ")
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
        case Pipe:
            flags[Pipe].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
            
        case Conf:
            flags[Conf].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            conf = strings.Join(flags[Conf].Args()[0:], " " )
            
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


    switch ( cmd ) {

        case Read:
            if renderer == nil { log.PANIC("renderer not available") }
            if scanner == nil { log.PANIC("scanner not available") }
            renderer.Init() 
            texts := make(chan render.Text)
            go scanner.ScanText(texts)
            go renderer.ReadText(texts)
            renderer.Render()

        case Listen:
            if server == nil { log.PANIC("server not available") }
            if renderer == nil { log.PANIC("renderer not available") }
            renderer.Init() 
            texts := make(chan render.Text)
            confs := make(chan render.Conf)
            go server.ListenText(texts)
            go server.ListenConf(confs)
            go renderer.ReadText(texts)
            go renderer.ReadConf(confs)
            renderer.Render()
                    
        case Send:
            if client == nil { log.PANIC("client not available") }
            client.SendText(text)
            
        case Pipe:
            if client == nil { log.PANIC("client not available") }
            client.ScanAndSendText()        
            
        case Conf:
            if client == nil { log.PANIC("client not available") }
            client.SendConf(conf)
            
        default:
            log.PANIC("inconsistent command")
    }
        
        
    
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
        if len(f.Name) >=  1 { flags +=    " [ -"+f.Name+name+" ]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s   ",BUILD_NAME,flags)
    for _,cmd := range cmds {
        fmt.Fprintf(os.Stderr,"%s | ",cmd)
    }
    fmt.Fprintf(os.Stderr,"%s\n",Version)
    fmt.Fprintf(os.Stderr,"\nCommands:\n")
    if render.RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"  %s    # %s\n",Read,"pipe stdin to display")
        fmt.Fprintf(os.Stderr,"  %s    # %s\n",Listen,"receive text and display")
    }
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Send,"send text to remote facade")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Pipe,"pipe stdin to remote facade")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Conf,"control remote facade")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Version,"show facade info")
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}
    

func ShowVersion() {
    fmt.Printf(AUTHOR)
    fmt.Fprintf(os.Stderr,"\n%s, version %s for %s, built %s\n",BUILD_NAME,BUILD_VERSION,BUILD_PLATFORM,BUILD_DATE)
}    
    






