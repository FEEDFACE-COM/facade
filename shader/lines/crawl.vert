uniform mat4 model;               // model transformation
uniform mat4 view;                // view transformation
uniform mat4 projection;          // projection transformation

uniform vec2 tileSize;            // size of largest glyph in font, as (width,height)
uniform vec2 tileCount;           // grid dimensions, as (columns,rows)
uniform vec2 tileOffset;          // grid center offset from (0,0), as (columns,rows)

uniform float now;                // time since start of program, as seconds
uniform float scroller;           // scroller amount, from 0.0 to 1.0
uniform float debugFlag;          // 0.0 unless -D flag given by user

attribute vec3 vertex;            // vertex position as (x,y,z) centered on (0,0,0)

attribute vec2 tileCoord;         // tile coordinates, as (x,y) centered on (0,0)
attribute vec2 gridCoord;         // tile coordinates, as (column,row)
attribute vec2 texCoord;          // texture coordinates in font texture atlas


varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying vec2 vTexCoord;


bool DEBUG = debugFlag > 0.0;
bool DEBUG_FREEZE = false;
bool DEBUG_TOP = false;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }
float Ease(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

mat3 rotx(float w) {return mat3(1.0,0.0,0.0,0.0,cos(w),sin(w),0.0,-sin(w),cos(w));}
mat3 roty(float w) {return mat3(cos(w),0.0,sin(w),0.0,1.0,0.0,-sin(w),0.0,cos(w));}
mat3 rotz(float w) {return mat3(cos(w),sin(w),0.0,-sin(w),cos(w),0.0,0.0,0.0,0.0);}



void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    
    vec4 pos = vec4(vertex,1);

    pos.y += (scroller * tileSize.y);
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    pos.y += 7.;

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);

    pos.z -= 2.*pos.y;

    gl_Position = projection * view * model * pos;
}



