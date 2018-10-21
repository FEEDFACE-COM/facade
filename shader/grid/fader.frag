
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

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;


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
        col.a = 1. - vScroller;    
    }
    
    if ( ! downward && lastLine() || downward && firstLine() ) {
       col.a = vScroller;
    }
    
    
    gl_FragColor = col;
}
