uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;
varying float vScroller;
varying float vDownward;

bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }


void main() {
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vDownward = downward;
    vScroller = scroller;
    vNow = now;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    
    
    float F = 0.25;
    float f0 = ease1( now /1.);  
    float f1 = 0.;
    
    

    // allow for scroller
    float from = cos( vTileCoord.y + 3. * now + PI/2.);
    float to =   cos( vTileCoord.y-1. + 3. * now + PI/2. );
    float delta =  to + scroller * (from - to);
    
    

    pos.z += F * cos( vTileCoord.x + 2. * now );
    pos.z += F * delta;


    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    
    gl_Position = projection * view * model * pos;
}
