
package main

import (
    "fmt"
    "flag"
    "log"
    "os"    
)

//const AUTHOR = "FACADE by FEEDFACE.COM"
//const AUTHOR = 
//`   _   _   _   _   _   _      _   _   _   _   _   _   _   _     _   _           
//  |_  |_| /   |_| | \ |_     |_  |_  |_  | \ |_  |_| /   |_    /   / \ |\/|  
//  |   | | \_  | | |_/ |_  BY |   |_  |_  |_/ |   | | \_  |_  * \_  \_/ |  |  
//
//`
const AUTHOR = 
`   _   _   _   _   _   _
  |_  |_| /   |_| | \ |_
  |   | | \_  | | |_/ |_  by FEEDFACE.COM

`


var BUILD_NAME, BUILD_VERSION, BUILD_PLATFORM, BUILD_DATE string

type Mode string
const (
    Test    Mode = "test"
    Version Mode = "version"
    Help    Mode = "help"
)





func main() {
    log.SetFlags(0)
    
    flag.Usage = ShowHelp
    testFlags   := flag.NewFlagSet(string(Test), flag.ExitOnError)
    for _,elem := range []*flag.FlagSet{flag.CommandLine,testFlags} {
        elem.BoolVar(&VERBOSE,"v", false, "show verbose messages")
        elem.BoolVar(&DEBUG,  "d", false, "show debug messages")
    }
    flag.Parse()
    


    if flag.NArg() < 1 { 
        ShowHelp(); 
        os.Exit(-2) 
    }
        
    switch ( Mode( flag.Args()[0] ) ) {
        case Test:
            err := doTest(testFlags, flag.Args()[1:])
            if err != nil {
                ERROR("could not test: %s",err)
                os.Exit(-1)    
            }
        case Version:
            ShowVersion()
            os.Exit(-2)
        case Help:
            ShowHelp()
            os.Exit(2)
        default:
            ShowHelp()
            os.Exit(-2)
    }
}

func doTest(testFlags *flag.FlagSet, testArgs []string) error {
    
    
    testFlags.Usage = func() {
        switches := ""
        flags := ""
        testFlags.VisitAll( func(f *flag.Flag) { 
            name,_ := flag.UnquoteUsage(f)
            if name != "" { name = "="+name }
            if len(f.Name) == 1 { switches += " [ -"+f.Name+name+" ]" }
            if len(f.Name) >  1 { flags += " [ -"+f.Name+name+" ]" }
        })
        fmt.Fprintf(os.Stderr,AUTHOR,)
        fmt.Fprintf(os.Stderr,"\nUsage:\n")
        fmt.Fprintf(os.Stderr,"  %s %s%s%s\n",BUILD_NAME,Test,switches,flags)
        fmt.Fprintf(os.Stderr,"\nFlags:\n")
        testFlags.PrintDefaults()
        fmt.Fprintf(os.Stderr,"\n")
    }
    
    passPtr := testFlags.Bool("pass", false, "autopass test")
    thingPtr := testFlags.String("do", "nothing", "the `thing` to do")
    errPtr  := testFlags.Bool("E", false, "error")
    failPtr := testFlags.Bool("F", false, "fail")
    testFlags.Parse(testArgs)

    Debug(AUTHOR)
    Info("doing %s...",*thingPtr)
    if *errPtr { 
        return fmt.Errorf("ERROR ERROR ERROR")
    }
    if *failPtr {
        FATAL("FAIL FAIL FAIL")
    }
    if *passPtr  {
        Log("...TEST PASSED!")
    }
    return nil
}



func ShowHelp() {
    flags := ""
    flag.CommandLine.VisitAll( func(f *flag.Flag) { 
        name,_ := flag.UnquoteUsage(f)
        if name != "" { name = "="+name }
        flags += " [ -"+f.Name+name+" ]"
    })
    fmt.Fprintf(os.Stderr,AUTHOR,)
    fmt.Fprintf(os.Stderr,"\nUsage:\n")
    fmt.Fprintf(os.Stderr,"  %s %s %s | %s \n",BUILD_NAME,flags,Test,Version)
    fmt.Fprintf(os.Stderr,"\nModes:\n")
    fmt.Fprintf(os.Stderr," %8s    # %s\n",Test,"test the testing test")
    fmt.Fprintf(os.Stderr," %8s    # %s\n",Version,"show version")
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
    

