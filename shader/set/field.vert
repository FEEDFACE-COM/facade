
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float wordMax;
uniform float wordIndex;

uniform float wordWidth;
uniform float wordFader;
uniform float wordCount;

uniform float ratio;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec4 vTexCoord;
varying vec4 vPosition;


bool DEBUG = debugFlag > 0.0;

mat4 scaleMatrix(float a) {
    return mat4(
         a, 0., 0., 0.,
        0.,  a, 0., 0.,
        0., 0., 1., 0.,
        0., 0., 0., 1.
    );
}


float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }


void main() {
    float fader = wordFader;
    
    vec4 pos;
    pos = vec4(vertex,1);

    float a;

    pos.x *= wordMax/10.;
    pos.y *= wordMax/10.;
    
    float rx = 0.5 * wordMax/2.;
    float ry = wordMax/2.;
    
    float idx = wordIndex/wordMax;
    
    float w = 2.*TAU;
    float o = PI/2. + PI/4.;
    a = idx * w + o;
    pos.x += rx * cos(a);
    pos.y += ry * sin(a);
        

    if (idx >= 0.5 ) {
        pos.x += wordMax/2.;
    } else {
        pos.x -= wordMax/2.;
    }


    pos.z -= wordMax; 
    pos.z += fader * 1.5 * wordMax;

    mat4 v = view;
//    v = view * scaleMatrix( 1.  );

    vPosition = pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    
        
    gl_Position = projection * v * model * pos;
}

