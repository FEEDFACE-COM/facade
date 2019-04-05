
uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying float vScroller;

bool DEBUG    = debugFlag > 0.0;


bool oddColCount() { return mod(tileCount.x, 2.0) == 1.0 ; }
bool oddRowCount() { return mod(tileCount.y, 2.0) == 1.0 ; }


bool firstLine() {
    float s = 0.0;
    float t = 0.0;
    if      ( ! (downward == 1.0) && ! oddRowCount()) { s = 0.0; t = 0.0; }
    else if ( ! (downward == 1.0) &&   oddRowCount()) { s = 1.0; t = 0.0; }
    else if (   (downward == 1.0) && ! oddRowCount()) { s = 0.0; t = 1.0; }
    else if (   (downward == 1.0) &&   oddRowCount()) { s = 0.0; t = 1.0; }
    return (  tileCount.y + (vTileCoord.y-s) - 1.) * 2. + t <= tileCount.y;
}


bool lastLine() { 
    float s = 0.0;
    float t = 0.0;
    if      ( ! (downward == 1.0) && ! oddRowCount()) { s = 0.0; t = 0.0; }
    else if ( ! (downward == 1.0) &&   oddRowCount()) { s = 1.0; t = 0.0; }
    else if (   (downward == 1.0) && ! oddRowCount()) { s = 0.0; t = 1.0; }
    else if (   (downward == 1.0) &&   oddRowCount()) { s = 0.0; t = 1.0; }
    return  (vTileCoord.y-s)*2.  > tileCount.y - t;
}

bool newestLine() {
    return  ( !(downward == 1.0) && firstLine() ) || ( (downward == 1.0) && lastLine() ) ;
}

bool oldestLine() {
    return ( !(downward == 1.0) && lastLine() ) || ( (downward == 1.0) && firstLine() );
}



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 
    
//    if (firstLine()) { col.r *= 0.; } //cyan
//    if (lastLine())  { col.g *= 0.; } //magenta
    
    
    if ( newestLine() )   { col.rb *= 0.; col.a *= (1.-vScroller); } //green
    if ( oldestLine() ) { col.gb *= 0.; col.a *= vScroller; }      //red

    if (!gl_FrontFacing) { col.a /= 2.; }

    gl_FragColor = col;
    
}
