
// +build linux,arm
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




"grid":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texcoord;
attribute vec2 tileCoord;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vDebugFlag;

bool DEBUG = debugFlag > 0.0;


void main() {
    vTexCoord = texcoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    
    vec4 pos = vec4(vertex,1);
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    
    if (mod(tileCount.x, 2.0) != 1.0 ) {
        pos.x -= tileSize.x/2.;
    }
    if (mod(tileCount.y, 2.0) != 1.0 ) {
        pos.y -= tileSize.y/2.;    
    }
    
    gl_Position = projection * view * model * pos;
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




"mask":`


uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texcoord;

varying vec2 vTexCoord;
varying float vDebugFlag;

bool DEBUG = debugFlag > 0.0;


void main() {
    vTexCoord = texcoord;
    vDebugFlag = debugFlag;
    
    gl_Position = vec4(vertex,1);
}
`,


}
