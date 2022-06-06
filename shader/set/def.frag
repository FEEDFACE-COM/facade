
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

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return        cos(x*PI/2. + 3.*PI/2. );        }
float EaseIn(float x) {  return  -1. * cos(x*PI/2.            ) + 1.  ; }


void main() {
    vec4 col;
    vec4 pos = vPosition;
    vec4 tex = vTexCoord;

    col = texture2DProj(texture, tex);
    col.a = wordFader;

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
    }

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
