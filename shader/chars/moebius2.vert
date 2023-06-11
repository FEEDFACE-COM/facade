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
mat3 rotz(float w) {return mat3(cos(w),sin(w),0.0,-sin(w),cos(w),0.0,0.0,0.0,1.0);}



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
    }


    float c = count * fontRatio;  // circumference
    float r = c / TAU;            // radius
    r += log(charCount);          // adjust smaller charcounts


    float a;
    
    // first half: start at index #0
    if (charIndex < count) {

        a = PI + (index - scroller) / (count+1.) * ARC;

    } 
    // second half: start at index count/2
    if (charIndex >= count) {

        a = PI + (count + index - scroller) / (count+1.) * ARC;
        
    }

    float d = TAU / count;
    float s = pos.y / fontRatio;
    float t = pos.x >= 0. ? a+d/2. : a-d/2. ;


    float x,y,z;

    { // simple plane
        x = s;
        y = 0.;
        z = 0.;
    }

    if ( false ) { // torus with circular cross section?
        float A = 1.;
        float B = 1.;
        
        x = A * cos(s);
        y = 0.;
        z = B * sin(s); 
        
    }
    
    
    // Moebius Strip: https://mathworld.wolfram.com/MoebiusStrip.html
    // pos.x = ( r + s * cos( t/2.) ) * cos(t);
    // pos.y = ( r + s * cos( t/2.) ) * sin(t);
    // pos.z =       s * sin( t/2.);
    

    
    // Moebius Strip: https://math.stackexchange.com/a/1396314   
    pos.x = cos(t) * ( r - sin(t/2.) * z + cos(t/2.) * x ) - sin(t) * y;
    pos.y = sin(t) * ( r - sin(t/2.) * z + cos(t/2.) * x ) + cos(t) * y;
    pos.z =                cos(t/2.) * z + sin(t/2.) * x ;


    mat3 rotate = mat3(1.0);

    rotate *= roty( PI + PI/4.);

    rotate *= rotz( -now/32. );
    rotate *= rotx(  PI - PI/8. * sin(now/5.) );
    rotate *= roty( PI/2. + 3.*PI/2. + sin(now/7.) + now/96. );


    
    float zoom = 1.0;
    {
        zoom = 7. / (8.*r);
    }
    pos.xyz *= zoom;
    

    mat4 rot = mat4(0.0);
    for (int i =0; i < 3; i++) {
        for (int j=0; j<3; j++) {
            rot[i][j] = rotate[i][j];
        }
    }
    
    rot[3][3] = 1.0;


                                                                                                                                                                                                                                            
    gl_Position = projection * view * rot* model * pos;
}
