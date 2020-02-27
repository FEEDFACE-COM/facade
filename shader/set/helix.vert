
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float wordMax;
uniform float wordIndex;

uniform float wordWidth;
uniform float wordFader;
uniform float wordCount;

uniform float ratio;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec4 vPosition;
varying vec4 vTexCoord;


bool DEBUG = debugFlag > 0.0;



float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }

mat4 rotationMatrix(vec3 axis, float angle)
{
    vec3 a  = normalize(axis);
    float s = sin(angle);
    float c = cos(angle);
    float oc = 1.0 - c;
    
    return mat4(
        oc*a.x*a.x + c,      oc*a.x*a.y - a.z*s,  oc*a.z*a.x + a.y*s,  0.0,
        oc*a.x*a.y + a.z*s,  oc*a.y*a.y + c,      oc*a.y*a.z - a.x*s,  0.0,
        oc*a.z*a.x - a.y*s,  oc*a.y*a.z + a.x*s,  oc*a.z*a.z + c,      0.0,
                       0.0,                 0.0,                 0.0,  1.0
    );
}

mat4 scaleMatrix(float z) {
    return mat4(
        z,0.0,0.0,0.0,
      0.0,  z,0.0,0.0,
      0.0,0.0,  z,0.0,
      0.0,0.0,0.0,1.0
    );
}


/*

     A          D
 -w/2,h/2____w/2,h/2
     |          |
     |          |
 -w/2,-h/2___w/2,-h/2
     B          C


     A          D
    0,0________1,0
     |          |
     |          |
    0,1________1,1
     B          C

*/


float H = 32.;


vec3 helix(float radius, float gamma, float alpha) {
    radius *= 8.;
    return vec3( cos(gamma+alpha)*radius, sin(gamma+alpha)*radius, H * (gamma+alpha) );
}

vec3 helix2(float radius, float gamma, float alpha) {
    radius *= -8.;
    return vec3( cos(gamma+alpha)*radius, sin(gamma+alpha)*radius, H * (gamma+alpha) );
}



void main() {

    vec4 pos = vec4(vertex,1);
    
    vec4 tex; 
    tex.xy = texCoord.xy;

    vTexCoord = tex;

    float X = 8.;

    float TOTAL =  X * TAU;
    float t = (wordIndex/wordMax) * TOTAL;
    

    float alpha,gamma;


    alpha = 1./4.;
    gamma = t ;
    gamma -= TOTAL/2.;


    float inner = 4.;
    float outer = inner + 4. * 4. * X ;
//    outer = 16.;
    
    //outer = inner + 4.;


    vPosition = pos;
    
    vec3 A,B,C,D;

     A = helix(inner,gamma,+alpha);
     B = helix(inner,gamma,-alpha);
     

     
     C = helix2(inner,gamma,-alpha);
     D = helix2(inner,gamma,+alpha);

     
    float w,n; //wide,narrow
    w = A.y - B.y;
    n = D.y - C.y;

    vec4 a,b,c,d;

    d = vec4(n,0.,0.,n);
    a = vec4(0.,0.,0.,w);
    b = vec4(0.,w,0.,w);
    c = vec4(n,n,0.,n);
    

    if        ( pos.x < 0. && pos.y > 0. ) {
        pos.xyz = A;
        tex = a;
    } else if ( pos.x < 0. && pos.y < 0. ) {
        pos.xyz = B;
        tex = b;
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xyz = C;
        tex = c;
    } else if ( pos.x > 0. && pos.y > 0. ) {
        pos.xyz = D;
        tex = d;
    }

    




    float z = 1./8.;
//    z = wordMax / 32.;    

    mat4 V = view;
    V *= rotationMatrix(vec3(0.,1.,0.), PI/2.);
    V *= rotationMatrix(vec3(1.,0.,0.), PI/6.);
    
//    R = rotationMatrix(vec3(1.,0.,0.), sin(now/2.) * PI/15.);
//    R *= rotationMatrix(vec3(0.,1.,0.), sin(now/2.) * PI/13.);
    V *= rotationMatrix(vec3(0.,0.,1.), now);
//    R *= rotationMatrix(vec3(0.,0.,1.), now/-1.);
//    V *= scaleMatrix(z);

    V *= scaleMatrix(z);
    
    
        
    vTexCoord = tex;
    gl_Position = projection * V * model * pos;
}





