
uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;

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
	col = texture2D(texture, vTexCoord); 
    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    }
    
    if ( newestLine() ) {
		col.r = 0.0;
	}
	else if ( oldestLine() ) {
		col.g = 0.0;
	}
	else {
		col.rg = vec2(0.0,0.0);
	}

    if (gl_FrontFacing) {
		col.r = 1.0 - col.r;
		col.g = 1.0 - col.g;
		col.b = 1.0 - col.b;
    } 
    
    gl_FragColor = col;
    

}
