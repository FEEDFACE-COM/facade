uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

attribute vec3 vertex;
attribute vec2 texcoord;

varying vec2 fragcoord;

void main() {
    fragcoord = texcoord;
    gl_Position = projection * view * model * vec4(vertex, 1);
}
