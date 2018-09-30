
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




"grid/grid":`

uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;


varying vec2 vTileCount;
varying float vDownwardFlag;
varying float vDebugFlag;
varying float vScroller;
varying float vTimer;

bool DEBUG    = vDebugFlag > 0.0;
bool downward = vDownwardFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;


void main() {
    float scroll = abs(vScroller);
    
    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;

    vec4 col = texture2D(texture, tex);
//    if (DEBUG) {
//        col = vec4(1.0,1.0,1.0,1.0);
//        col = vec4(1.0,scroll,scroll,1.0);
//    } 
    
    
    bool lastLine = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;

    if (downward) {
        lastLine =  0.5*vTileCount.y == vTileCoord.y ;
    }

    if (lastLine) {
//        col.rgb = col.rgb * scroll;
//        col.rgb = vec3( col.r * scroll,col.g * scroll,col.b * scroll);
//        col.rg = vec2(0.,0.);
//        col.r = 1.0 * scroll;
col.r=1.0;
        col.a = 1.0;
    }    


//    if ( -0.5*vTileCount.y + 1.0 == vTileCoord.y ) { //last row blends in
//        col.rgb = vec3( col.r * scroll,col.g * scroll,col.b * scroll);
//        col.rg = vec2(0.,0.);
//        col.r = 1.0;
//        col.a = 1.0;
//    }

    

    
    gl_FragColor = col;
}
`,




"mask/mask":`

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
