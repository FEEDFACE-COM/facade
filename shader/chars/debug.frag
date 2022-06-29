
uniform sampler2D texture;        
uniform float debugFlag;          // 0.0 unless -D flag given by user

bool DEBUG = debugFlag > 0.0;

void main() {
    vec4 col;
    gl_FragColor = col;    
}


