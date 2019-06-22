
uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;
uniform vec2 cursorPos;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying float vScroller;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    if ( vGridCoord.y == 0.0 ) { // oldest line
        col.a *= (1.-vScroller);
    }
    
    if ( vGridCoord.y == tileCount.y ) { // newest line
        col.a *= vScroller;    
    }
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) { // invert cursor
        col.rgba = 1. - col.rgba;
    }

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
