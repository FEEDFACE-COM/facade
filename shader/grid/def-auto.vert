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
varying float vScroller;

bool DEBUG = debugFlag > 0.0;


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

    float fontRatio = tileSize.x/tileSize.y;
    float screenRatio = (tileCount.x*tileSize.x)/(tileCount.y*tileSize.y);
    float ratio = screenRatio / fontRatio;
    float scaleWidth = ratio * 2. / tileCount.x;
    float scaleHeight =        2. / tileCount.y;
    

    float zoom = 1.;
    if ( scaleWidth < scaleHeight ) {
        zoom = scaleWidth;
    } else {            
        zoom = scaleHeight;
    }


//    pos.xyz *= zoom;
    pos.xyz *= model[0][0];
    gl_Position = projection * view * pos;
}

