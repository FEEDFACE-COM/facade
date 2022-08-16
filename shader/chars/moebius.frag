
uniform float charCount;
uniform float charLast;

uniform float scroller;

uniform sampler2D texture;        
uniform float now;                // time since start of program, as seconds
uniform float debugFlag;          // 0.0 unless -D flag given by user


varying vec4 vTexCoord;
varying vec4 vPosition;
varying float vCharIndex;

bool DEBUG = debugFlag > 0.0;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

vec3    BLACK=vec3(0.,0.,0.);   vec3 WHITE=vec3(1.,1.,1.);
vec3 DARKGRAY=vec3(.25,.25,.25); vec3 GRAY=vec3(.5,.5,.5); vec3 LIGHTGRAY=vec3(.75,.75,.75);
vec3      RED=vec3(1.,0.,0.);   vec3 GREEN=vec3(0.,1.,0.);     vec3 BLUE=vec3(0.,0.,1.);
vec3     CYAN=vec3(0.,1.,1.); vec3 MAGENTA=vec3(1.,0.,1.);   vec3 YELLOW=vec3(1.,1.,0.);


void main() {
    vec4 col;
    vec4 pos = vPosition;
    vec4 tex = vTexCoord;

    col = texture2DProj(texture, vTexCoord);


    float count = floor(charCount/2.);
    float index = vCharIndex;
    if (index >= count) {
        index -= count;
    }

    if (gl_FrontFacing) { 

        if (vCharIndex == 0.0 ) {
            col.a *= 1. - scroller;
        }
        else if (vCharIndex == charCount-2.) {
            col.a *= scroller;
        } 
        else if (vCharIndex == charCount-1.) {
            col.a = 0.;
        }

        
    } else { // hide back side
        col.a = 0.;
    }


    if (DEBUG) { 
        col.rgb = WHITE;
    }

    gl_FragColor = col;
    
}


