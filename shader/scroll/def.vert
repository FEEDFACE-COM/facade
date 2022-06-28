uniform mat4 model;               // model transformation
uniform mat4 view;                // view transformation
uniform mat4 projection;          // projection transformation

uniform float now;                // time since start of program, as seconds
uniform float debugFlag;          // 0.0 unless -D flag given by user

attribute vec3 vertex;            // vertex position as (x,y,z) centered on (0,0,0)


bool DEBUG = debugFlag > 0.0;

void main() {
    vec4 pos = vec4(vertex,1);
    gl_Position = projection * view * model * pos;
}
