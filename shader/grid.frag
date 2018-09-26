uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vDebugFlag;

bool DEBUG = vDebugFlag > 0.0;


void main() {

    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;


    vec4 col = texture2D(texture, tex);

       
    if (DEBUG && pos.x == 0.0 && pos.y == 0.0 ) {
        col.r += 0.5;
        col.g += 0.5;
        col.b += 0.5;
    }

    if (DEBUG && pos.x == 0.0  ) {
        col.g += 0.5;
    }

    if (DEBUG && pos.y == 0.0  ) {
        col.r += 0.5;
    }
    
    gl_FragColor = vec4(col.rgb,1.0);
}
