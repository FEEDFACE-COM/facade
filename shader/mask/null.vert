

uniform float debugFlag;
uniform float ratio;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vTexCoord;
varying float vRatio;
varying float vDebugFlag;

bool DEBUG = debugFlag > 0.0;


void main() {
    vTexCoord = texCoord;
    vRatio = ratio;
    vDebugFlag = debugFlag;
    
    
    gl_Position = vec4(vertex,1);
}
