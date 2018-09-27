
uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vDebugFlag;
varying float vScroller;

bool DEBUG = vDebugFlag > 0.0;


void main() {
    float scroll = vScroller;
    
    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;

    vec4 col = texture2D(texture, tex);
    if (DEBUG) {
//        col = vec4(1.0,1.0,1.0,1.0);
        col = vec4(1.0,scroll,scroll,1.0);
    } 
    
    gl_FragColor = vec4(col.rgb,1.0);
}
