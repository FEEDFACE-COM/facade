
uniform float debugFlag;
uniform float ratio;

varying vec2 vTexCoord;


bool DEBUG = debugFlag > 0.0;
float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.5) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}


void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    float a = 1.0;


    vec2 pos = vTexCoord;


    if ( true && grid(pos) ) {
        float gray = 0.5;
        col = gray * vec3(1.,1.,1.);
    }

    
    
    gl_FragColor = vec4(col.rgb, a);
}


