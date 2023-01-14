
uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;
uniform vec2 cursorPos;



varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

//    col.rgb = vec3(0.,0.,0.);
    col.rgb /= 2.;
    col.a = 1.;
    
//    if ( vTileCoord.x == 0.0 || vTileCoord.y == 0.0 ) {
//        
//        col.rgb += vec3(1.,1.,1.);
//    
//    } 
    if ( vTileCoord.x + tileCount.x/2. <= .5 ) {
    
        col.rgb += vec3(0.,1.,0.);
    
    }
    if ( vTileCoord.x + tileCount.x/2. >= tileCount.x-2. ) {
    
        col.rgb += vec3(1.,0.,1.);
    
    }
    if ( vTileCoord.y + tileCount.y/2. <= .5 ) {
    
//        col.rgb += vec3(1.,0.,0.);
        col.r += 1.0;

    }
    if ( vTileCoord.y + tileCount.y/2. >= tileCount.y-1. ) {
    
        col.rgb += vec3(0.,1.,1.);
    
    }
    if ( mod(-abs(vTileCoord.x) , 2.) == 0.0  ^^ mod(-abs(vTileCoord.y) , 2.) == 0.0 ) {
    
        col.rgb += 0.5 * vec3(1.,1.,1.);
    
    }
    
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) {
        col.rgba = 1. - col.rgba;
    }

    if (!gl_FrontFacing) { 
        col.a /= 4.;
    }

    gl_FragColor = col;
    
}
