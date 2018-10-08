
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
    gfx "./gfx"
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
var commands = []Command{CONF,PIPE,TEST}



var (
    textPort       uint     = 0xfcd
    confPort       uint     = 0xfcc
    connectHost    string   = "fcd.hq.feedface.com"
    connectTimeout float64  = 5.0
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

    if RENDERER_AVAILABLE {
        commands = append(commands, READ)
        commands = append(commands, RECV)
    }
    
    for _,cmd := range commands {
        flags[cmd] = flag.NewFlagSet(string(cmd), flag.ExitOnError)
    }

    for _,cmd := range []Command{PIPE} {
        flags[cmd].UintVar(&textPort, "tp", textPort, "connect to `port` for text" )
    }
    
    for _,cmd := range []Command{PIPE,CONF} {
        flags[cmd].UintVar(&confPort, "cp", confPort, "connect to `port` for config" )
        flags[cmd].StringVar(&connectHost, "h", connectHost, "connect to `host`" )
        flags[cmd].Float64Var(&connectTimeout, "t", connectTimeout, "timeout after `seconds`") 
    }

    if RENDERER_AVAILABLE {
        flags[RECV].UintVar(&confPort, "cp", confPort, "listen on `port` for config" )
        flags[RECV].UintVar(&textPort, "tp", textPort, "listen on `port` for text" )
        flags[RECV].StringVar(&listenHost, "h", listenHost, "listen on `host`" )
//        flags[RECV].BoolVar(&daemonize, "D",         daemonize, "daemonize" )
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
    
    
    
    cmd := Command(flag.Args()[0])
    switch ( cmd ) {

        case READ:
            if !RENDERER_AVAILABLE {
                ShowHelp()
                os.Exit(-2)    
            }
            flags[READ].Usage = func() { ShowHelpCommand(READ,flags) }
            flags[READ].Parse( flag.Args()[1:] )
            renderer = NewRenderer(directory)
            scanner = NewScanner()

        case RECV:
            if !RENDERER_AVAILABLE {
                ShowHelp()
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
            ShowAssets()
            os.Exit(-2)

        case HELP:
            ShowHelp()
            os.Exit(-2)
            
        
        case TEST:
            flags[TEST].Usage = func() {ShowHelpCommand(TEST,flags) }
            flags[TEST].Parse( flag.Args()[1:] )
            tester = NewTester(directory)

        default:
            ShowHelp()
            os.Exit(-2)
    }
    
    
    
    var config *facade.Config = facade.NewConfig( facade.DEFAULT_MODE )
    var modeflags *flag.FlagSet = config.FlagSet()
    args := flags[cmd].Args()
    
    //no more args after cmd
    if len(args) < 1 {
        if cmd == CONF { //conf needs args so bail
            ShowHelpCommand(CONF,flags)
            os.Exit(-2)
        } else { // otherwise huh?
            
                modeflags.Usage = func() { ShowHelpMode(facade.Mode(config.Mode),cmd,modeflags) }
                modeflags.Parse( args[0:] )
                
        }
         
            
    } else {
        mode := facade.Mode(args[0])
        switch facade.Mode(mode) {
            
            case facade.GRID:
                config = facade.NewConfig(mode)
                modeflags = config.FlagSet()
                modeflags.Usage = func() { ShowHelpMode(facade.Mode(mode),cmd,modeflags) }
                modeflags.Parse( args[1:] )

            case facade.LINES:
                config = facade.NewConfig(mode)
                modeflags = config.FlagSet()
                modeflags.Usage = func() { ShowHelpMode(facade.Mode(mode),cmd,modeflags) }
                modeflags.Parse( args[1:] )

            case facade.TEST:
                config = facade.NewConfig(mode)
                modeflags = config.FlagSet()
                modeflags.Usage = func() { ShowHelpMode(facade.Mode(mode),cmd,modeflags) }
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


func ShowHelpMode(mode facade.Mode, cmd Command, flagset *flag.FlagSet) {
    switches := "-"
    flags := ""
    flagset.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) == 1 && name == "" { switches += f.Name }
        if len(f.Name) >  1 || name != "" { flags += " [-"+f.Name+name+"]" }
    })
    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s %s [%s]%s\n",BUILD_NAME,cmd,mode,switches,flags)
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagset.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}



func ShowHelpCommand(cmd Command, flagSetMap map[Command]*flag.FlagSet) {
    modes := JoinModes(facade.Modes,"|")
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
    fmt.Fprintf(os.Stderr,"  %s %s [%s] %s  %s\n",BUILD_NAME,cmd,switches,flags,modes)
    ShowModes()

    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flagSetMap[cmd].PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
}

func ShowCommands() {
    fmt.Fprintf(os.Stderr,"\nCommands:\n")
    if RENDERER_AVAILABLE {
        fmt.Fprintf(os.Stderr,"%6s     %s\n",READ,"read stdin and display")
        fmt.Fprintf(os.Stderr,"%6s     %s\n",RECV,"receive and display ")
    }
    fmt.Fprintf(os.Stderr,"%6s     %s\n",PIPE,"pipe stdin to remote")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",CONF,"configure remote")
}


func ShowModes() {
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr,"%6s     %s\n",facade.GRID,"character grid")
}

func ShowHelp() {
    flags := ""
    flag.CommandLine.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        if len(f.Name) >=  1 { flags +=    " [-"+f.Name+name+"]" }
    })
    cmds := JoinCommands(commands,"|")
    modes := JoinModes(facade.Modes,"|")

    ShowVersion()
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s  %s  %s\n",BUILD_NAME,flags,cmds,modes)
    ShowCommands()
    ShowModes()
    fmt.Fprintf(os.Stderr,"\nFlags:\n")
    flag.PrintDefaults()
    fmt.Fprintf(os.Stderr,"\n")
    
    
}
    
func JoinCommands(cmds []Command, sep string) string { 
    var strs []string 
    for _,cmd := range cmds { strs = append(strs,string(cmd)) }
    return strings.Join(strs,sep)
}

func JoinModes(modes []facade.Mode, sep string) string { 
    var strs []string 
    for _,mode := range modes { strs = append(strs,string(mode)) }
    return strings.Join(strs,sep)
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
    






