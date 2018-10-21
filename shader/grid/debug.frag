
uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;
varying float vScroller;
varying float vDownward;

bool DEBUG    = vDebugFlag > 0.0;
bool downward = vDownward == 1.0;

bool firstLine() {
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        return 0.5*vTileCount.y       == vTileCoord.y ;
    } else {
        return  0.5*(vTileCount.y+1.) == vTileCoord.y + 1. ;
    }
}

bool lastLine() {
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        return -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
    } else {
        return -0.5*(vTileCount.y+1.) == vTileCoord.y - 1. ;
    }

}


void main() {

    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { 
        col = texture2D(texture, vTexCoord); 
    }
    
    
    if ( ! downward && firstLine() || downward && lastLine() ) {
        col.g = 0.0;
    }

    if ( ! downward && lastLine() || downward && firstLine() ) {
        col.r = 0.0;
        col.b = 0.0;
    }
    
    gl_FragColor = col;
}
