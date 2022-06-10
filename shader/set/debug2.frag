
uniform float now;
uniform float debugFlag;
uniform sampler2D texture;

uniform float wordCount;
uniform float wordMaxWidth;
uniform float wordFader;
uniform float wordWidth;
uniform float wordIndex;
uniform float wordValue;

uniform float charCount;

varying vec4 vPosition;
varying vec4 vTexCoord;
varying float vCharIndex;
varying vec4 vCoord;


bool DEBUG    = debugFlag > 0.0;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;
vec3 BLACK = vec3(0.,0.,0.); vec3 GRAY = vec3(.5,.5,.5);    vec3 WHITE = vec3(1.,1.,1.);
vec3 RED = vec3(1.,0.,0.);   vec3 GREEN = vec3(0.,1.,0.);   vec3 BLUE = vec3(0.,0.,1.);
vec3 CYAN = vec3(0.,1.,1.);  vec3 MAGENTA = vec3(1.,0.,1.); vec3 YELLOW = vec3(1.,1.,0.);

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return        cos(x*PI/2. + 3.*PI/2. );        }
float EaseIn(float x) {  return  -1. * cos(x*PI/2.            ) + 1.  ; }

void main() {
    vec4 col;
    vec4 pos = vPosition;
    vec4 tex = vTexCoord;

    vec2 grid = vCoord.xy;
    vec2 word = vCoord.zw;

    col.a = 0.5 * wordFader;

    col.rgb = vec3(0.);

    if (word.x <= 0.0) {
        col.r = 1.0;
    } else if (word.x <= 1.0) { 
        col.rb = vec2(1.0);
    } else if (word.x <= 2.0) { 
        col.rgb = vec3(1.);
    } else if (word.x <= 3.0) { 
        col.gb = vec2(1.0);
    } else if (word.x <= 4.0) { 
        col.b = 1.0;
    } else {
        col.rgb = vec3(1.);
    }

    

    
    if (DEBUG) {
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    }
    
    
    
    gl_FragColor = col;
    
}
