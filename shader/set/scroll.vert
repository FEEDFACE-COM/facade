
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
    vTexCoord.xy = texCoord.xy;
    
    vec4 pos,tex;
    pos = vec4(vertex,1);
    tex = vec4(texCoord.xy, 1., 1.);

    pos.y -= wordCount/2.;
    pos.y += (wordIndex/wordCount) * wordCount;
    pos.y -= 0.5;
    
    pos.x += (wordCount/2.) * ratio ;
    pos.x -= wordFader * ( wordCount * ratio);
    
    vPosition = pos;
    vTexCoord = tex;
    vCharIndex = charIndex;
    
    gl_Position = projection * view * model * pos;
}

