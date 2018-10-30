
uniform float debugFlag;

bool DEBUG = debugFlag > 0.0;


void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    float a = 1.0;

    if ( DEBUG ) {
        col.r = 0.5; 
    }
    
    gl_FragColor = vec4(col.rgb, a);
}
