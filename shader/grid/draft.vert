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

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);
    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);



    vec4 displacedVertex;
    displacedVertex = pos;
    
    float len = length( displacedVertex );
    float freq = 1./4.;

    {
        displacedVertex.x += -4.0 * sin(now * 1./4. + len + 1.);
        displacedVertex.y += 0.5 * sin(now * 1./5. + len);
        displacedVertex.z += 5.0 * cos(now * 1./3. + len);
    }

//    gl_Position = gl_ModelViewProjectionMatrix * displacedVertex;      


    pos = displacedVertex;


    gl_Position = projection * view * model * pos;
}


//float fun0(float x) {
//    float a = 1.0;
//    float f = 1.0;
//    float p = 0.0;
//    float b = 0.0;
//    return a * cos( f * x * PI + p ) + b;
//}
//
//float fun(float x, float a, float f, float p, float b) {
//    return a * cos( f * x * PI + p ) + b;
//}
//
//



//void main() {
//    vTexCoord = texCoord;
//    vTileCoord = tileCoord;
//    vGridCoord = gridCoord;
//    vScroller = abs(scroller);
//    
//    vec4 pos = vec4(vertex,1);
//
//    pos.y += scroller;
//    pos.x += (tileCoord.x * tileSize.x);
//    pos.y += (tileCoord.y * tileSize.y);
//
//    pos.x += ( tileOffset.x * tileSize.x);
//    pos.y += ( tileOffset.y * tileSize.y);
//    
//    float foo = fun(now, 1., 1./3., 0., 0.);
//    
//    float f = abs( fun(now, 1., 1./8., 0., 0.) );
//    f += 0.0001;
//    f /= 8.;
//    pos.z += fun(pos.x+now, 1., 1./8., 0., 0.);
//    pos.z +=  fun(pos.y+now, 1., 1./8., 0., 0. );
//    
//    pos.x += 0.125 * ( fun(pos.x+now*8., 1., 1./8., 0., foo) );
//    pos.y += 0.125 * ( fun(pos.y+now*8., 1., 1./4., 0., 0.) );
//
//    pos.z += fun0( pos.x + now );
//    pos.z += fun0( pos.y + now );
//
//    pos.xyz = Wave(pos.xyz, vec2(0.,0.);
//
//
//
//    gl_Position = projection * view * model * pos;
//}

