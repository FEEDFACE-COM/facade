uniform mat4 model;               // model transformation
uniform mat4 view;                // view transformation
uniform mat4 projection;          // projection transformation

uniform float now;                // time since start of program, as seconds
uniform float debugFlag;          // 0.0 unless -D flag given by user

uniform float scroller;

uniform float charCount;

uniform float screenRatio;
uniform float fontRatio;

attribute vec3 vertex;            // vertex position as (x,y,z) centered on (0,0,0)
attribute vec2 texCoord;
attribute float charIndex;

varying vec4 vTexCoord;
varying vec4 vPosition;
varying float vCharIndex;

bool DEBUG = debugFlag > 0.0;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

void main() {
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;



    pos.x += abs(pos.x);
//    pos.x += 1. * fontRatio;
    pos.x += (charIndex) * fontRatio;
    pos.x -= (charCount*fontRatio)/2.;
    pos.x -= scroller * fontRatio;


    pos.y += 2.*sin( (charIndex+1.-scroller) * TAU/charCount );

    float zoom = 1.0;
    {
        zoom = 1./(((charCount+1.)*fontRatio)/2.) * screenRatio;
    }
    pos.xy *= zoom;
                                                                                                                                                                                                                                            
    gl_Position = projection * view * model * pos;
}
