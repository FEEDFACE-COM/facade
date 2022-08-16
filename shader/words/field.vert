uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float wordCount;
uniform float wordMaxWidth;
uniform float wordMaxLength;

uniform float wordFader;
uniform float wordAge;

uniform float screenRatio;
uniform float fontRatio;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute float charIndex;
attribute float wordIndex;
attribute float wordLength;
attribute float wordWidth;

varying vec4 vTexCoord;
varying vec4 vPosition;
varying float vCharIndex;
varying float vWordIndex;
varying float vWordWidth;

bool DEBUG = debugFlag > 0.0;


float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }
float Ease(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

mat3 rotx(float w) {return mat3(1.0,0.0,0.0,0.0,cos(w),sin(w),0.0,-sin(w),cos(w));}
mat3 roty(float w) {return mat3(cos(w),0.0,sin(w),0.0,1.0,0.0,-sin(w),0.0,cos(w));}
mat3 rotz(float w) {return mat3(cos(w),sin(w),0.0,-sin(w),cos(w),0.0,0.0,0.0,1.0);}

vec3 translate(float w, float r, vec3 v) {
    v.x += cos(w)*r;
    v.y += sin(w)*r;
    return v;
}


void main() {
    float fader = wordFader;
    
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;
    vWordIndex = wordIndex;
    vWordWidth = wordWidth;


    float r0 = 2.;
    float r1 = r0 * screenRatio;
    float sector = TAU / wordCount;

    float gamma =  PI/2. + (sector) * -wordIndex;


    pos.x += cos(gamma) * r1;
    pos.y += sin(gamma) * r0;
    pos.z = 0.0;

    pos.xy /= 4.;
    pos.z += 2.*wordAge;

    float zoom = 1.0;
    {
        zoom = 1./ (r0+log(wordMaxWidth));
    }
    pos.xy *= zoom;


    gl_Position = projection * view * model * pos;
}

