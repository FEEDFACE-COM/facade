
uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;



varying vec2 vTexCoord;
varying vec2 vTileCoord;



bool DEBUG    = debugFlag > 0.0;
bool DOWNWARD = downward == 1.0;

bool firstLine() {
	float t = 2.;
    if (mod(tileCount.y, 2.0) != 1.0 ) {
    	t = 1.;
    }
	return (  tileCount.y + vTileCoord.y - t) * 2. <= tileCount.y;
}


bool lastLine() { 
	return  vTileCoord.y*2.  > tileCount.y + 1. ;
}

bool newestLine() {
	return ! DOWNWARD && firstLine() || DOWNWARD && lastLine() ;
}

bool oldestLine() {
	return ! DOWNWARD && lastLine()  || DOWNWARD && firstLine();
}


void main() {
    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { 
        col = texture2D(texture, vTexCoord);
    }

	if ( newestLine() ) {
		col.rgba *= (1.0 - scroller);
	} else if ( oldestLine() ) {
		col.rgba *= scroller;
	}


    if (!gl_FrontFacing) {
		col.a /= 2.;
	}
	
    gl_FragColor = col;

}

