
uniform float debugFlag;

attribute vec3 vertex;


bool DEBUG = debugFlag > 0.0;


void main() {
    gl_Position = vec4(vertex,1.);
}
