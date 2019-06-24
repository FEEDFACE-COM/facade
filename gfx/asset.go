
package gfx

import "fmt"
import "sort"

func ListShaderNames() []string {
    var ret []string
    for n, _ := range VertexShader {
        ret = append(ret,fmt.Sprintf("%s.vert",n)) 
    }
    for n, _ := range FragmentShader {
        ret = append(ret,fmt.Sprintf("%s.frag",n)) 
    }
    sort.Strings(ret)
    return ret
}


func ListFontNames() []string {
    var ret []string
    for n,_ := range VectorFont {
        ret = append(ret,n)
    }
    sort.Strings(ret)
    return ret
}


