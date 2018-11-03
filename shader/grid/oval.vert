uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;



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
    
    vec4 pos = vec4(vertex,1);
    
//    pos = rotationMatrix(vec3(1., 0., 0.), PI  ) * pos;
    
    pos.z = pos.y;
    pos.y = 0.;



    vec2 w = vec2(tileSize.x*tileCount.x, tileSize.y*tileCount.y);

    vec2 coord = vec2( 
        tileCoord.x + (tileCount.x/2.  - 1.),
        tileCoord.y + (tileCount.y/1.  - 1.)
    );
    
    vec2 grad = vec2(coord.x / tileCount.x, coord.y / tileCount.y);

    // circum = 2π * radius <=> circum/2π = radius //
    float circum = tileCount.x * tileSize.x;
    float radius = circum / TAU;
    
    float a,b;
    float zoom = 2.;
    a = zoom * radius;
    b = zoom * radius/2.;

    

    float phase = -1. * PI/8.; 
    
    
//    phase += PI/2. * (ease1( now/8. ) );
    
    float alpha = grad.x * TAU + phase;
 
//    alpha += now/10.;
 
    
    alpha = alpha * -1. + PI/4.;
        
    pos = rotationMatrix(vec3(0., 0., 1.), alpha) * pos;


    

//    a += 0.25 * cos(0.3*now * ease1((  mod(now,TAU) +scroller)/100.) * vTileCoord.y);

//    vec3 tmp = vec3(pos);
    

    pos.x -= cos(alpha) * a;
    pos.y -= sin(alpha) * b;

    pos.z += (tileCoord.y*tileSize.y);
    pos.z += scroller;

//    pos = rotationMatrix(vec3(-1.,1.,0.), PI/4.) * pos;
//    pos = rotationMatrix(vec3(1.,0.,0.), PI/2.) * pos;

    
    gl_Position = projection * view * model * pos;
}
