package log


import (
    "fmt"
    "os"
    "errors"
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


func NewError(format string, args ...interface{}) error { return errors.New( fmt.Sprintf(format,args...) ) }



func Dump(in []byte, count,offset int) string {
    off := offset % (4*4)
    ret := ""
    left,right := "",""
    //first line spaces for offset
    for i:=0;i<off;i++ {
        left  += "  "
        right += " " 
        if (i+1) % 4 == 0 {
            left += " "
        } else if (i+1) % (4*4) != 0 {
            left += " "
        }
    }
    for i,s := range(in) {
        
        if count != 0 && i >= count {
            break    
        }
        
        if i>0 && (i+off) % (4*4) == 0 {
            ret += left + "    " + right + "\n"
            left,right = "",""
        }
        left  += fmt.Sprintf("%02x",s)
        if s >= 0x20 && s <= 0x7f {
            right += fmt.Sprintf("%c",s)
        } else {
            right += "."
        }
        if (i+off+1) % 4 == 0 {
            left += " "
        } else if (i+off+1) % (4*4) != 0 {
            left += ":"
        }
    
    }
    //fill up remaining space
    for i:=len(left);i<len("00:00:00:00 00:00:00:00 00:00:00:00 00:00:00:00 ");i++ {
        left += " "    
    }
    for i:=len(right);i<len("XXXXXXXXXXXXXXXX");i++ {
        right += " "    
    }
    ret += left + "    " + right
    return ret 
}

