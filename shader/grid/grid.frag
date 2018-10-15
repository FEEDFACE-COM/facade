
uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;


varying vec2 vTileCount;
varying float vDownwardFlag;
varying float vDebugFlag;
varying float vScroller;
varying float vTimer;

bool DEBUG    = vDebugFlag > 0.0;
bool downward = vDownwardFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;


void main() {
    float scroll = abs(vScroller);
    
    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;

    vec4 col;
    if (DEBUG) { col = vec4(1.,1.,1.,1.); }
    else       { col = texture2D(texture, tex); }
    

    bool firstLine, lastLine;
    if (downward) {
        firstLine = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
        lastLine =   0.5*vTileCount.y       == vTileCoord.y ;
    } else {
        firstLine =  0.5*vTileCount.y       == vTileCoord.y ;
        lastLine  = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
    }

    if (firstLine && scroll > 0.5) { //oldest line vanishes later
//        col.rgb = col.rgb * (1.- 2.*(scroll-0.5));
    }

    if (lastLine) { //newest line blends in
        col.rgb = col.rgb * scroll;
    }    
    
    gl_FragColor = col;
}
