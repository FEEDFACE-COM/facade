
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
    float t = 0.0;
    if (DOWNWARD) {
        t = 1.;
    }    
    return (  tileCount.y + vTileCoord.y - 1.) * 2. + t <= tileCount.y;
}


bool lastLine() { 
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
        col.rb *= 0.; //green
		col.a *= (1.-vScroller);
	}
	else if ( oldestLine() ) {
        col.gb *= 0.; //red
		col.a *= vScroller;
	}


    gl_FragColor = col;
    
}
