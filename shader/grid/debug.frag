
uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying float vScroller;

bool DEBUG    = debugFlag > 0.0;
bool DOWNWARD = downward == 1.0;
bool EVENLINES = mod(tileCount.y, 2.0) != 1.0;


bool firstLine() {
	float t = 1.;
//    if (EVENLINES ) {
//    	t = 2.;
//    }
    
    
    
//    if (DOWNWARD) {
//    	return (  tileCount.y + vTileCoord.y - t) * 2. < tileCount.y;
//    } else {
//    	return (  tileCount.y + vTileCoord.y - t) * 2. <= tileCount.y;
//    }


    float d = 0.0;
    if (DOWNWARD) {
        d = 1.;
    }    
    return (  tileCount.y + vTileCoord.y - t) * 2. + d <= tileCount.y;
}


bool lastLine() { 
//    if (DOWNWARD) {
//    	return  vTileCoord.y*2.  > tileCount.y - 1. ;
//    } else {
//    	return  vTileCoord.y*2.  > tileCount.y;
//    }
    float t = 0.0;
    if (DOWNWARD) {
        t = 1.0;
    }
    return  vTileCoord.y*2.  > tileCount.y - t;
}

bool newestLine() {
	return  ( !DOWNWARD && firstLine() ) || ( DOWNWARD && lastLine() ) ;
}

bool oldestLine() {
	return ( !DOWNWARD && lastLine() ) || ( DOWNWARD && firstLine() );
}



void main() {
    vec4 col;
	col = texture2D(texture, vTexCoord); 
    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    }
    
    if ( newestLine() ) {
//		col.rgb = vec3(0.,1.,0.);
        col.rb *= 0.; //green
		col.a *= (1.-vScroller);
	}
	else if ( oldestLine() ) {
//		col.rgb = vec3(1.,0.,0.);
        col.gb *= 0.; //red
		col.a *= vScroller;
	}
	else {
//		col.rgb = vec3(0.0,0.0,1.0);
	}

//    if (!gl_FrontFacing) {
//		col.r = 1.0 - col.r;
//		col.g = 1.0 - col.g;
//		col.b = 1.0 - col.b;
//    } 
    
    gl_FragColor = col;
    

}
