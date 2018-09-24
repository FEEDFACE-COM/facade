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
