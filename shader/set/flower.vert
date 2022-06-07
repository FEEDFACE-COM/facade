
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float wordCount;
uniform float wordIndex;

uniform float wordWidth;
uniform float wordFader;
uniform float wordValue;

uniform float charCount;

uniform float screenRatio;

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


vec2 fun(float gamma, vec2 v) {

    float R0 = 4.;
    float R1 = R0 + 4.;

    float r0 = R0;
    float r1 = R1;

    float ARC = TAU;
    float sector = ARC / (wordCount*1.);
    

    float alpha = sector;
    
    float x = v.x / wordWidth;
    float y = v.y / 2.;
    
    float r = r0 + x*(r1-r0);
    float a = alpha * y;

    v.xy = vec2( cos(gamma+a)*r, sin(gamma+a)*r);

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

    float rho = now/8.;

    float ARC = TAU;
    float sector = ARC / wordCount;
    float gamma = wordIndex * sector;

//    pos.xyz = rotate(gamma+phi,pos.xyz);
//    pos.xyz = translate(gamma+phi,pos.xyz);

    pos.xy = fun(gamma+phi, pos.xy);
    
//    pos.xyz = rot2( Ease(rho)*PI/8., pos.xyz);
//    pos.xyz = rot3( -PI/8. + Ease(rho+PI/3.)*-PI/8., pos.xyz);
    

    gl_Position = projection * view * model * pos;
}


