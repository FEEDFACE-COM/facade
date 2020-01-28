
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vTexCoord;



bool DEBUG = debugFlag > 0.0;


void main() {
    vTexCoord = texCoord;
    vec4 pos = vec4(vertex,1);
    gl_Position = projection * view * model * pos;
}

