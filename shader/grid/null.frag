
uniform float debugFlag;
uniform sampler2D texture;

varying vec2 vTexCoord;


bool DEBUG    = debugFlag > 0.0;

void main() {
    if (DEBUG) { 
        gl_FragColor = vec4(1.,1.,1.,1.); 
    } else { 
        gl_FragColor = texture2D(texture, vTexCoord); 
    }
}
