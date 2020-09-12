
uniform sampler2D texture;        // font glyph atlas texture

uniform vec2 tileCount;           // grid dimensions, as (columns,rows)
uniform float scroller;           // scroller amount, from 0.0 to 1.0
uniform vec2 cursorPos;           // terminal cursor position, in (column,row)

uniform float debugFlag;          // 0.0 unless -D flag given by user

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;



bool DEBUG = debugFlag > 0.0;

void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    if ( vGridCoord.y == 0.0 ) { // oldest line
        col.a *= (1.-abs(Scroller));
    }
    
    if ( vGridCoord.y == tileCount.y ) { // newest line
        col.a *= abs(scroller);
    }
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) { // invert cursor
        col.rgba = 1. - col.rgba;
    }

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
