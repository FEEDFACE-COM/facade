
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
        col.a = 1.0;
    } 

    float F = 1.;

    float x = vGridCoord.x / tileCount.x;
    float y = vGridCoord.y / tileCount.y;
    
    col.r *= F * (1. - x);
    col.g *= F * (1. - y);

    if ( abs(vGridCoord.y) == tileCount.y  ) {
        col.r = 1.0;
        col.g = 1.0;
        col.b = 0.;
    }
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) {
        col.rgba = 1. - col.rgba;
    }

    if (gl_FrontFacing) { 
        vec3 tmp = vec3(col.rgb);
        col.r = tmp.g;
        col.g = tmp.b;
        col.b = tmp.r;
    }

    gl_FragColor = col;
    
}
