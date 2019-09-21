
// +build linux,arm

package gfx

import (
//    gl "github.com/FEEDFACE-COM/piglet/gles2"
)    



func QuadVertices(w,h float32) []float32 {
    return []float32{
        

/*
    A        D
     +------+
     |      |
     |      |
     +------+
    B        C

	         x,          y,   z,           u,   v,
*/
    -1.0 * w/2,  1.0 * h/2, 0.0,         0.0, 0.0,    // A
    -1.0 * w/2, -1.0 * h/2, 0.0,         0.0, 1.0,    // B
     1.0 * w/2, -1.0 * h/2, 0.0,         1.0, 1.0,    // C
     1.0 * w/2, -1.0 * h/2, 0.0,         1.0, 1.0,    // C
     1.0 * w/2,  1.0 * h/2, 0.0,         1.0, 0.0,    // D
    -1.0 * w/2,  1.0 * h/2, 0.0,         0.0, 0.0,    // A

        
    }    
}