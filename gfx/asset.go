
package gfx

import "fmt"
import "sort"

func ListShaderNames() []string {
    var ret []string
    for n, _ := range ShaderAsset {
        ret = append(ret,fmt.Sprintf("%s",n)) 
    }
    sort.Strings(ret)
    return ret
}


func ListFontNames() []string {
    var ret []string
    for n,_ := range FontAsset {
        ret = append(ret,n)
    }
    sort.Strings(ret)
    return ret
}


