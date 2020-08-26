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

float dy(float y) {
    return 1. - (  0.25 + 0.5 * ((y+tileCount.y/2.) / tileCount.y ) );
}

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    
    float w0 = dy(tileCoord.y+1.) * tileSize.x;
    float w1 = dy(tileCoord.y)    * tileSize.x;
    float h  = tileSize.y;
    

    vec4 pos = vec4(0.,0.,0.,1.);

    vec2 A = vec2((w0*(tileCoord.x-1.)) , (h * (tileCoord.y   )));
    vec2 B = vec2((w1*(tileCoord.x-1.)) , (h * (tileCoord.y-1.)));
    vec2 C = vec2((w1*(tileCoord.x   )) , (h * (tileCoord.y-1.)));
    vec2 D = vec2((w0*(tileCoord.x   )) , (h * (tileCoord.y   )));
    
    
   
    if        ( vertex.x < 0. && vertex.y > 0. ) {
        pos.xy = A;
    } else if ( vertex.x < 0. && vertex.y < 0. ) {
        pos.xy = B;
    } else if ( vertex.x > 0. && vertex.y < 0. ) {
        pos.xy = C;
    } else if ( vertex.x > 0. && vertex.y > 0. ) {
        pos.xy = D;
    }
    
    pos.xy += (scroller * (A-B));


    float ratio = screenRatio / fontRatio;
    float zoom_cols = ratio * 2. / tileCount.x * 1.5;
    float zoom_rows = 2./(tileCount.y+1.);
    float zoom = min(zoom_cols,zoom_rows);
    

    mat4 mdl;
    mdl = mat4(1.0);
    mdl[0][0] = zoom;
    mdl[1][1] = zoom;

    gl_Position = projection * view * mdl * pos;
}

