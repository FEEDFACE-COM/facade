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

bool oddColCount() { return mod(tileCount.x, 2.0) == 1.0 ; }
bool oddRowCount() { return mod(tileCount.y, 2.0) == 1.0 ; }


bool upperVertex(vec2 vec) {
    return vec.y >= 0.; 
}


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
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);


    float rx = tileCoord.x / tileCount.x;

    float r0 = 4.;
    float r1 = 1.;
    
    float phi = (tileCoord.x / tileCount.x) * TAU;
    
    if ( upperVertex(pos.xy) ) {
        pos.x += cos(phi) * r0;
    } else {
        pos.x += cos(phi) * r1;
    }
//    if (pos.y > 0.) {
//        pos.x = cos(phi) * r0;
//    } else {
//        pos.y = cos(phi) * r1;
//    }


    vec3 axis = vec3(-1.,-1.,0.);
    mat4 rot = rotationMatrix(axis, PI/4.);
//    pos = rot * pos;
    
    gl_Position = projection * view * model * pos;
}

