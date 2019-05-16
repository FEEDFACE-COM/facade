
package log

import (
    "fmt"
    "errors"    
)


func NewError(format string, args ...interface{}) error { 
    return errors.New( fmt.Sprintf(format,args...) ) 
}

