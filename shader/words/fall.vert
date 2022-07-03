
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

void main() {
    float fader = wordFader;
    
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;
    vWordIndex = wordIndex;
    vWordWidth = wordWidth;


    float s = 1.0;
    if (mod(wordIndex,2.) == 1.0) {
        s = -1.0;
    }


//    pos.y += .5*cos( s * (pos.x/(wordMaxWidth/2.)*PI/2.) + wordAge * PI/2.);


//
    pos.y -= .25;






    
    if (wordCount > 1.0) { // distribute across screen width
        float d = (wordCount) * fontRatio;
        pos.x -= d/2.;
        pos.x += (wordIndex/(wordCount-1.)) * d;
    }


    float zoom = 1.0;
    zoom = 2./(wordMaxWidth + wordCount*fontRatio) * screenRatio;
    pos.xy *= zoom;
    
    if (wordAge < 0.) {

        pos.y -= -cos( (wordIndex/wordCount) * (TAU * wordCount*wordMaxWidth ));
        
    } else {
        
        
        pos.y += 1.;
        if (false) {
            pos.y -= 2. * wordAge;
        } else {
            pos.y -= 2. * (sin(wordAge * PI/2. ));
            pos.x += .5 * s * sin(wordAge * PI);
        }

    }




    
    gl_Position = projection * view * model * pos;
}


