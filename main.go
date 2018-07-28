
package main

import (
    "fmt"
    "strings"
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
var modes = []Mode{Send,Conf}



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

    flags := make(map[Mode] *flag.FlagSet)

    if SERVER_MODE_AVAILABLE {
        modes = append(modes, Beam)
    }
    
    for _,mode := range modes {
        flags[mode] = flag.NewFlagSet(string(mode), flag.ExitOnError)
        flags[mode].Usage = func() { ShowModeHelp(mode,flags) }
    }

    for _,mode := range []Mode{Send,Conf} {
        flags[mode].UintVar(&connectPort, "port", connectPort, "connect to `port`" )
        flags[mode].StringVar(&connectHost, "host", connectHost, "connect to `host`" )
        flags[mode].Float64Var(&connectTimeout, "timeout", connectTimeout, "timeout after `seconds`") 
    }

    if SERVER_MODE_AVAILABLE {
        flags[Beam].UintVar(&listenPort, "port", listenPort, "listen on `port`" )
        flags[Beam].StringVar(&listenHost, "host", listenHost, "listen on `host`" )
        flags[Beam].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
    }

    all := []*flag.FlagSet{flag.CommandLine}
    for _,mode := range modes {
        all = append(all,flags[mode])
    }
    for _,flagSet := range all {
        flagSet.BoolVar(&VERBOSE,"v", false, "show verbose messages")
        flagSet.BoolVar(&DEBUG,  "d", false, "show debug messages")
    }
    
    flag.Parse()
    if flag.NArg() < 1 { 
        ShowHelp(); 
        os.Exit(-2) 
    }
    
    var client *Client
    var server *Server

    
    switch ( Mode(flag.Args()[0]) ) {

        case Beam:
            if !SERVER_MODE_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[Beam].Parse( flag.Args()[1:] )
            server = NewServer(listenHost,listenPort)
            server.Serve()
            
        case Send:
            flags[Send].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,connectPort,connectTimeout)
            client.text = strings.Join(flags[Send].Args()[0:], "  ")
            client.Send()
            
        case Conf:
            flags[Conf].Parse( flag.Args()[1:] )
            client = NewClient(connectHost,connectPort,connectTimeout)
            
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
    fmt.Fprintf(os.Stderr,"  %s %s   ",BUILD_NAME,flags)
    for _,mode := range modes {
        fmt.Fprintf(os.Stderr,"%s | ",mode)
    }
    fmt.Fprintf(os.Stderr,"%s\n",Version)
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    if SERVER_MODE_AVAILABLE {
        fmt.Fprintf(os.Stderr,"  %s    # %s\n",Beam,"receive text and display")
    }
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
    






