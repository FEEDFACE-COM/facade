uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float scroller;
uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texcoord;
attribute vec2 tileCoord;

attribute float totalWidth;

varying vec2 vTexCoord;
varying vec2 vTileCoord;


varying float vDebugFlag;
varying float vScroller;



bool DEBUG = debugFlag > 0.0;



void main() {
    vTexCoord = texcoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vScroller = scroller;
    
    vec4 pos = vec4(vertex,1);

    pos.y += vScroller;    
    
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    
    gl_Position = projection * view * model * pos;
}
