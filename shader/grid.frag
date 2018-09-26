uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying vec2 vDebugFlag;


bool debug = false;

void main() {

    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;


    vec4 col = texture2D(texture, tex);




    bool debug = false; 
    
    if ( vDebugFlag.x > 0.0 ) {
        debug = true;
    }
       
    if (debug && pos.x == 0.0 && pos.y == 0.0 ) {
        col.g += 0.5;
    }

    if (debug && pos.x == 1.0 && pos.y == 1.0 ) {
        col.b += 0.5;
    }
    
    gl_FragColor = vec4(col.rgb,1.0);
}
