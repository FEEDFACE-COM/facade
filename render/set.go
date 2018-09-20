
package render

import (
    "fmt"
)

type Set struct {
    items map[string] *Item
}


type Item struct {
    text string    
}

func NewSet() *Set {
    ret := &Set{}
    ret.items = make(map[string] *Item )
    return ret    
}


func (set *Set) Desc() string { return fmt.Sprintf("set[%d]",len(set.items) ) }
func (set *Set) Dump() string { 
    ret := ""
    for item := range set.items {
        ret = ret + " " +set.items[item].text    
    }
    return ret
}