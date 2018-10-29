
varying vec2 vTexCoord;
varying float vDebugFlag;
varying float vRatio;


bool DEBUG = vDebugFlag > 0.0;
float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.5) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}


void main() {

    float MARGIN = 0.25;

    vec2 pos = vTexCoord;
    vec3 col = vec3(0.,1.,1.);
    

    if (pos.x > MARGIN || pos.x < -1. * MARGIN  ) {
//        col -= 0.5;
    }
    
    if (pos.y > MARGIN || pos.y < -1. * MARGIN  ) {
//        col -= 0.5;
    }
    

    float a = vTexCoord.x * -1.;           
    gl_FragColor = vec4( col.rgb, a );
}


