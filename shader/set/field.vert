
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float tagMax;
uniform float tagIndex;

uniform float tagWidth;
uniform float tagFader;
uniform float tagCount;

uniform float ratio;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vTexCoord;


bool DEBUG = debugFlag > 0.0;


float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }


void main() {
    vTexCoord = texCoord;
    float fader = tagFader;
    
    vec4 pos = vec4(vertex,1);

    float a;
    float x,y;

    pos.x *= tagMax/10.;
    pos.y *= tagMax/10.;
    
    float rx = tagMax/4.;
    float ry = tagMax/2.;
    
    float idx = tagIndex/tagMax;
    
    float w = 2.*TAU;
    float o = PI/2. + PI/2.;
    a = idx * w + o;

    pos.x += rx * cos(a);
    pos.y += ry * sin(a);
        

    if (idx >= 0.5 ) {
        pos.x -= -tagMax/2.;
    } else {
        pos.x -= tagMax/2.;
    
    }
    pos.z -= tagMax; 
    pos.z += fader * 1.5 * tagMax;

        
    gl_Position = projection * view * model * pos;
    vPos =    gl_Position;
}

