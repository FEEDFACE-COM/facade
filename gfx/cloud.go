
package gfx

import (
    "fmt"
)
type Cloud struct {
    font Font    
}


func NewCloud() *Cloud {
    ret := &Cloud{}
    
    return ret
}


func (cloud *Cloud) Render() {}


func (cloud *Cloud) Desc() string { return fmt.Sprintf("cloud[]") }

