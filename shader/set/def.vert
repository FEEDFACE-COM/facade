
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float tagCount;
uniform float tagIndex;

uniform float tagWidth;
uniform float tagFader;

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
    
    vec4 pos = vec4(vertex,1);

    pos.y -= 0.5;
    pos.y -= tagCount/2.;
    pos.y += mod( tagIndex, tagCount);    
  
  
    pos.x += tagWidth/2.;
    pos.x -= 5.*ratio;
    
    pos.x += tagFader * 10. ;
    
    
    
    gl_Position = projection * view * model * pos;
}

