uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;
uniform vec2 tileOffset;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;
attribute vec2 gridCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;

varying float vScroller;

bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;


void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vScroller = abs(scroller);
    vGridCoord = gridCoord;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);

    float F = 0.75;
    float x = pos.x;
    float y = pos.y;
    
    float freq = -1./24.;
    pos.y += F * cos( 2. * freq * x * PI + now         );
    pos.x += F * cos( 3. * freq * y * PI + now + PI/2. );
	
	pos.z += F * cos( 5. * freq * (x+y) * PI + now + PI/2. );

    gl_Position = projection * view * model * pos;
}

