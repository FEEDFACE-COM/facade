uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileCount;
uniform vec2 tileSize;

uniform vec2 glyphCount;
uniform vec2 glyphSize;

attribute vec3 vertex;
attribute vec2 texcoord;

attribute vec2 tileCoord;
//attribute vec2 glyphCoord;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGlyphCoord;
varying vec2 vGlyphSize;

void main() {
    vTexCoord = texcoord;
    vTileCoord = tileCoord;
//    vGlyphCoord = glyphCoord;
    vGlyphCoord = vec2(1.,1.);
    vGlyphSize = glyphSize;

    
    vec4 pos = vec4(vertex,1);
  
//    pos.x += (gridcoord.x * tilesize.x) - gridsize.x;
//    pos.y += (gridcoord.y * tilesize.y) - gridsize.y;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    
//    if (mod(tileCount.x, 2.0) != 1.0 ) {
//        pos.x -= tileSize.x/2.;
//        pos.y -= tileSize.y/2.;    
//    }
    
    gl_Position = projection * view * model * pos;
}
