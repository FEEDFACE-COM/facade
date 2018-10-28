
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

    float MARGIN = 0.5;

    vec2 pos = vTexCoord;
    vec3 col = vec3(1.0,1.0,1.0);

//    if ( grid(pos) ) { col = vec4(1.,1.,1.,0.5); }
//    if ( pos.y > 0.0 && pos.y < 1.0 && abs(pos.x) <= w ) { col = vec4(0.,1.,0.); }
//    if ( pos.x > 0.0 && pos.x < 1.0 && abs(pos.y) <= w ) { col = vec3(1.,0.,0.); }
  
  

    if (pos.x > MARGIN || pos.x < -1. * MARGIN  ) {
        col -= 0.5;
    }
    
    if (pos.y > MARGIN || pos.y < -1. * MARGIN  ) {
        col -= 0.5;
    }
    
       
    gl_FragColor = vec4(col.rgb, 1.0);
}


