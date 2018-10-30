
uniform float scroller;
uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;

uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

bool DEBUG    = debugFlag > 0.0;
bool DOWNWARD = downward == 1.0;


bool firstLine() {
    if (mod(tileCount.y, 2.0) != 1.0 ) { 
        return 0.5*tileCount.y       == vTileCoord.y ;
    } else {
        return  0.5*(tileCount.y+1.) == vTileCoord.y + 1. ;
    }
}

bool lastLine() {
    if (mod(tileCount.y, 2.0) != 1.0 ) { 
        return -0.5*tileCount.y + 1.0 == vTileCoord.y ;
    } else {
        return -0.5*(tileCount.y+1.) == vTileCoord.y - 1. ;
    }
}



void main() {

    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { 
        col = texture2D(texture, vTexCoord); 
    }
    
    if (DOWNWARD) {
        if ( firstLine() ) {
            col.rgba *= (-1.0 * scroller);
        }
        if ( lastLine() ) {
            col.rgba *= (1.0 - -1.0 * scroller);
        }
    } else { // ! downward
        if ( firstLine() ) {
            col.rgba *= (1.0 - scroller);
        }
        if ( lastLine() ) {
            col.rgba *= scroller; 
        }
    }
        
    gl_FragColor = col;
}
