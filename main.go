
package main

import (
    "fmt"
//    "strings"
    "flag"
    "os"    
    "os/signal"
    log "./log"
    render "./render"
    proto "./proto"
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
    Pipe    Command = "pipe"
    Conf    Command = "conf"    
    Version Command = "info"
    Help    Command = "help"
)
var cmds = []Command{Conf,Pipe}



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
    }

    for _,cmd := range []Command{Pipe} {
        flags[cmd].UintVar(&textPort, "textport", textPort, "connect to `port` for text" )
    }
    
    for _,cmd := range []Command{Pipe,Conf} {
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
    
    
    
    cmd := Command(flag.Args()[0])
    switch ( cmd ) {

        case Read:
            if !render.RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[Read].Usage = func() { ShowCommandHelp(Read,flags) }
            flags[Read].Parse( flag.Args()[1:] )
            renderer = render.NewRenderer()
            scanner = NewScanner()

        case Listen:
            if !render.RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[Listen].Usage = func() { ShowCommandHelp(Listen,flags) }
            flags[Listen].Parse( flag.Args()[1:] )
            server = NewServer(listenHost,confPort,textPort)
            renderer = render.NewRenderer()

        case Pipe:
            flags[Pipe].Usage = func() { ShowCommandHelp(Pipe,flags) }
            flags[Pipe].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
            
        case Conf:
            flags[Conf].Usage = func() { ShowCommandHelp(Conf,flags) }
            flags[Conf].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,confPort,textPort,connectTimeout)
            
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
    
    
    
    var config *proto.Config = nil
    args := flags[cmd].Args()
    if len(args) < 1 {
    } else {
        mode := proto.Mode(args[0])
        switch (mode) {
            
            case proto.PAGER:
                config = proto.NewConfig(mode)
                cflags := config.FlagSet()
                cflags.Usage = func() { ShowModeHelp(mode,cmd,cflags) }
                cflags.Parse( args[1:] )
            
            default:
                ShowHelp()
                os.Exit(-2)    
        }
    }
        
    
    
    
    
    
    
        
    log.Info(AUTHOR)

    switch ( cmd ) {

        case Read:
            if renderer == nil { log.PANIC("renderer not available") }
            if scanner == nil { log.PANIC("scanner not available") }
            if config == nil {
                config = proto.NewConfig(proto.PAGER)
            }
            renderer.Init() 
            renderer.Config(config)
            texts := make(chan proto.Text)
            go scanner.ScanText(texts)
            go renderer.ReadText(texts)
            renderer.Render()

        case Listen:
            if server == nil { log.PANIC("server not available") }
            if renderer == nil { log.PANIC("renderer not available") }
            if config == nil {
                config = proto.NewConfig(proto.PAGER)
            }
            renderer.Init() 
            renderer.Config(config)
            texts := make(chan proto.Text)
            confs := make(chan proto.Config)
            go server.ListenText(texts)
            go server.ListenConf(confs)
            go renderer.ReadText(texts)
            go renderer.ReadConf(confs)
            renderer.Render()
                    
        case Pipe:
            if client == nil { log.PANIC("client not available") }
            if config != nil {
                client.SendConf(config)
            }
            client.ScanAndSendText()
            
        case Conf:
            if client == nil { log.PANIC("client not available") }
            if config == nil { log.PANIC("config not available") }
            client.SendConf(config)
            
        default:
            log.PANIC("inconsistent command")
    }
        
        
    
}


func ShowModeHelp(mode proto.Mode, cmd Command, flagset *flag.FlagSet) {
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
    for _,m := range proto.Modes {
        fmt.Fprintf(os.Stderr,"%s|",m)
    }
    fmt.Fprintf(os.Stderr,"\n")
    fmt.Fprintf(os.Stderr,"\nCommands:\n")
    if render.RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"  %s    # %s\n",Read,"pipe stdin to display")
        fmt.Fprintf(os.Stderr,"  %s    # %s\n",Listen,"receive text and display")
    }
//    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Send,"send text to remote facade")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Pipe,"pipe stdin to remote facade")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Conf,"control remote facade")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Version,"show facade info")
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",proto.PAGER,"console pager")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",proto.CLOUD,"wordcloud")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",proto.SCROLL,"scroller")
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}
    

func ShowVersion() {
    fmt.Printf(AUTHOR)
    fmt.Fprintf(os.Stderr,"\n%s, version %s for %s, built %s\n",BUILD_NAME,BUILD_VERSION,BUILD_PLATFORM,BUILD_DATE)
}    
    






