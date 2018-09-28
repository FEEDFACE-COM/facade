
uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vDebugFlag;
varying float vScroller;
varying float vTimer;

bool DEBUG = vDebugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;


void main() {
    float scroll = vScroller;
    
    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;

    vec4 col = texture2D(texture, tex);
    if (DEBUG) {
//        col = vec4(1.0,1.0,1.0,1.0);
        col = vec4(1.0,scroll,scroll,1.0);
    } 
    
    col.r = cos(vTimer + vTileCoord.x) * col.r;
    col.gb = cos(vTimer + vTileCoord.y) * col.gb;
    
    gl_FragColor = vec4(col.rgb,1.0);
}
