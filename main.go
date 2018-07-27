
package main

import (
    "fmt"
    "flag"
    "log"
    "os"    
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
    Beam    Mode = "beam"
    Send    Mode = "send"
    Conf    Mode = "conf"    
    Version Mode = "info"
    Help    Mode = "help"
    Test    Mode = "test"
)



var (
    connectPort    uint     = 0xfcd
    connectHost    string   = "localhost"
    connectTimeout float64  = 5.0
    listenPort     uint     = 0xfcd
    listenHost     string   = "0.0.0.0"
    daemonize      bool     = false
)


func main() {
    log.SetFlags(0)
    
    flag.Usage = ShowHelp


    var flags map[Mode] *flag.FlagSet
    
    flags = make(map[Mode] *flag.FlagSet)
    flags[Beam] = flag.NewFlagSet(string(Beam), flag.ExitOnError)
    flags[Send] = flag.NewFlagSet(string(Send), flag.ExitOnError)
    flags[Conf] = flag.NewFlagSet(string(Conf), flag.ExitOnError)

    for _,elem := range []*flag.FlagSet{flag.CommandLine,flags[Beam],flags[Send],flags[Conf]} {
        elem.BoolVar(&VERBOSE,"v", false, "show verbose messages")
        elem.BoolVar(&DEBUG,  "d", false, "show debug messages")
    }
    
    for _,mode := range []Mode{Send,Conf} {
        flags[mode].UintVar(&connectPort, "port", connectPort, "connect to `port`" )
        flags[mode].StringVar(&connectHost, "host", connectHost, "connect to `host`" )
        flags[mode].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout after `seconds`") 
    }

    for _,mode := range []Mode{Beam} {
        flags[mode].UintVar(&listenPort, "port", listenPort, "listen on `port`" )
        flags[mode].StringVar(&listenHost, "host", listenHost, "listen on `host`" )
        flags[mode].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
    }
    
    flag.Parse()
    if flag.NArg() < 1 { 
        ShowHelp(); 
        os.Exit(-2) 
    }
    
    var beamer *FcdBeamer
    var sender *FcdSender
    var confer *FcdConfer

    
    switch ( Mode(flag.Args()[0]) ) {

        case Beam:
            flags[Beam].Usage = func() { ShowModeHelp(Beam,flags) }
            flags[Beam].Parse( flag.Args()[1:] )
            beamer = NewFcdBeamer(listenHost,listenPort)
            
        case Send:
            flags[Send].Usage = func() { ShowModeHelp(Send,flags) }
            flags[Send].Parse( flag.Args()[1:] )
            sender = NewFcdSender(connectHost,connectPort,connectTimeout)
            
        case Conf:
            flags[Conf].Usage = func() { ShowModeHelp(Conf,flags) }
            flags[Conf].Parse( flag.Args()[1:] )
            confer = NewFcdConfer(connectHost,connectPort,connectTimeout)
            
        case Test:
            FATAL("TEST TEST TEST")

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
    
    if beamer != nil {
        beamer.beam()    
    }
    
    if sender != nil {
        sender.send()
    }
    
    if confer != nil {
        confer.conf()
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
    fmt.Fprintf(os.Stderr,AUTHOR,)
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
    fmt.Fprintf(os.Stderr,AUTHOR,)
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s   %s | %s | %s\n",BUILD_NAME,flags,Beam,Send,Conf)
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Beam,"receive text and display")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Send,"send text to display")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Conf,"control display style")
    fmt.Fprintf(os.Stderr,"  %s    # %s\n",Version,"show version info")
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}
    

func ShowVersion() {
    fmt.Fprintf(os.Stderr, 
`%s
%s, version %s for %s, built %s
`, AUTHOR, BUILD_NAME, BUILD_VERSION, BUILD_PLATFORM, BUILD_DATE )
}    
    






