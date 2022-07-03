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

float Ease(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

mat3 rotx(float w) {return mat3(1.0,0.0,0.0,0.0,cos(w),sin(w),0.0,-sin(w),cos(w));}
mat3 roty(float w) {return mat3(cos(w),0.0,sin(w),0.0,1.0,0.0,-sin(w),0.0,cos(w));}
mat3 rotz(float w) {return mat3(cos(w),sin(w),0.0,-sin(w),cos(w),0.0,0.0,0.0,0.0);}

void main() {
    vec4 pos = vec4(vertex,1);

    vPosition =    pos;
    vTexCoord = vec4(texCoord.xy,1.,1.);
    vCharIndex = charIndex;


    float w = -PI/4. * -sin( (charIndex+1.-scroller) * PI/charCount  - PI/2. + PI/2.);

    float ARC = TAU;
    float GAP = 1.;
    
    
    float c = (charCount+GAP) * fontRatio;  // circumference
    float r = c / TAU;                      // radius
    r += log(charCount + GAP);              // adjust smaller charcounts
    
    
    float gamma = 0.;
    gamma += 1.*PI/64.;                     // adjust start point
    
    float x = (charIndex-scroller) / (charCount + GAP);
    float b = 2.;

    pos.xyz *= rotz( b*PI/32. * cos(x*b*TAU  ) );
    pos.y +=                 cos(x*b*TAU - 3.*PI/2. );

    gamma -= x;


    pos.xyz *= roty( -gamma*ARC + PI/2. );  // face outwards
    
    pos.x += r * cos( gamma * ARC );
    pos.z += r * sin( gamma * ARC );
    



    pos.xyz *= rotx( -PI/16. * cos(now/2.) ); // animated swifel
//   pos.xyz *= roty( now / 8.); // rotate whole ring

    float zoom = 1.0;
    {
        zoom = 4. / (3.*r);
    }
    pos.xyz *= zoom;
    
    
                                                                                                                                                                                                                                            
    gl_Position = projection * view * model * pos;
}
