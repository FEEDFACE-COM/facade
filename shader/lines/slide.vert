uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;
uniform vec2 tileOffset;

uniform float now;
uniform float scroller;
uniform float screenRatio;
uniform float fontRatio;
uniform float debugFlag;

attribute vec3 vertex;

attribute vec2 texCoord;
attribute vec2 tileCoord;
attribute vec2 gridCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying float vScroller;

bool DEBUG = debugFlag > 0.0;


float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;



void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);
    
    
    float zoom;
    float ratio = screenRatio / fontRatio;
    
    float cols = ratio * 2. / tileCount.x;
    float rows = model[0][0];
    

    zoom = rows;


    float p = (pos.y + (tileSize.y*tileCount.y/2.)) / (tileSize.y * tileCount.y); 
    
    if ( p < 0.5) {
        pos.x *= cols/(1.+4.*(p));
    } else {
        pos.x *= cols/(1.+4.*(1.-p));
    
    }
    pos.y *= rows;
    

    gl_Position = projection * view * mat4(1.0) * pos;
}

