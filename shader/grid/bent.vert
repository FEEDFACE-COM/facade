uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;
uniform vec2 tileOffset;

uniform float now;
uniform float scroller;
uniform float debugFlag;

attribute vec3 vertex;

attribute vec2 texCoord;
attribute vec2 tileCoord;
attribute vec2 gridCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying float vScroller;

bool DEBUG = debugFlag > 0.0;


float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU= 6.2831853071795864769252867665590057683943387987502116419498891840;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

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



void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);


    float RADIUS = 10.;
    float R0 = 5.0;
    float rad = RADIUS / (tileCount.y + R0); 


    float delta = 0.0;
//    delta += now/10.;
    delta += ease1(now/4.) - 0.5;
    

    float ARC = TAU;
    float A0 = 0.0;
  
    float alpha,gamma;
    
    float row = (-tileCoord.y+tileCount.y/2.);


    alpha = ARC / (A0 + tileCount.x);
    gamma += delta;
    gamma += ( ARC / (tileCount.x+A0)) * tileCoord.x;


    
    float r0 = R0 + (rad * row ) ;
    float r1 = r0 + rad;

    r0 -= (scroller*rad);
    r1 -= (scroller*rad);

    
    vec2 A = vec2( cos(gamma+alpha)*r0, sin(gamma+alpha)*r0);
    vec2 B = vec2( cos(gamma+alpha)*r1, sin(gamma+alpha)*r1);
    vec2 C = vec2( cos(gamma      )*r1, sin(gamma      )*r1);
    vec2 D = vec2( cos(gamma      )*r0, sin(gamma      )*r0);
    
   
   
    if        ( pos.x > 0. && pos.y > 0. ) {
        pos.xy = A;
        pos.z +=  tileSize.x*(scroller+tileCoord.y+1.)/tileCount.y * 16.;
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xy = B;
        pos.z +=  tileSize.x*(scroller+tileCoord.y   )/tileCount.y * 16.;
    } else if ( pos.x < 0. && pos.y > 0. ) {
        pos.xy = D;
        pos.z +=  tileSize.x*(scroller+tileCoord.y+1.)/tileCount.y * 16.;
    } else if ( pos.x < 0. && pos.y < 0. ) {
        pos.xy = C;
        pos.z +=  tileSize.x*(scroller+tileCoord.y   )/tileCount.y * 16.;
    }

    pos.z +=  tileCoord.x/tileCount.x * 8. ;



    mat4 R = mat4(1.0);
    R = rotationMatrix(vec3(1.,0.,0.), sin(now/2.) * PI/15.);
    R *= rotationMatrix(vec3(0.,1.,0.), sin(now/2.) * PI/13.);
    pos = R * pos;

    gl_Position = projection * view * model * pos;
}

