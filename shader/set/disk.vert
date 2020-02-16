
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

    vec4 pos = vec4(vertex,1);

/*
    d   c
    +---+
    |   |
    +---+
    a   b


     c
    /\
  d/  \
  /\   \
 /  \   \
+----\---\---
     a   b



*/

    float RADIUS = 4.;
    float R0 = 1.0;


    float ARC = TAU;
    float A0 = 2.0;
  
    float alpha,gamma;
    
    float row = (-tagIndex+tagMax/2.);


//    alpha = ARC / (A0 + tagMax);
//    gamma = ( ARC / (tagMax+A0)) * tagIndex;

    
    float r0 = R0;
    float r1 = r0 + tagWidth;


    alpha = ARC / tagMax;
    gamma = (ARC/tagMax) * tagIndex;


    vTexCoord = texCoord;


    vec2 A = vec2( cos(gamma      )*r0, sin(gamma      )*r0);
    vec2 B = vec2( cos(gamma      )*r1, sin(gamma      )*r1);
    vec2 C = vec2( cos(gamma+alpha)*r1, sin(gamma+alpha)*r1);
    vec2 D = vec2( cos(gamma+alpha)*r0, sin(gamma+alpha)*r0);

    if        ( pos.x < 0. && pos.y < 0. ) {
        pos.xy = A;
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xy = B;
    } else if ( pos.x > 0. && pos.y > 0. ) {
        pos.xy = C;
    } else if ( pos.x < 0. && pos.y > 0. ) {
        pos.xy = D;
    }
     

    
    
        
    gl_Position = projection * view * model * pos;
}

