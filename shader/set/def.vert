
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


void main() {
    float fader = wordFader;
    
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;

//    float ratio = wordIndex/wordCount;
//    
//    float MAX_WORD_LENGTH = 12;
//    float ratio_width = wordLength / MAX_WORD_LENGTH;
//
//    pos.y += wordIndex/2. - wordCount/4.;    
//    
//    if (wordIndex >= wordCount/2.) {
//        pos.x += 3.;
//    } else {
//        pos.x -= 3.;
//    }
//
//
//    if (wordIndex >= wordCount/2.) {
//        pos.y += wordIndex-wordCount/2.;
//    }


//    pos.z += wordIndex;


    float WORD_LENGTH_MAX = 16.0;
    
    float rows = 8.0; // standard size
    float cols = 2.0;

    float row = mod(wordIndex, rows/cols);
    float col = 0.0;
    if (wordIndex >= wordCount/2.) {
        col = 1.0;
    }

    pos.y += 3.;
    pos.x += wordWidth/2.;
    pos.x -= 4. * screenRatio;


    pos.x += col * 4. * screenRatio;

    pos.y -= row * (rows/ (wordCount/cols));


 
    gl_Position = projection * view * model * pos;
}

