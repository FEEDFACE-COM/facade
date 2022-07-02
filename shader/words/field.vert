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
attribute float charOffset;
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
mat3 rotz(float w) {return mat3(cos(w),sin(w),0.0,-sin(w),cos(w),0.0,0.0,0.0,0.0);}

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


    float radius = 8.;
    float gamma = 0.;
    float ARC = TAU * 2./6.;
    float sector = (2.*ARC) / wordCount;



    if ( wordIndex <= wordCount/2. ) {
 
        float index = wordIndex;
 
        gamma =  PI/2. + ARC/4. + sector * index;

 
        
    } else {

        float index = wordIndex;
 
        gamma =  PI + ARC/4. + sector * index;
        
    }

    pos.x += cos(gamma) * radius;
    pos.y += sin(gamma) * radius;


    pos.xy /= 4.;

    float zoom = 1.0;
    {
        zoom = 1./ (radius+log(wordMaxWidth));
    }
    pos.xy *= zoom;

    pos.z += 2.*wordAge;


    gl_Position = projection * view * model * pos;
}

