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


    float ARC = TAU;
    
    float count = floor(charCount/2.);
    float index = charIndex;
    if (index >= count) {
        index -= count;
    } else {
    }



    float c = count * fontRatio;  // circumference
    float r = c / TAU;                      // radius
    r += log(charCount);              // adjust smaller charcounts


    bool one = false;
    bool two = false;


    float over = 2.;
    if ( mod(now, over) >= over/2. ) {
        one = true;            
    } else {
        two = true;
    }

    if (one && charIndex < count) {


        float a = (/*1.+*/index - scroller) / (count+1.) * ARC;
        float d = TAU / count;

        float s = pos.y / fontRatio;
        float t = pos.x >= 0. ? a+d/2. : a-d/2. ;
    
    

        // Moebius Strip: https://mathworld.wolfram.com/MoebiusStrip.html
    
        pos.x = ( r + s * cos( t/2.) ) * cos(t);
        pos.y = ( r + s * cos( t/2.) ) * sin(t);
        pos.z = s * sin( t/2. );

        
        
    } 
    
    if (two && charIndex >= count) {

        float a = (count+index - scroller) / (count+1.) * ARC;
        float d = TAU / count;
    
        float s = pos.y / fontRatio;
        float t = pos.x >= 0. ? a+d/2. : a-d/2. ;
    
    

        // Moebius Strip: https://mathworld.wolfram.com/MoebiusStrip.html
    
        pos.x = ( r + s * cos( t/2.) ) * cos(t);
        pos.y = ( r + s * cos( t/2.) ) * sin(t);
        pos.z = s * sin( t/2. );
        
        
        
    }
    

    mat4 rotate = mat4(1.0);
//    rotate = mat4( rotx( PI/2. - now*PI/16. ) );
//    rotate = mat4( rotx( PI/2. )  );
//    rotate *= mat4( roty(now/3.) );
//    rotate *= mat4( rotx(now/4. + PI));


    rotate *= mat4( roty( sin(now/3.) * PI/4. ) );

    rotate *= mat4( rotx(0.) );
    rotate *= mat4( roty(PI/2.) );
    rotate *= mat4( rotz(0.) );
    float zoom = 1.0;
    {
        zoom = 15. / (8.*r);
    }
    pos.xyz *= zoom;
    



                                                                                                                                                                                                                                            
    gl_Position = projection * view * rotate* model * pos;
}
