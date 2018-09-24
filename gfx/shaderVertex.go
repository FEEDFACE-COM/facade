
// +build linux,arm
package gfx
var VertexShader = map[string]string{


"color":`
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
`,




"grid":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 gridsize;

attribute vec3 vertex;
attribute vec2 texcoord;

attribute vec2 gridcoord;
attribute vec2 texoffset;

varying vec2 offset;
varying vec2 fragcoord;
varying vec2 fraggridcoord;

void main() {
    fragcoord = texcoord.xy;
    fraggridcoord = gridcoord;
    offset = texoffset;
    vec4 pos = vec4(vertex,1);
  
    pos.x += gridcoord.x - gridsize.x;
    pos.y += gridcoord.y - gridsize.y;
    
    gl_Position = projection * view * model * pos;
}
`,




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
