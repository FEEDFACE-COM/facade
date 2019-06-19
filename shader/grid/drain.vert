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


bool oddColCount() { return mod(tileCount.x, 2.0) == 1.0 ; }
bool oddRowCount() { return mod(tileCount.y, 2.0) == 1.0 ; }

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);
    if (pos.x < 0.) {
//        pos.x *= 2.;
    }
    if (pos.y < 0.) {
//        pos.y *= 2.;
    }

    float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;

    float a = tileCoord.x / tileCount.x * 2. * PI;// + 0.25* sin(now);
    float a0 = a + PI/2.;
    float a1 = a;
    float radius = 4. + tileCoord.y / tileCount.y * 5.;
    
    
    radius -= vScroller;


    pos.x *= radius/2.;

    vec2 p = vec2(pos.x,pos.y);
    
    
    mat2 rotate = mat2(
        cos(-a0),  -1. * sin(-a0),
        sin(-a0),  cos(-a0) 
    );
    
    
    pos.xy = rotate * p;
    
//    pos.xy *= 1. * radius;
    
    pos.x += cos(a1) * radius;
    pos.y += sin(a1) * radius;

//    vec2 v = vec2(10., 20.);
//    mat2 m = mat2(1., 2.,  3., 4.);
//    vec2 w = rotate * v; // = vec2(1. * 10. + 3. * 20., 2. * 10. + 4. * 20.)


    gl_Position = projection * view * model * pos;
}

