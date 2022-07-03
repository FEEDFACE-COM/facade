
uniform float charCount;
uniform float charLast;

uniform float scroller;


uniform sampler2D texture;        
uniform float debugFlag;          // 0.0 unless -D flag given by user

varying vec4 vTexCoord;
varying vec4 vPosition;
varying float vCharIndex;
varying float vCharOffset;

bool DEBUG = debugFlag > 0.0;

vec3    BLACK=vec3(0.,0.,0.);   vec3 WHITE=vec3(1.,1.,1.);
vec3 DARKGRAY=vec3(.25,.25,.25); vec3 GRAY=vec3(.5,.5,.5); vec3 LIGHTGRAY=vec3(.75,.75,.75);
vec3      RED=vec3(1.,0.,0.);   vec3 GREEN=vec3(0.,1.,0.);     vec3 BLUE=vec3(0.,0.,1.);
vec3     CYAN=vec3(0.,1.,1.); vec3 MAGENTA=vec3(1.,0.,1.);   vec3 YELLOW=vec3(1.,1.,0.);

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float EaseOut(float x) { return        sin( x * PI/2. ); }
float EaseIn(float x)  { return   1. - cos( x * PI/2. ); }


void main() {
    vec4 col;
    vec4 pos = vPosition;

    col = texture2DProj(texture, vTexCoord);
    
    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    }
    
    if (charCount >= 4.) {

        if ( vCharIndex == 0.0 ) {
            col.a *= (1. - EaseIn( .5 + scroller/2.) );
        } else if ( vCharIndex == 1.0 ) {
            col.a *= (1. - EaseIn( scroller/2. ) );
        } else if (vCharIndex == charCount-2.) {
            col.a *= EaseIn( .5 + scroller/2.);
        } else if (vCharIndex == charCount-1.) {
            col.a *= EaseIn( scroller/2. );
        }


    } else {
        if ( vCharIndex == 0.0 ) {
            col.a *= (1. - EaseIn(scroller) );
        } else if (vCharIndex >= charCount-1.) {
            col.a *= EaseIn(scroller);
        }
    }

    if (!gl_FrontFacing) { col.a /= 4.; }


    gl_FragColor = col;
    
}


