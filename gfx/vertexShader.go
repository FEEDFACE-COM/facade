
// +build linux,arm
package gfx
var VertexShader = map[string]string{


"ident":`
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
`,




"mask":`
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
`,


}
