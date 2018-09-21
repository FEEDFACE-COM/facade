
// +build linux,arm

package gfx

import (
    gl "src.feedface.com/gfx/piglet/gles2"
)    

type Quad struct {
    object uint32
    vertAttrib uint32
    texCoordAttrib uint32
}


func NewQuad(width float32, height float32) *Quad {
    ret := &Quad{}
    gl.GenBuffers(1,&ret.object)
    vrts := QuadVertices(1.0,height/width)
	gl.BindBuffer(gl.ARRAY_BUFFER, ret.object)
	gl.BufferData(gl.ARRAY_BUFFER, len(vrts)*4, gl.Ptr(vrts), gl.STATIC_DRAW)
	return ret
}

func (quad *Quad) VertexAttribPointer(program uint32) {
	quad.vertAttrib = uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(quad.vertAttrib) 
	gl.VertexAttribPointer(quad.vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	quad.texCoordAttrib = uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(quad.texCoordAttrib)
	gl.VertexAttribPointer(quad.texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
        
}

func (quad *Quad) Bind() { 
    gl.BindBuffer(gl.ARRAY_BUFFER,quad.object) 
}




func QuadVertices(w,h float32) []float32 {
    return []float32{
        

	//  X, Y, Z, U, V

//	// Front
    -1.0 * w/2,  1.0 * h/2, 0.0, 0.0, 0.0,
    -1.0 * w/2, -1.0 * h/2, 0.0, 0.0, 1.0,
     1.0 * w/2, -1.0 * h/2, 0.0, 1.0, 1.0,
     1.0 * w/2, -1.0 * h/2, 0.0, 1.0, 1.0,
     1.0 * w/2,  1.0 * h/2, 0.0, 1.0, 0.0,
    -1.0 * w/2,  1.0 * h/2, 0.0, 0.0, 0.0,

        
    }    
}