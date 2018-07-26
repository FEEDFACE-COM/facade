
package main


import (
    "fmt"
    "log"
    "os"
)

var VERBOSE bool = false
var DEBUG   bool = false

func Print(format string, args ...interface{}) { fmt.Fprintf(os.Stdout, format, args...) }

func   Log(format string, args ...interface{}) {                       log.Printf(format, args...)   }
func Debug(format string, args ...interface{}) { if DEBUG            { log.Printf(format, args...) } }
func  Info(format string, args ...interface{}) { if DEBUG || VERBOSE { log.Printf(format, args...) } }
func ERROR(format string, args ...interface{}) { log.Printf("3RR0R: "+format, args...) }
func FATAL(format string, args ...interface{}) { log.Printf("F4T4L: "+format, args...); os.Exit(2) }



