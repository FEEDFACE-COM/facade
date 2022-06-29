
uniform sampler2D texture;        // font glyph atlas texture

uniform vec2 tileCount;           // grid dimensions, as (columns,rows)
uniform float scroller;           // scroller amount, from 0.0 to 1.0
uniform vec2 cursorPos;           // terminal cursor position, in (column,row)

uniform float now;                // time since start of program, as seconds
uniform float debugFlag;          // 0.0 unless -D flag given by user

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;


bool DEBUG = debugFlag > 0.0;

void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 
    
    col.rgb /= 4.;

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    if ( vGridCoord.y == 0.0 ) { // oldest line
        col.a *= (1.-abs(scroller));
    }
    
    if ( vGridCoord.y == tileCount.y ) { // newest line
        col.a *= abs(scroller);
    }
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) { // invert cursor
        col.rgba = 1. - col.rgba;
    }
    
    {
        float ARC = TAU;
        float sector = ARC / (tileCount.y+4.);
        float gamma = (tileCount.y-vGridCoord.y) * sector;


        float phi = 0.0;
        phi += now;
        phi += abs(scroller) * sector;
        

        col.rgb /= 4.;
    
        col.rgb += ( .5 + sin(gamma+phi)/2. );
    }
    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
