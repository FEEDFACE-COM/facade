uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texcoord;
attribute vec2 tileCoord;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vDebugFlag;

void main() {
    vTexCoord = texcoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    
    vec4 pos = vec4(vertex,1);
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    
    
    if (mod(tileCount.x, 2.0) != 1.0 ) {
        pos.x -= tileSize.x/2.;
        pos.y -= tileSize.y/2.;    
    }
    
    gl_Position = projection * view * model * pos;
}
