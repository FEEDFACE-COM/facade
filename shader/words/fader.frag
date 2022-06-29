
uniform float now;
uniform float debugFlag;
uniform sampler2D texture;

uniform float wordMax;
uniform float wordMaxWidth;
uniform float wordTimer;
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


float Log(float x) { return log(x); }


void main() {
    vec4 col;
    vec4 pos = vPosition;
    vec4 tex = vTexCoord;
    
    
    col = texture2DProj(texture, tex);
    
    if (0==1) {
        col.a = 1.;
        col.rgb = WHITE;   
    }
    
    float z = 1.;
    if (col.a >= 0.0 ) {
        if ( wordTimer >= 0. && wordTimer <= 1. ) {
            z = EaseOut(wordTimer);
        } else {

            z = 1. - EaseIn(wordFader);
//            z = Log(wordFader);
        }
    }
    
//    col.rgba = vec4(z,z,z,1.);
    
    col.a = col.a * z;
    

    
//    col.rgb *= col.a;
    
    if (DEBUG) { 
        col.rgb = vec3(0.,0.,1.);
        if (wordTimer >= 0. && wordTimer <= 1. ) { //early
            col.rgb = RED;
        } else  {
            col.rgb = WHITE;
        }
        col.a = z;
    } 


    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
