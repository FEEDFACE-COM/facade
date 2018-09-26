
uniform vec2 fragCoord;
uniform float debugFlag;

attribute vec2 texCoord;
attribute vec3 vertex;
attribute vec4 color;


varying vec4 vFragColor;
varying vec2 vFragCoord;
varying float vDebugFlag;


bool DEBUG = debugFlag > 0.0;

void main() {
    vFragColor = vec4( vertex, 1.0);
    vFragCoord = texCoord;
    vDebugFlag = debugFlag;
    gl_Position = vec4(vertex,1);
}
