
uniform float now;
uniform float debugFlag;
uniform sampler2D texture;

uniform float wordMax;
uniform float wordMaxWidth;
uniform float wordFader;
uniform float wordIndex;
uniform float wordCount;

varying vec4 vPosition;
varying vec4 vTexCoord;



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


vec3 color(vec4 pos, bool front) {
    vec3 ret;
    vec3 one = MAGENTA;
    vec3 two = GREEN;
    if (mod(wordIndex,2.)==1.) {
        one = CYAN;
        two = RED;
    }
    
    if (front) {
        return one;
    }
    return two;        
}

bool check() {
    float CHECK = 16.;
    if (wordIndex == CHECK || wordIndex+1. == CHECK ) {
        return true;
    }
    return false; 
}

void main() {
    vec4 col;
    vec4 pos = vPosition;
    vec4 tex = vTexCoord;


    if (mod(wordIndex,2.)==1.) {


    }    
    
    if (DEBUG) {
        return;
        col.a = 1.0;
        col.rgb = WHITE;
    }

    if (!DEBUG) {
        col.a = 1.0;
        col.rgb += color(vPosition,gl_FrontFacing);
    
    }
    

    
    if (check()) {
        col.rgb = WHITE;
    }
        
    
    if (!check() && !gl_FrontFacing) {
        col.rgb = WHITE;
        col.a = 0.5;
    }
    
    if (gl_FrontFacing) {
        if  (pos.x < 0.) {
            col.b = 0.;    
        } else {
            col.b = 1.;
        }
    }    
    gl_FragColor = col;
    
}
