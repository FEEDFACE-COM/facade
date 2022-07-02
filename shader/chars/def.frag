
uniform sampler2D texture;        
uniform float debugFlag;          // 0.0 unless -D flag given by user

varying vec4 vTexCoord;
varying vec4 vPosition;
varying float vCharIndex;
varying float vCharOffset;

bool DEBUG = debugFlag > 0.0;


void main() {
    vec4 col;
    vec4 pos = vPosition;
    vec4 tex = vTexCoord;

    col = texture2DProj(texture, tex);
//    col.a *= wordFader;

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
    }

    if (!gl_FrontFacing) { col.a /= 4.; }


    gl_FragColor = col;
    
}


