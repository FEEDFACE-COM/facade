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

bool DEBUG = debugFlag > 0.0;


float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU= 6.2831853071795864769252867665590057683943387987502116419498891840;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

mat3 rotx(float w) {return mat3(1.0,0.0,0.0,0.0,cos(w),sin(w),0.0,-sin(w),cos(w));}
mat3 roty(float w) {return mat3(cos(w),0.0,sin(w),0.0,1.0,0.0,-sin(w),0.0,cos(w));}
mat3 rotz(float w) {return mat3(cos(w),sin(w),0.0,-sin(w),cos(w),0.0,0.0,0.0,0.0);}



void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    
    vec4 pos = vec4(vertex,1);

    float offset = PI/4.;
//    offset = 0.0;
    float ARC = 6.*PI/4. - offset;
    float RADIUS = tileCount.y/2. * tileSize.y /2. ;

//     RADIUS += .25 * cos( tileCoord.x/tileCount.x * PI + PI*now );

    

    float delta = 0.0;
    float alpha,beta;

    float x = tileCoord.x-0.5+tileOffset.x;
    float y = tileCoord.y-1.5+tileOffset.y;
    

    alpha = -1. * ARC / (tileCount.y);
    delta = PI/2. - offset + alpha;
    beta = delta + ( alpha * (scroller+y) ) ;


    float r = RADIUS * 2.;
    
    
    vec3 A = vec3( (x+1.)*tileSize.x, cos(alpha+beta)*r, -tileCount.y/2.+sin(alpha+beta)*r);
    vec3 B = vec3( (x+1.)*tileSize.x, cos(beta)*r,       -tileCount.y/2.+sin(beta)*r);
    vec3 C = vec3( (x)*tileSize.x,    cos(alpha+beta)*r, -tileCount.y/2.+sin(alpha+beta)*r);
    vec3 D = vec3( (x)*tileSize.x,    cos(beta)*r,       -tileCount.y/2.+sin(beta)*r);
    
   
    if ( pos.x > 0. && pos.y > 0. ) {
        pos.xyz = A;
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xyz = B;
    } else if ( pos.x < 0. && pos.y > 0. ) {
        pos.xyz = C;
    } else if ( pos.x < 0. && pos.y < 0. ) {
        pos.xyz = D;
    }

    float ratio = screenRatio / fontRatio;
    float zoom = ratio * 2. / ( tileCount.x );
    
    mat4 mdl = mat4(1.0);
    mdl[0][0] = zoom;
    mdl[1][1] = zoom;
    mdl[2][2] = zoom;

    gl_Position = projection * view * mdl * pos;
}

