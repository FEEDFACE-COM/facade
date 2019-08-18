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




void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    float offset = PI/16.;
    offset = 0.0;
//    offset = 0.0;
    float ARC = PI/2.;// - offset;
    float RADIUS = tileCount.y/2. * tileSize.y /2. ;


    

    float delta = 0.0;
    float alpha,beta;
    

    alpha = -1. * ARC / (tileCount.y);
    delta = PI/2. - offset + alpha;
    beta = delta + ( alpha * (scroller+tileCoord.y) ) ;


    float r = RADIUS * 2.;
    
    vec3 A = vec3( (tileCoord.x+1.)*tileSize.x, cos(alpha+beta)*r, sin(alpha+beta)*r);
    vec3 B = vec3( (tileCoord.x+1.)*tileSize.x, cos(beta)*r,       sin(beta)*r);
    vec3 C = vec3( tileCoord.x*tileSize.x,      cos(alpha+beta)*r, sin(alpha+beta)*r);
    vec3 D = vec3( tileCoord.x*tileSize.x,      cos(beta)*r,       sin(beta)*r);
    
   
    if ( pos.x > 0. && pos.y > 0. ) {
        pos.xyz = A;
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xyz = B;
    } else if ( pos.x < 0. && pos.y > 0. ) {
        pos.xyz = C;
    } else if ( pos.x < 0. && pos.y < 0. ) {
        pos.xyz = D;
    }

    float zoom = 0.65;
    

//
//
    float fontRatio = tileSize.x/tileSize.y;
    float screenRatio = (tileCount.x*tileSize.x)/((tileCount.y)*tileSize.y);
    float ratio = screenRatio / fontRatio;

    float scaleWidth = ratio * 2. / tileCount.x;
    float scaleHeight =        2. / tileCount.y;
    


//    if ( scaleWidth < scaleHeight/2. ) {
//        zoom = scaleWidth;
//    } else {            
//        zoom = scaleHeight;
//    }
    zoom *= scaleWidth ;

    pos.xyz *= zoom;
//    pos.xyz *= model[0][0]  * 0.8;

    gl_Position = projection * view * pos;
}

