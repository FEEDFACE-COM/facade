
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

    if (true) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.;
    } 

    float F = 1.;

    float x = vGridCoord.x / tileCount.x;
    float y = vGridCoord.y / tileCount.y;
    


    
    if (gl_FrontFacing) {
        col.rgb = vec3(1.,0.,0.);
        col.a = 1.;
//        col.r *= F * (1. - x);
    } else {
        col.rgb = vec3(0.0,1.0,0.0);
//        col.a = .25;
//        col.g *= F * (1. - y);
    }

    if ( abs(vGridCoord.y) == tileCount.y  ) {
    }
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) {
        col.rgba = 1. - col.rgba;
    }


    gl_FragColor = col;
    
}
