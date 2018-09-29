

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
