
package gfx
var VertexShader = map[string]string{


"color":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform float debugFlag;

attribute vec3 vertex;
attribute vec4 color;

varying vec4  vFragColor;
varying float vDebugFlag;


void main() {
    vFragColor = color;
    vDebugFlag = debugFlag;
    gl_Position = projection * view * model * vec4(vertex, 1);
}
`,




"ident":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vFragCoord;
varying float vDebugFlag;

bool DEBUG = debugFlag > 0.0;

void main() {
    vFragCoord = texCoord;
    gl_Position = projection * view * model * vec4(vertex, 1);
}
`,




"grid/grid":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float timer;
uniform float scroller;
uniform float debugFlag;
uniform float downwardFlag;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying vec2 vTileCount;

varying float vDebugFlag;
varying float vDownwardFlag;
varying float vScroller;
varying float vTimer;


bool DEBUG = debugFlag > 0.0;





void main() {
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vDownwardFlag = downwardFlag;
    vScroller = scroller;
    vTimer = timer;

    
    vec4 pos = vec4(vertex,1);

    pos.y += vScroller;
    
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    

    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    
    gl_Position = projection * view * model * pos;
}
`,




"mask/mask":`


uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vTexCoord;
varying float vDebugFlag;

bool DEBUG = debugFlag > 0.0;


void main() {
    vTexCoord = texCoord;
    vDebugFlag = debugFlag;
    
    gl_Position = vec4(vertex,1);
}
`,


}


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
    
    bool firstLine =  0.5*vTileCount.y       == vTileCoord.y ;
    bool lastLine  = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;

    if (downward) {
        firstLine = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
        lastLine =   0.5*vTileCount.y       == vTileCoord.y ;
    }

    if (firstLine && scroll > 0.5) { //oldest line vanishes later
        col.rgb = col.rgb * (1.- 2.*(scroll-0.5));
    }

    if (lastLine) { //newest line blends in
        col.rgb = col.rgb * scroll;
    }    
    
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
