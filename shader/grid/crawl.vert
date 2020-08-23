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


void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

//    pos.y += scroller;
//    pos.x += (tileCoord.x * tileSize.x);
//    pos.y += (tileCoord.y * tileSize.y);
//
//    pos.x += ( tileOffset.x * tileSize.x);
//    pos.y += ( tileOffset.y * tileSize.y);
    
//    float w = tileCount.x * tileSize.x;
//    float h = tileCount.y * tileSize.y;
    
    
    
//    vec2 A = vec2(-w/2., h/2.);
//    vec2 B = vec2(-w/2.,-h/2.);
//    vec2 C = vec2( w/2., h/2.);
//    vec2 D = vec2( w/2.,-h/2.);
    
//    vec2 Ai = vec2(-w/4., h/2.);
//    vec2 Bi = vec2(-w/2.,-h/2.);
//    vec2 Ci = vec2( w/2., h/2.);
//    vec2 Di = vec2( w/4.,-h/2.);

//    float dx,dy;
//    dx = (tileCoord.x+tileCount.x/2.) / tileCount.x;
//    dy = (tileCoord.y+tileCount.y/2.) / tileCount.y;


    float w = tileSize.x;
    float h = tileSize.y;
    
    float cols = tileCount.x;
    float rows = tileCount.y;

    float x = tileCoord.x;
    float y = tileCoord.y;    
    
    float dx = 0.;
    float dy = 0.;
    
//    pos.x = w * dx(x);
//    pos.y = h * dy(y);
        
//    float H = tileCount.y * tileSize.y;
//    float alpha = asin(1./ H );
//    float x = H * cos(alpha);    

    vec2 orig = pos.xy;
    if (        orig.x < 0.0 && orig.y > 0.0 ) {
        dx = 0.0;
        dy = 1.0;
    } else if ( orig.x < 0.0 && orig.y < 0.0 ) {
        dx = 0.0;
        dy = 0.0;
    } else if ( orig.x > 0.0 && orig.y < 0.0 ) {
        dx = 1.0;
        dy = 0.0;
    } else if ( orig.x > 0.0 && orig.y > 0.0 ) {
        dx = 1.0;
        dy = 1.0;
    }        

    float px = 0.5*((x+dx+cols/2.) / cols) + 0.5;
    float py = 0.5*((-y+dy+rows/2.) / rows) + 0.5;


    float h1 = 1.;
    float h0 = 2.;

    pos.x = (x+dx) * w * (py*(h0-h1));
    pos.y = (y+dy) * h ;


//    pos = pos*rollMatrix(now);
    
    float zoom = 1.0;
    float ratio = screenRatio / fontRatio;
    float zoom_rows = model[0][0];
    float zoom_cols = ratio * 2. / tileCount.x;
    
    zoom = zoom_rows;
        
    
//    float ALPHA;
//    ALPHA = PI/4.;
//    ALPHA = PI * 3./8.;
//    ALPHA = alpha;
//    mat4 rot;
//    rot = rollMatrix(alpha);
    


//    pos.y += h/2.;
//    pos = rot * pos;
//    pos.y -= tileCount.y/1.;
    
    mat4 mdl;
    mdl = mat4(1.0);
    mdl[0][0] = zoom;
    mdl[1][1] = zoom;
    
    

    gl_Position = projection * view * mdl * pos;
}

