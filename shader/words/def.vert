
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float wordCount;
uniform float wordMaxWidth;
uniform float wordMaxLength;

uniform float wordFader;
uniform float wordValue;

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


void main() {
    float fader = wordFader;
    
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;
    vWordIndex = wordIndex;
    vWordWidth = wordWidth;


    pos.y += 0.5;
    pos.x += wordWidth/2.;
    

    float rows,cols;
    float slots = wordCount;
    float row,col;

    float SPACER = 1.;
    float colWidth = (wordMaxLength+SPACER)*fontRatio;
    
    
    colWidth = wordMaxWidth;
    colWidth += SPACER * fontRatio;

    float wordRatio = wordMaxWidth / 1.;
    float ratio = wordRatio / screenRatio;
    float a = sqrt( wordCount );
    float b = floor(a/ratio);

    if (wordCount == 1.0) {
        cols = 1.;
        rows = 1.;
        col = 0.;
        row = 0.;
    } else if (wordCount == 2.0) {
        cols = 2.;
        rows = 1.;
        col = wordIndex;
        row = 0.;
    } else {
        if (wordCount <= 8.) {
            cols = 2.;
        } else if (wordCount <= 24.) {
            cols = 3.;
        } else if (wordCount <= 48.) {
            cols = 4.;
        } else {
            cols = floor( sqrt(wordCount) / 1.6);
        } 
    
        rows = ceil( wordCount / cols );
        row = mod(wordIndex+1., rows) -1.;
        row = rows - row - 1.; // switch top-down
        col = floor((wordIndex+1.) / rows);

    }

    pos.x += col * colWidth;
    pos.y += row;

    pos.x -= (cols/2.) * colWidth;
    pos.y -= rows/2.;

//    //DEBUG
//    if (wordIndex == floor(0)) { // 2
//        pos.z+=100.;
//    }
//    if (wordIndex == floor(1)) { // 0
//       pos.z+=100.;
//    }


    float zoom = 1.0;
    {
        float zr = 2./rows;
        float zc = 2./(cols * colWidth) * screenRatio;
        zoom = min(zr,zc);
    }
    pos.x *= zoom;
    pos.y *= zoom;
    
    
    gl_Position = projection * view * model * pos;
}

