
uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;

bool DEBUG    = vDebugFlag > 0.0;

void main() {
    if (DEBUG) { 
        gl_FragColor = vec4(1.,1.,1.,1.); 
    } else { 
        gl_FragColor = texture2D(texture, vTexCoord); 
    }
}
