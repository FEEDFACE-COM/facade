
uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;

bool DEBUG    = vDebugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;


void main() {

    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { col = 
        texture2D(texture, vTexCoord); 
    }
    

    bool firstLine =  0.5*vTileCount.y       == vTileCoord.y ;
    bool lastLine  = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;

    if (DEBUG && (firstLine || lastLine) ) {
        col.rgb = vec3(1.,0.,0.);
    }
    
    gl_FragColor = col;
}
