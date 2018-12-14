
uniform float debugFlag;
uniform float ratio;

varying vec2 vTexCoord;

bool DEBUG = debugFlag > 0.0;

bool grid(vec2 pos) {
    float w = 0.005;
    for (float d = -2.0; d<=2.0; d+=1.0) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}


//float mask(vec2 pos) {
//    float MARGIN = 0.125;
//    if ( abs(pos.x) > (ratio*1.0) - MARGIN ) { return 1.0; }
//    if ( abs(pos.y) >        1.0  - MARGIN ) { return 1.0; }
//    return 0.0;
//}

float MAX(float a, float b) { 
    if (a>=b) 
        return a; 
    else 
        return b; 
}

float mask(vec2 p) { 
    float x = p.x; float y = p.y;
    float ff = 1.;
    float A = 3. * ff;
    float B = 4. * ff; 
    if (abs(x) >= A/B && abs(y) >= A/B)
        return 1.0 - ( B * MAX(abs(x),abs(y)) - A );
    else if (abs(x) >= A/B) 
        return 1.0 - ( B * abs(x) - A);
    else if (abs(y) >= A/B)
        return 1.0 - ( B * abs(y) - A);
    return 1.0;
}



void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    float a = 1.0;


    vec2 pos = vTexCoord;

    a = mask(vec2(pos.x/ratio,pos.y));


    if ( true && DEBUG && grid(pos) ) {
        float gray = 0.5;
        col = gray * vec3(1.,1.,1.);
    }

    if (DEBUG) {
        col.r =    1. - mask(vec2(pos.x/ratio,pos.y)) ;
    }
    
    
    gl_FragColor = vec4(col.rgb, a);
}
