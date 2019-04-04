
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
bool ODDLINES = mod(tileCount.y, 2.0) == 1.0;


bool firstLine() {
    float s = 0.0;
    float t = 0.0;
    if ( ! DOWNWARD && ! ODDLINES) { s = 0.0; t = 0.0; }
    if ( ! DOWNWARD &&   ODDLINES) { s = 1.0; t = 0.0; }
    if (   DOWNWARD && ! ODDLINES) { s = 0.0; t = 1.0; }
    if (   DOWNWARD &&   ODDLINES) { s = 0.0; t = 1.0; }
    return (  tileCount.y + (vTileCoord.y-s) - 1.) * 2. + t <= tileCount.y;
}


bool lastLine() { 
    float s = 0.0;
    float t = 0.0;
    if ( ! DOWNWARD && ! ODDLINES) { s = 0.0; t = 0.0; }
    if ( ! DOWNWARD &&   ODDLINES) { s = 1.0; t = 0.0; }
    if (   DOWNWARD && ! ODDLINES) { s = 0.0; t = 1.0; }
    if (   DOWNWARD &&   ODDLINES) { s = 0.0; t = 1.0; }
    return  (vTileCoord.y-s)*2.  > tileCount.y - t;
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
        col = vec4(1.,1.,1.,1.); 
    }

	if ( newestLine() ) { col.a *= (1.0 - vScroller); }
	if ( oldestLine() ) { col.a *= vScroller; }

    if (!gl_FrontFacing) { col.a /= 2.; }

    gl_FragColor = col;

}

