uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

attribute vec3 vertex;
attribute vec4 color;

varying vec4 fragcolor;

void main() {
    fragcolor = color;
    gl_Position = projection * view * model * vec4(vertex, 1);
}
