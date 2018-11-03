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


    float c = tileCount.x * tileSize.x;

    // c = 2π * r <=> c/2π = r //
    float r = c / TAU;
    
    float a;
    a = (tileCoord.x / (0.5*tileCount.x + 2.)) * PI - PI/8.;

    a += now/10.;
    
    
    a += ease1(now/2.);

    pos = rotationMatrix(vec3(1.,0.,0.), PI/2.) * pos;
    pos = rotationMatrix(vec3(0.,0.,1.), -a-PI/2.) * pos;
    

    pos.x +=  cos(a) * r;
    pos.y +=  sin(a) * r;


    pos.z -= tileCoord.y;
    pos.z -= scroller;


    
    vec3 axis = vec3(-1.,-1.,0.);
    mat4 rot = rotationMatrix(axis, PI/2.);
    pos = rot * pos;

//    pos.x += (tileCoord.x * tileSize.x);
//    pos.y += (tileCoord.y * tileSize.y);
    
        
    

    
    gl_Position = projection * view * model * pos;
}
