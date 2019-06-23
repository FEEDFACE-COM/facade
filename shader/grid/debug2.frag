
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

    
    
    if ( 
        vTileCoord.x == 0.0
     || vTileCoord.y == 0.0
     || vTileCoord.x+1. >= (tileCount.x/2.)
     || vTileCoord.y+1. >= (tileCount.y/2.)
     || vTileCoord.x <= -(tileCount.x/2.)
     || vTileCoord.y <= -(tileCount.y/2.)
    ) {
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
     
    } else if (
         mod(-abs(vTileCoord.x) , 2.) == 0.0 
     ^^ mod(-abs(vTileCoord.y) , 2.) == 0.0
    
    ) {
        col.rgb = 0.75 * vec3(1.,1.,1.);
        col.a = 0.5;
    } else {
        col = vec4(0.);
    }
         


    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) {
        col.rgba = 1. - col.rgba;
    }

    if (gl_FrontFacing) { 
        col.rgb /= 2.;
    }

    gl_FragColor = col;
    
}
