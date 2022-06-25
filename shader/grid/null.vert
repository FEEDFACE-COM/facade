uniform mat4 model;               // model transformation
uniform mat4 view;                // view transformation
uniform mat4 projection;          // projection transformation

uniform vec2 tileSize;            // size of largest glyph in font, as (width,height)
uniform vec2 tileCount;           // grid dimensions, as (columns,rows)
uniform vec2 tileOffset;          // grid center offset from (0,0), as (columns,rows)

uniform float now;                // time since start of program, as seconds
uniform float scroller;           // scroller amount, from 0.0 to 1.0
uniform float debugFlag;          // 0.0 unless -D flag given by user

attribute vec3 vertex;            // vertex position as (x,y,z) centered on (0,0,0)

attribute vec2 tileCoord;         // tile coordinates, as (x,y) centered on (0,0)
attribute vec2 gridCoord;         // tile coordinates, as (column,row)
attribute vec2 texCoord;          // texture coordinates in font texture atlas


varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying vec2 vTexCoord;



bool DEBUG = debugFlag > 0.0;

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    
    vec4 pos = vec4(vertex,1);
    pos.x += tileCoord.x * tileSize.x;
    pos.y += tileCoord.y * tileSize.y;



    gl_Position = projection * view * model * pos;
}

