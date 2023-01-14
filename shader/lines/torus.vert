uniform mat4 model;               // model transformation
uniform mat4 view;                // view transformation
uniform mat4 projection;          // projection transformation

uniform vec2 tileSize;            // size of largest glyph in font, as (width,height)
uniform vec2 tileCount;           // grid dimensions, as (columns,rows)
uniform vec2 tileOffset;          // grid center offset from (0,0), as (columns,rows)

uniform float now;                // time since start of program, as seconds
uniform float scroller;           // scroller amount, from 0.0 to 1.0
uniform float debugFlag;          // 0.0 unless -D flag given by user

uniform float screenRatio;
uniform float fontRatio;

attribute vec3 vertex;            // vertex position as (x,y,z) centered on (0,0,0)

attribute vec2 tileCoord;         // tile coordinates, as (x,y) centered on (0,0)
attribute vec2 gridCoord;         // tile coordinates, as (column,row)
attribute vec2 texCoord;          // texture coordinates in font texture atlas


varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying vec2 vTexCoord;



bool DEBUG = debugFlag > 0.0;

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }
float Ease(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

mat3 rotx(float w) {return mat3(1.0,0.0,0.0,0.0,cos(w),sin(w),0.0,-sin(w),cos(w));}
mat3 roty(float w) {return mat3(cos(w),0.0,sin(w),0.0,1.0,0.0,-sin(w),0.0,cos(w));}
mat3 rotz(float w) {return mat3(cos(w),sin(w),0.0,-sin(w),cos(w),0.0,0.0,0.0,1.0);}


vec3 sphere(vec2 count, vec2 index) {
    float RADIUS = 8.;

    vec3 ret;
    
    float count_theta = count.x;
    float idx_theta   = index.x;

    float count_phi = count.y;
    float idx_phi   = count.y-index.y;

    float theta = TAU * (idx_theta/count_theta);
    float phi  =  PI *  (idx_phi/count_phi);

    float slice_theta = (TAU / count_theta) / 2.;
    float slice_phi =   (PI / count_phi) / 2.;

    float t = theta;
    float p = phi;

//    p += (scroller * slice_phi);


    t = vertex.x >= 0. ? theta + slice_theta : theta - slice_theta;
    p = vertex.y >= 0. ?   phi + slice_phi   :   phi - slice_phi;


    ret.x = RADIUS * cos(t) * sin(p);
    ret.y = - RADIUS * sin(t) * sin(p);
    ret.z = RADIUS * cos(p);

    
    ret.xyz *= roty( PI );
    return ret;

}




vec3 torus(vec2 count, vec2 index) {
    float RADIUS = 4.;
    float THICKNESS = 3. * RADIUS;
    vec3 ret;
    
    float count_theta = count.x;
    float idx_theta   = - mod( index.x, count.x);

    float count_phi = count.y + 1.;
    float idx_phi   = mod(index.y,count.y);
//    float idx_phi   = mod( count_phi/2. + index.y - 1., count.y);
////    float idx_phi   = mod( -index.y + count.y/2. - 1., count.y);
//    
    
//    float theta = PI/2. + TAU * (idx_theta/count_theta);
    float theta = TAU * (idx_theta/count_theta);
    float phi  =  TAU * (idx_phi/count_phi);
    
    float slice_theta = PI / count_theta;
    float slice_phi =   PI / count_phi;
    


    float p = vertex.x >= 0. ? theta - slice_theta : theta + slice_theta;
    float t = vertex.y >= 0. ? phi   - slice_phi   : phi   + slice_phi;
    
    t -= 4.*PI/5.;
    
    t -= scroller * 2.*slice_phi;

    ret.x = (RADIUS * cos(t) + THICKNESS) * sin(p);
    ret.y = - (RADIUS * cos(t) + THICKNESS) * cos(p);
    ret.z =  RADIUS * sin(t);




    
    return ret;
}

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;

    
    vec2 count = vec2(
        tileCount.x, // rows
        tileCount.y // cols
    );

    vec2 index = vec2(
        gridCoord.x, // row
        gridCoord.y // col
    );
    
    vec4 pos = vec4(vertex,1.);

    
    if ( true ) {    
        pos.xyz = torus(count,index);
    } else {
        pos.xyz = sphere(count,index);
    }
    




//    pos.xyz = roty( PI/4.) * pos.xyz;
    float rho = now;
    if ( true ) {
        pos.xyz = rotz( -2.*PI/3. + now/16. ) * pos.xyz;
        pos.xyz = roty( -PI/8. + Ease(rho)*PI/8.) * pos.xyz;
        pos.xyz = rotx( -PI/4. + Ease(rho+PI/3.)*-PI/8.) * pos.xyz;
    
    } else {
    
        //pos.xyz 
        
    }



    gl_Position = projection * view * model * pos;
}

