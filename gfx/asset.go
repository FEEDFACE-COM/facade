
package gfx

import "fmt"

func ListShaderNames() []string {
    var ret []string
    for n, _ := range VertexShader {
        if FragmentShader[n] == "" {
            ret = append(ret,fmt.Sprintf("%s.vert",n)) 
        } else {
            ret = append(ret,fmt.Sprintf("%s.vert %s.frag",n,n)) 
        }        
    }
    for n, _ := range FragmentShader {
        if VertexShader[n] == "" {
            ret = append(ret,fmt.Sprintf("%s.frag",n)) 
        }        
    }
    return ret
}


func ListFontNames() []string {
    var ret []string
    for n,_ := range VectorFont {
        ret = append(ret,n)
    }
    return ret
}


