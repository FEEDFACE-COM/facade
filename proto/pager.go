
package proto

import (
    "fmt"
    "flag"
)


type Pager struct {
    Width  uint
    Height uint    
}
func NewPager() *Pager {
    ret := &Pager{Width:16, Height:8}
    return ret
}
func (pager *Pager) AddFlags(flags *flag.FlagSet) {
    flags.UintVar(&pager.Width,"width",pager.Width,"pager width")    
    flags.UintVar(&pager.Height,"height",pager.Height,"pager height")    
}
func (pager *Pager) Desc() string {
    ret := fmt.Sprintf("pager[%dx%d]",pager.Width,pager.Height)
    return ret
}
