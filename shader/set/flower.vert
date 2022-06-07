
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float wordCount;
uniform float wordMaxWidth;

uniform float wordIndex;
uniform float wordWidth;
uniform float wordFader;
uniform float wordValue;

uniform float charCount;

uniform float screenRatio;
uniform float fontRatio;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute float charIndex;
attribute float charOffset;

varying vec4 vTexCoord;
varying vec4 vPosition;
varying float vCharIndex;

bool DEBUG = debugFlag > 0.0;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }
float Ease(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }


vec3 rotate(float w, vec3 v) {
    mat3 M = mat3(
       cos(w), sin(w), 0.0,
      -sin(w), cos(w), 0.0,
          0.0,    0.0, 0.0
    );
    return M*v;
}

vec3 translate(float w, vec3 v) {
    float r = 3.0;
    v.x += cos(w)*r;
    v.y += sin(w)*r;
    return v;
}

vec3 curve(vec3 v) {
    v.z += .25 * sin(now + charIndex + wordIndex);
    return v;   
}

float len(vec3 v) {
    return sqrt( v.x*v.x + v.y*v.y + v.z*v.z );
}

vec3 fun(float phi, vec3 v) {

    float ARC = TAU;
    float sector = ARC / wordCount;
    float gamma = wordIndex * sector;
//    gamma += phi;

    float R0 = 0.5 * 10.;
    float R1 = 0.1;

    float r0 = R0;
    float r1 = R0+R1;


    vec3 p0,p1;
    
    p0.x = cos(gamma)*r0;
    p0.y = sin(gamma)*r0;

    p1.x = cos(gamma)*r1;
    p1.y = cos(gamma)*r1;
    
    vec3 n0;

    vec3 e2 = vec3(0.,0.,1.);
    n0 = cross(p0,e2);
    
    vec3 q0,q1;
    
    q0 = p0+n0;
    q1 = p0-n0;
    

    float x = v.x / (wordWidth * fontRatio);
    float y = v.y / 1.;


    v.x += p0.x;
    v.y += p0.y;
    
//    v.x += x * (p1-p0).x;
//    v.y += y * (q0-q1).y;
    

    
    

    
//    v += p0;
    
//    v.x = x * (p1 - p0).x;
//    v.y = y * (q1 - q0).y;


//    float alpha = sector;
    
    
//    v.x = pX.x * x;
//    v.y = pX.y * y;

//    {
//        float a1 = len(v) * cos(gamma);
//
//        vec3 b = 1./len(p0) * p0;
//        
//        v = a1*b;
//
//        
//    }
    
//    v.x *= x;
//    v.y *= y;
    
//    v.x = x * pX.x;
//    v.y = y * pX.y;
    
    
//    float r = r0 + x*(r1-r0);
//    float a = alpha * 1.*y;
//
//    v.x = cos(gamma+a)*r;
//    v.y = sin(gamma+a)*r;

    return v;
}

vec3 rot2(float w, vec3 v) {

    mat3 M = mat3(
       cos(w), 0.0, sin(w),
          0.0, 1.0,    0.0,
      -sin(w), 0.0, cos(w)
    );
    
    return M*v;
}

vec3 rot3(float w, vec3 v) {
    mat3 M = mat3(
          1.0,    0.0,    0.0,
          0.0, cos(w), sin(w),
          0.0,-sin(w), cos(w)
    );
    return M*v;
}


void main() {
    float fader = wordFader;
    
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;

    pos.y += .5;
    pos.x += wordWidth/2.;

    float phi = -now/4.;
    phi += .5 * sin( Ease( now ) );
    phi += .25 * cos( Ease(now) );

    float ARC = TAU;
    float sector = ARC / wordCount;
    float gamma = wordIndex * sector;

    pos.xyz = rotate(gamma+phi,pos.xyz);
    pos.xyz = curve(pos.xyz);
    pos.xyz = translate(gamma+phi,pos.xyz);
    
    float rho = now/1.;

//    pos.xyz = fun(phi, pos.xyz);
    
    pos.xyz = rot2( Ease(rho)*PI/8., pos.xyz);
     pos.xyz = rot3( -PI/8. + Ease(rho+PI/3.)*-PI/8., pos.xyz);
    

    gl_Position = projection * view * model * pos;
}


