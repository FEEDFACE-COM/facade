uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;
uniform vec2 tileOffset;

uniform float now;
uniform float scroller;
uniform float debugFlag;

attribute vec3 vertex;

attribute vec2 texCoord;
attribute vec2 tileCoord;
attribute vec2 gridCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;

bool DEBUG = debugFlag > 0.0;


float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;


void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);

    float F = 1. + tileCount.y/8.;

    float y  =  vTileCoord.y       / (tileCount.y/2.);
    float yy = (vTileCoord.y + ((scroller)) ) / (tileCount.y/2.);


    float freq = -1.;
    float f0 = cos( freq * y  * PI + now + PI/2. );
    float f1 = cos( freq * yy * PI + now + PI/2. );
    float d =  f0 + /*(scroller) * */(f1 - f0);
    pos.z += F * d;
    pos.z -= F;


    
    gl_Position = projection * view * model * pos;
}
