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


mat4 rotationMatrix(vec3 axis, float angle)
{
    vec3 a  = normalize(axis);
    float s = sin(angle);
    float c = cos(angle);
    float oc = 1.0 - c;
    
    return mat4(
        oc*a.x*a.x + c,      oc*a.x*a.y - a.z*s,  oc*a.z*a.x + a.y*s,  0.0,
        oc*a.x*a.y + a.z*s,  oc*a.y*a.y + c,      oc*a.y*a.z - a.x*s,  0.0,
        oc*a.z*a.x - a.y*s,  oc*a.y*a.z + a.x*s,  oc*a.z*a.z + c,      0.0,
                       0.0,                 0.0,                 0.0,  1.0
    );
}




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
    
    


    
    float zoom;
    float ratio = screenRatio / fontRatio;
    
    float cols = ratio * 2. / tileCount.x;
    float rows = model[0][0];
    
//    float w = (tileCount.x/2.+tileCoord.x)/tileCount.x;
//    float h = (tileCount.y/2.+tileCoord.y)/tileCount.y;


    zoom = rows;

    float h0 = cols;
    float h1 = cols/8.;

//    float p = (tileCoord.y + tileCount.y/2.) / tileCount.y;

    float p = (pos.y + (tileSize.y*tileCount.y/2.)) / (tileSize.y * tileCount.y); 
    
    
    
    
    
  



//    pos.x *= h1;

    pos.x *= cols/(1.+4.*p);
    pos.y *= rows;


    
//    mat4 shear = mat4(
//        1.0, 0.0, 0.0, 0.0,
//        0.0, 1.0, 0.0, 0.0,
//        0.0, 0.0, 1.0, 0.0,
//        0.0, 0.0, 0.0, 1.0
//    
//    );

//    pos.y += ( (1.+(tileCount.y/2.)) * tileSize.y);
//    pos = shear * pos;
//    pos.y -= 5.;
    
//    float ALPHA;
//    ALPHA = PI * 2.5/8.;
//    mat4 rot;
//    
//
//    pos.y += ( (1.+(tileCount.y/2.)) * tileSize.y);
//
//    rot = rotationMatrix(vec3(1.,0.,0.), ALPHA);
//    pos = rot * pos;
//    pos.y -= 5.;
    
    mat4 mdl;
    mdl = mat4(1.0);
//    mdl[0][0] = zoom;
//    mdl[1][1] = zoom;

    gl_Position = projection * view * mdl * pos;
}

