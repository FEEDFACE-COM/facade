package log


import (
    "fmt"
    "os"
)

type Verbosity uint

const (
    DEBUG     Verbosity = 7
    INFO      Verbosity = 6
    NOTICE    Verbosity = 5
    WARNING   Verbosity = 4
)


type Logger struct { verbosity Verbosity }
var logger Logger = Logger{verbosity: NOTICE}

func SetVerbosity(verbosity Verbosity) { logger.verbosity = verbosity }


func   Debug(format string, args ...interface{})  { if logger.verbosity >= DEBUG   { fmt.Fprintf(os.Stderr,format+"\n", args...) } }
func    Info(format string, args ...interface{})  { if logger.verbosity >= INFO    { fmt.Fprintf(os.Stderr,format+"\n", args...) } }
func  Notice(format string, args ...interface{})  { if logger.verbosity >= NOTICE  { fmt.Fprintf(os.Stderr,format+"\n", args...) } }
func Warning(format string, args ...interface{})  { if logger.verbosity >= WARNING { fmt.Fprintf(os.Stderr,format+"\n", args...) } }
func   Error(format string, args ...interface{})  {                                  fmt.Fprintf(os.Stderr,"ERROR: "+format+"\n", args...) } 
func   PANIC(format string, args ...interface{})  {                                  fmt.Fprintf(os.Stderr,"FACADE PANIC: "+format+"\n", args...); os.Exit(2) } 



