
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
bool DEBUG_FREEZE = false;

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


/*************************************/
/*  mat3 rotz(float w) {             */
/*      return mat3(                 */
/*         cos(w), sin(w), 0.0,      */
/*        -sin(w), cos(w), 0.0,      */
/*            0.0,    0.0, 0.0       */
/*      );                           */
/*  }                                */
/*  mat3 roty(float w) {             */
/*      return mat3(                 */
/*         cos(w), 0.0, sin(w),      */
/*            0.0, 1.0,    0.0,      */
/*        -sin(w), 0.0, cos(w)       */
/*      );                           */
/*  }                                */
/*  mat3 rotx(float w) {             */
/*      return mat3(                 */
/*            1.0,    0.0,    0.0,   */
/*            0.0, cos(w), sin(w),   */
/*            0.0,-sin(w), cos(w)    */
/*      );                           */
/*  }                                */
/*************************************/


float len(vec3 v) {
    return sqrt( v.x*v.x + v.y*v.y + v.z*v.z );
}



vec3 rotate(float w, vec3 v) {
    return rotz(w)*v;
}

vec3 translate(float w, vec3 v) {
    float r = 3.0;
    v.x += cos(w)*r;
    v.y += sin(w)*r;
    return v;
}

vec3 wave(vec3 v) {
    float run = 0.0;
    run = 4.*now/1.;
    v.z += .25 * sin(run + v.x + 1.);
    return v;
}

vec3 curve(vec3 v,float x) {
    float run = 0.0;
    run = 4.*now/1.;
    v.z += 4.*x + .25 * -cos(x*PI);
    v.z += .25 * sin(run + x*PI + 2.*PI*(wordIndex/wordCount));
    return v;   
}




void main() {
    float fader = wordFader;
    
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;
    

    pos.y += .5;
    pos.x += wordWidth/2.;

    float phi = 0.0;
    if (! DEBUG_FREEZE ) {
        phi += -now/4.;
        phi += .5 * sin( Ease( now ) );
        phi += .25 * cos( Ease(now) );
    }

    float ARC = TAU;
    float sector = ARC / wordCount;
    float gamma = wordIndex * sector;

    float XXX = pos.x / wordWidth;
    
    pos.xyz = rotate(gamma+phi,pos.xyz);
    pos.xyz = curve(pos.xyz,XXX);
    pos.xyz = translate(gamma+phi,pos.xyz);
    

    float rho = 0.0;
    if (! DEBUG_FREEZE ) {
        rho = now/5.;
    }
    pos.xyz = roty( -PI/8. + Ease(rho)*PI/8.) * pos.xyz;
    pos.xyz = rotx( -PI/4. + Ease(rho+PI/3.)*-PI/8.) * pos.xyz;


    gl_Position = projection * view * model * pos;
}


