attribute vec2 texcoord;
attribute vec3 vertex;
attribute vec4 color;

varying vec4 fragcolor;
varying vec2 fragcoord;

void main() {
    fragcolor = vec4( vertex, 1.0);
    fragcoord = texcoord;
    gl_Position = vec4(vertex,1);
}
