
uniform float debugFlag;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
