
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
	col = texture2D(texture, vTexCoord); 
    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    }
    
    col.b =  -0.5 + (tileCount.x+vTileCoord.x) / tileCount.x;
    
    col.rg *= 0.75;
    
    if ( ! DOWNWARD && firstLine() || DOWNWARD && lastLine() ) {
        col.g = 0.0;
    }

    if ( ! DOWNWARD && lastLine() || DOWNWARD && firstLine() ) {
        col.r = 0.0;
//        col.b = 0.0;
    }
    
    if ( vTileCoord.y+1. >= tileCount.y ) {
    	col.rgb = vec3(0., 0., 1.);
   	}

    if (!gl_FrontFacing) {
        col.rgb = 0.35 * vec3(1., 1., 1.);
    } 
    
    
    gl_FragColor = col;
    

}
