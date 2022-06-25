uniform mat4 model;               // model transformation
uniform mat4 view;                // view transformation
uniform mat4 projection;          // projection transformation

uniform float screenRatio;

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

mat4 ident() {return mat4(1.,0.0,0.0,0.0,0.0,1.,0.0,0.0,0.0,0.0,1.,0.0,0.0,0.0,0.0,1.0);}
mat4 scale(float s) {return mat4(s,0.0,0.0,0.0,0.0,s,0.0,0.0,0.0,0.0,s,0.0,0.0,0.0,0.0,1.0);}


/*************************************/
/*  mat3 rotz(float w) {             */
/*      return mat3(                 */
/*         cos(w), sin(w), 0.0,      */
/*        -sin(w), cos(w), 0.0,      */
/*            0.0,    0.0, 0.0       */
/*      );                           */
/*  }                                */
/*  mat3 roty(float w) {             */
/*      return mat3(                 */
/*         cos(w), 0.0, sin(w),      */
/*            0.0, 1.0,    0.0,      */
/*        -sin(w), 0.0, cos(w)       */
/*      );                           */
/*  }                                */
/*  mat3 rotx(float w) {             */
/*      return mat3(                 */
/*            1.0,    0.0,    0.0,   */
/*            0.0, cos(w), sin(w),   */
/*            0.0,-sin(w), cos(w)    */
/*      );                           */
/*  }                                */
/*************************************/


float len(vec3 v) {
    return sqrt( v.x*v.x + v.y*v.y + v.z*v.z );
}



vec3 rotate(float w, vec3 v) {
    return rotz(w)*v;
}

vec3 translate(float w, float r, vec3 v) {
    v.x += cos(w)*r;
    v.y += sin(w)*r;
    return v;
}

vec3 wave(vec3 v) {
    float run = 0.0;
    run = 4.*now/1.;
    v.z += .25 * sin(run + v.x + 1.);
    return v;
}


vec3 scale(vec3 v) {
    float s = 2. * gridCoord.x / tileCount.x;

    if (mod(gridCoord.x,2.) == 1.0) {
//        return v;
    }

    float t = log(s+0.);
    t = s;
    
    v.x *= (1.+t);
    v.y *= (1.+t);
    
    return v;
}


vec3 shape(vec3 v) {
    float s = gridCoord.x / tileCount.x;
    v.xy *= 1.+8.*s;
    v.x += log(1.+2.*s) * gridCoord.x;
    v.x += (tileCoord.x * tileSize.x);
    return v;    
}

mat4 zoom(float radius) {
    float s = 1.0;
    s = 2./ ((tileCount.x * tileSize.x)+2.*radius);
    s *= screenRatio;
    mat4 ret = scale(s);
    return ret;
}



void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    mat4 mdl = ident();
    vec4 pos = vec4(vertex,1);

    DEBUG_FREEZE = true;
    DEBUG_TOP = true;


//    pos.y += .5;
//    pos.x += wordWidth/2.;

    float phi = 0.0;
    if (! DEBUG_FREEZE ) {
        phi += -now/TAU;
        phi += (PI/(2.*tileCount.y)) * -cos( PI * Ease(now) );
    }



    float ARC = TAU;
    float sector = ARC / (tileCount.y+4.);
    float gamma = (tileCount.y-gridCoord.y) * sector;



    float XXX = gridCoord.x / tileCount.x;
    float radius1 = 24.;

    phi += scroller * sector;


//    pos.y -= (scroller * tileSize.y);



//    pos.x += (tileOffset.x * tileSize.x);




    pos.xyz = shape(pos.xyz);

    pos.xyz = rotate(gamma+phi,pos.xyz);
    pos.xyz = translate(gamma+phi,radius1,pos.xyz);


    float rho = 0.0;
    if (! DEBUG_FREEZE ) {
        rho = now/5.;
    }
    if (! DEBUG_TOP ) {
        pos.xyz = roty( -PI/8. + Ease(rho)*PI/8.) * pos.xyz;
        pos.xyz = rotx( -PI/4. + Ease(rho+PI/3.)*-PI/8.) * pos.xyz;
    }

    mdl = zoom(radius1);
    gl_Position = projection * view  * mdl * pos;
}


