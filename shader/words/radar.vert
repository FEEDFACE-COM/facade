
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;


uniform float now;
uniform float debugFlag;

uniform float wordCount;
uniform float wordIndex;

uniform float wordWidth;
uniform float wordFader;
uniform float wordValue;

uniform float charCount;


uniform float screenRatio;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute float charIndex;

varying vec4 vPosition;
varying vec4 vTexCoord;
varying float vCharIndex;


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



       A_____D
       /|   /
      / |  /
     /  | /
    /___|/
   B     C

  

 
                A______D
                 \     \
                  \     \
                   \     \
                    \     \
                     \_____\
                     B      C
   


                    D
                   /\
                  /  \
                A/    \
                 \     \
             /\   \  X  \
            /  \   \     \
            \   \   \     \
          /_ \___\   \_____\
                     B      C



  
*/

mat4 scaleMatrix(float a) {
    return mat4(
         a, 0., 0., 0.,
        0.,  a, 0., 0.,
        0., 0., 1., 0.,
        0., 0., 0., 1.
    );
}



vec2 disk(vec2 pos, float gamma) {

//    vec2 A = vec2( cos(gamma+alpha)*inner, sin(gamma+alpha)*inner);
//    vec2 B = vec2( cos(gamma-alpha)*inner, sin(gamma-alpha)*inner);
//    vec2 C = vec2( cos(gamma-alpha)*outer, sin(gamma-alpha)*outer);
//    vec2 D = vec2( cos(gamma+alpha)*outer, sin(gamma+alpha)*outer);


    float angle, radius;

    float x = pos.x + wordWidth/2.;
    float y = pos.y + 0.5;

    

    vec2 ret = vec2( cos(angle)*radius, sin(angle)*radius );
    
    return ret;
}


void main() {

//    if (wordIndex != 2. ) {
//        return;
//    }
//    
//    if (charIndex != 0.) {
//        return;
//    }


    vec4 pos = vec4(vertex,1);
    
//    pos.x -= charOffset;
    
    
    float charWidth = abs(pos.x*2.);
    
    vec4 tex; 
    tex.xy = texCoord.xy;
    tex.w = 1.;

    vTexCoord = tex;

    
    float center = 1. + wordCount / 8. + wordWidth/2.;
    
    float inner = center + /*charOffset*/ - charWidth/2.;
    float outer = inner + charWidth;

    float alpha = (TAU/wordCount) / 2.0;
    float gamma = (TAU/wordCount) * -1. * wordIndex;

    vPosition = pos;
    vTexCoord = tex;
    vCharIndex = charIndex;


    vec2 A = vec2( cos(gamma+alpha)*inner, sin(gamma+alpha)*inner);
    vec2 B = vec2( cos(gamma-alpha)*inner, sin(gamma-alpha)*inner);
    vec2 C = vec2( cos(gamma-alpha)*outer, sin(gamma-alpha)*outer);
    vec2 D = vec2( cos(gamma+alpha)*outer, sin(gamma+alpha)*outer);

     
    float w,n; //wide,narrow
    w = A.y - B.y;
    n = D.y - C.y;
    



    vec4 a,b,c,d;

    a = vec4(tex.x+0.,tex.y+0.,0.,w);
    b = vec4(tex.x+0.,tex.y+w ,0.,w);
    c = vec4(tex.x+n ,tex.y+n ,0.,n);
    d = vec4(tex.x+n ,tex.y+0.,0.,n);

//    a = vec4(0.,0.,0.,w);
//    b = vec4(0.,w,0.,w);
//    c = vec4(n,n,0.,n);
//    d = vec4(n,0.,0.,n);
    

    if        ( pos.x < 0. && pos.y > 0. ) {
        pos.xy = A;
        tex = a;
    } else if ( pos.x < 0. && pos.y < 0. ) {
        pos.xy = B;
        tex = b;
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xy = C;
        tex = c;
    } else if ( pos.x > 0. && pos.y > 0. ) {
        pos.xy = D;
        tex = d;
    }
    

    pos.z += + 0.25 * cos( now + outer + PI/2.);
    pos.x += + 0.25 * sin( now + inner);
    pos.y += + 0.25 * cos( now + inner);


    mat4 R = mat4(1.0);
    R = rotationMatrix(vec3(1.,0.,0.), sin(now/2.) * PI/15.);
    R *= rotationMatrix(vec3(0.,1.,0.), sin(now/2.) * PI/13.);
    R *= rotationMatrix(vec3(0.,0.,1.), now/-8.);
    R *= scaleMatrix(wordCount/16.);
    pos = R * pos;

    
        
    gl_Position = projection * view * model * pos;
}






