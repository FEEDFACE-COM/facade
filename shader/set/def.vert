
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

uniform float ratio;

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
    pos.x += charOffset;

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;

    bool odd = mod(wordIndex,2.0) == 1.0;
    
    if (odd) {
//        pos.x += 2.43/2.;
        pos.x += wordWidth/2.;
    } else {
//        pos.x -= 2.43/2.;
        pos.x -= wordWidth/2.;
    }
    pos.y -= (wordIndex - wordCount/2.)/2.;
    

        
    gl_Position = projection * view * model * pos;
}

