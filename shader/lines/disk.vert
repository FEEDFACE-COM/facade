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




void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    
    vec4 pos = vec4(vertex,1);


    float RADIUS = tileCount.y/2. ;
    float R0 = RADIUS/4.;
    float rad = RADIUS / (tileCount.y + R0); 


    float delta = 0.0;
//    delta += now/10.;
    delta += ease1(now/4.) - 0.5;


    float ARC = TAU;
    float A0 = 0.0;
  
    float alpha,gamma = 0.;
    
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
    } else if ( pos.x > 0. && pos.y < 0. ) {
        pos.xy = B;
    } else if ( pos.x < 0. && pos.y > 0. ) {
        pos.xy = D;
    } else if ( pos.x < 0. && pos.y < 0. ) {
        pos.xy = C;
    }

    float zoom =  7./ (4.*tileCount.y);
    mat4 mdl = mat4(1.0);
    mdl[0][0] = zoom;
    mdl[1][1] = zoom;
    

    gl_Position = projection * view * mdl * pos;

}

