
// +build linux,arm
package gfx
var FragmentShader = map[string]string{


"color":`


varying float vDebugFlag;
varying vec4 vFragColor;

bool DEBUG = vDebugFlag > 0.0;

void main() {
    gl_FragColor = vFragColor;
}
`,




"grid":`

uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vDebugFlag;
varying float vScroller;

bool DEBUG = vDebugFlag > 0.0;


void main() {
    float scroll = vScroller;
    
    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;

    vec4 col = texture2D(texture, tex);
    if (DEBUG) {
//        col = vec4(1.0,1.0,1.0,1.0);
        col = vec4(1.0,scroll,scroll,1.0);
    } 
    
    gl_FragColor = vec4(col.rgb,1.0);
}
`,




"ident":`
uniform sampler2D texture;

varying vec2 vFragCoord;
varying float vDebugFlag;

bool DEBUG = vDebugFlag > 0.0;


void main() {
    vec4 tex = texture2D(texture,vFragCoord);
    gl_FragColor = tex;
}
`,




"mask":`

varying vec2 vTexCoord;
varying float vDebugFlag;

bool DEBUG = vDebugFlag > 0.0;
float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.25) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}


void main() {

    vec2 pos = vTexCoord;
    vec4 col = vec4(0.0,0.0,0.0,0.0);

    if ( grid(pos) ) { col = vec4(1.,1.,1.,0.5); }
//    if ( pos.y > 0.0 && pos.y < 1.0 && abs(pos.x) <= w ) { col = vec4(0.,1.,0.,1.); }
//    if ( pos.x > 0.0 && pos.x < 1.0 && abs(pos.y) <= w ) { col = vec4(1.,0.,0.,1.); }
       
    gl_FragColor = col;
}
`,


}
