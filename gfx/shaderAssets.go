
package gfx
var VertexShader = map[string]string{


"identity":`
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




"mask/identity":`


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




"grid/identity":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
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
varying float vNow;


bool DEBUG = debugFlag > 0.0;

void main() {
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vNow = now;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    //
    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    gl_Position = projection * view * model * pos;
}
`,




"grid/sinezfield":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
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
varying float vNow;


bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }


void main() {
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vNow = now;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    
    
    float F = 0.25;
    float f0 = ease1( pos.y + now );  
    float f1 = 0.;
    
    
    pos.z += F * cos( pos.x + 2. * now         );
    pos.z += F * cos( pos.y + 3. * now + PI/2. );
    pos.z += F * f0;


    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    
    gl_Position = projection * view * model * pos;
}
`,


}


var FragmentShader = map[string]string{


"identity":`
uniform sampler2D texture;

varying vec2 vFragCoord;
varying float vDebugFlag;

bool DEBUG = vDebugFlag > 0.0;


void main() {
    vec4 tex = texture2D(texture,vFragCoord);
    gl_FragColor = tex;
}
`,




"color":`


varying float vDebugFlag;
varying vec4 vFragColor;

bool DEBUG = vDebugFlag > 0.0;

void main() {
    gl_FragColor = vFragColor;
}
`,




"mask/identity":`

varying vec2 vTexCoord;
varying float vDebugFlag;

bool DEBUG = vDebugFlag > 0.0;


void main() {

    vec2 pos = vTexCoord;
    vec4 col = vec4(0.0,0.0,0.0,0.0);

    gl_FragColor = col;
}
`,




"mask/debug":`

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




"grid/identity":`

uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;

bool DEBUG    = vDebugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;



void main() {

    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { col = 
        texture2D(texture, vTexCoord); 
    }
    
    bool firstLine,lastLine;
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        firstLine =  0.5*vTileCount.y       == vTileCoord.y ;
        lastLine  = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
    } else {
        firstLine =  0.5*(vTileCount.y+1.) == vTileCoord.y + 1. ;
        lastLine  = -0.5*(vTileCount.y+1.) == vTileCoord.y - 1. ;
    }

    if (DEBUG && firstLine ) {
        col.rgb = vec3(1.,0.,1.);
    }

    if (DEBUG && lastLine ) {
        col.rgb = vec3(0.,1.,0.);
    }
    
    gl_FragColor = col;
}
`,


}
