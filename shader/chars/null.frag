
uniform sampler2D texture;        
uniform float debugFlag;          // 0.0 unless -D flag given by user

varying vec4 vTexCoord;

bool DEBUG = debugFlag > 0.0;


void main() {
    vec4 col;

    col = texture2DProj(texture, vTexCoord);
    
    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    }

    if (!gl_FrontFacing) { col.a /= 4.; }


    gl_FragColor = col;
    
}


