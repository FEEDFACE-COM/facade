
uniform float debugFlag;
uniform float screenRatio;

varying vec2 vTexCoord;


bool DEBUG = debugFlag > 0.0;
float w = 0.002;

bool major(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=1.0) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}

bool minor(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.25) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}



void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    float a = 1.0;


    vec2 pos = vTexCoord;


    if ( true && major(pos) ) {
        float gray = 1.0;
        col = gray * vec3(1.,1.,1.);
    } else if ( true && minor(pos) ) {
        float gray = 0.5;
        col = gray * vec3(1.,1.,1.);
    }

    gl_FragColor = vec4(col.rgb, a);
}



