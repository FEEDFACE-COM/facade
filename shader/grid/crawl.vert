uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;
uniform vec2 tileOffset;

uniform float now;
uniform float scroller;
uniform float screenRatio;
uniform float fontRatio;
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


float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }


mat4 yawMatrix(float angle) 
{
    return mat4(
        cos(angle), -sin(angle), 0.0, 0.0,
        sin(angle),  cos(angle), 0.0, 0.0,
        0.0,                0.0, 1.0, 0.0,
        0.0,                0.0, 0.0, 1.0
    
    );
}


mat4 rollMatrix(float angle) 
{
    return mat4(
        1.0,        0.0,         0.0, 0.0,
        0.0, cos(angle), -sin(angle), 0.0,
        0.0, sin(angle),  cos(angle), 0.0,
        0.0,        0.0,         0.0, 1.0
    
    );
}

mat4 pitchMatrix(float angle)
{
    return mat4(
        cos(angle), 0.0, sin(angle), 0.0,
               0.0, 1.0,        0.0, 0.0,
       -sin(angle), 0.0, cos(angle), 0.0,
               0.0, 0.0,        0.0, 1.0
    );

}

//mat4 rotationMatrix(vec3 axis, float angle)
//{
//    vec3 a  = normalize(axis);
//    float s = sin(angle);
//    float c = cos(angle);
//    float oc = 1.0 - c;
//    
//    return mat4(
//        oc*a.x*a.x + c,      oc*a.x*a.y - a.z*s,  oc*a.z*a.x + a.y*s,  0.0,
//        oc*a.x*a.y + a.z*s,  oc*a.y*a.y + c,      oc*a.y*a.z - a.x*s,  0.0,
//        oc*a.z*a.x - a.y*s,  oc*a.y*a.z + a.x*s,  oc*a.z*a.z + c,      0.0,
//                       0.0,                 0.0,                 0.0,  1.0
//    );
//}

//float dx(float x) {
//    return (x+tileCount.x/2.)/tileCount.x;
//}
//
//float dy(float y) {
//    return (y+tileCount.y/2.)/tileCount.y;
//}

float dy(float y) {
    return 1. - (  0.25 + 0.5 * ((y+tileCount.y/2.) / tileCount.y ) );
}

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    float dy0 = dy(tileCoord.y+1.);
    float dy1 = dy(tileCoord.y);


    
    float w0 = dy0 * tileSize.x;
    float w1 = dy1 * tileSize.x;
    float h  = tileSize.y;
    
    float x = tileCoord.x;
    float y = tileCoord.y;


    vec2 A = vec2((w0*(x-1.)) , (h * (y   )));
    vec2 B = vec2((w1*(x-1.)) , (h * (y-1.)));
    vec2 C = vec2((w1*(x   )) , (h * (y-1.)));
    vec2 D = vec2((w0*(x   )) , (h * (y   )));
    
    
   
    if        ( pos.x < 0. && pos.y > 0. ) {
        pos.xy = A;
    } else if ( pos.x < 0. && pos.y < 0. ) {
        pos.xy = B;
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xy = C;
    } else if ( pos.x > 0. && pos.y > 0. ) {
        pos.xy = D;
    }
    
    pos.xy = pos.xy + (scroller * (A-B));

    float ratio = screenRatio / fontRatio;
    float cols = ratio * 2. / tileCount.x;
    float rows = model[0][0];
    float zoom = rows;
    mat4 mdl;
    mdl = mat4(1.0);
    mdl[0][0] = zoom;
    mdl[1][1] = zoom;
    
    

    gl_Position = projection * view * mdl * pos;
}

