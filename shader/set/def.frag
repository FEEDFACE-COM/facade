
uniform float now;
uniform float debugFlag;
uniform sampler2D texture;

uniform float tagCount;
uniform float tagMaxWidth;
uniform float tagFader;
uniform float tagIndex;


varying vec2 vTexCoord;

bool DEBUG    = debugFlag > 0.0;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return        cos(x*PI/2. + 3.*PI/2. );        }
float EaseIn(float x) {  return  -1. * cos(x*PI/2.            ) + 1.  ; }


void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 


    if (col.a > 0.0 ) {
        if ( tagFader < .125 ) {
    //        col.g = 0.0;
    //        col.b = 0.0;
            col.a = EaseOut(tagFader * 8.);
        } else if (tagFader < .75 ) {
    //        col.r = 0.0;
    //        col.b = 0.0;
            col.a = 1.0;
        } else {
    //        col.r = 0.0;
            col.a = 1. - EaseOut(3./4. + tagFader * 4.);
        }
    }
    
    col.rgb *= col.a;
    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
//        col.g = tagIndex / tagCount;
        col.a = 1.0;
    } 

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
