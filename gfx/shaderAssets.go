
package gfx
var VertexShader = map[string]string{


"color":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform float debugFlag;

attribute vec3 vertex;
attribute vec4 color;

varying vec4  vFragColor;
varying float vDebugFlag;


void main() {
    vFragColor = color;
    vDebugFlag = debugFlag;
    gl_Position = projection * view * model * vec4(vertex, 1);
}
`,




"identity":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vFragCoord;
varying float vDebugFlag;

bool DEBUG = debugFlag > 0.0;

void main() {
    vFragCoord = texCoord;
    gl_Position = projection * view * model * vec4(vertex, 1);
}
`,




"grid/identity":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;
varying float vScroller;
varying float vDownward;


bool DEBUG = debugFlag > 0.0;

void main() {
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vNow = now;
    vScroller = scroller;
    vDownward = downward;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    //
    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    gl_Position = projection * view * model * pos;
}
`,




"grid/oval":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vScroller;
varying float vDebugFlag;
varying float vNow;
varying float vDownward;


bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU= 6.2831853071795864769252867665590057683943387987502116419498891840;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

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
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vScroller = scroller;
    vNow = now;
    vDownward = downward;
    
    vec4 pos = vec4(vertex,1);
    
//    pos = rotationMatrix(vec3(1., 0., 0.), PI  ) * pos;
    
    pos.z = pos.y;
    pos.y = 0.;



    vec2 w = vec2(tileSize.x*tileCount.x, tileSize.y*tileCount.y);

    vec2 coord = vec2( 
        tileCoord.x + (tileCount.x/2.  - 1.),
        tileCoord.y + (tileCount.y/1.  - 1.)
    );
    
    vec2 grad = vec2(coord.x / tileCount.x, coord.y / tileCount.y);

    // circum = 2π * radius <=> circum/2π = radius //
    float circum = tileCount.x * tileSize.x;
    float radius = circum / TAU;
    
    float a,b;
    
    a = radius;
    b = radius/2.;

    

    float phase = -1. * PI/8.; 
    
//    phase += PI/2. * (ease1( now/8. ) );
    
    float alpha = grad.x * TAU + phase;
    
    pos = rotationMatrix(vec3(0., 0., 1.), alpha * -1. + PI/2.  ) * pos;


    

    a += 0.25 * cos(0.3*now * ease1((  mod(now,TAU) +scroller)/100.) * vTileCoord.y);

//    vec3 tmp = vec3(pos);
    
//    pos.x = cos(alpha) * 1.;
//    pos.y = sin(alpha) * 1.;

    pos.x -= cos(alpha) * a;
    pos.y -= sin(alpha) * b;

    pos.z -= (tileCoord.y*tileSize.y);
    pos.z -= scroller;

//    pos = rotationMatrix(vec3(-1.,1.,0.), PI/4.) * pos;
//    pos = rotationMatrix(vec3(1.,0.,0.), PI/2.) * pos;

    
    gl_Position = projection * view * model * pos;
}
`,




"grid/tunnel":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vScroller;
varying float vDebugFlag;
varying float vNow;
varying float vDownward;


bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU= 6.2831853071795864769252867665590057683943387987502116419498891840;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }


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
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vScroller = scroller;
    vNow = now;
    vDownward = downward;
    
    vec4 pos = vec4(vertex,1);

//    pos.y += scroller;

    float c = tileCount.x * tileSize.x;

    // c = 2π * r <=> c/2π = r //
    float r = c / TAU;
    
    float a;
    a = (tileCoord.x / (0.5*tileCount.x + 2.)) * PI - PI/8.;

    pos = rotationMatrix(vec3(1.,0.,0.), PI/2.) * pos;
    pos = rotationMatrix(vec3(0.,0.,1.), -a-PI/2.) * pos;
    

    pos.x +=  cos(a) * r;
    pos.y +=  sin(a) * r;


    pos.z -= tileCoord.y;
    pos.z -= scroller;


    
//    pos.z += 2.;
//    if (downward == 1.0) {
//        pos.z += tileCoord.y;
//        pos.z -= scroller;
//    } else {
//    }
    
//    pos.y += tileCoord.y * tileSize.y;


//    vec3 axis = vec3(1.,0.,0.);
//    mat4 rot = rotationMatrix(axis, PI/2.);
//    pos = rot * pos;

//    pos.x += (tileCoord.x * tileSize.x);
//    pos.y += (tileCoord.y * tileSize.y);
    
        
    

    
    gl_Position = projection * view * model * pos;
}
`,




"grid/zstep":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;
varying float vScroller;
varying float vDownward;

bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }


void main() {
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vDownward = downward;
    vScroller = scroller;
    vNow = now;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    
    
    float F = 0.25;
    float f0 = ease1( now /1.);  
    float f1 = 0.;
    
    

    // allow for scroller
    float from = cos( vTileCoord.y + 3. * now + PI/2.);
    float to =   cos( vTileCoord.y-1. + 3. * now + PI/2. );
//    float from = to;
    float delta =  to + scroller * (from - to);
    
    

//    pos.z += F * cos( vTileCoord.x + 2. * now );
    pos.z += F * delta;


    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    
    gl_Position = projection * view * model * pos;
}
`,




"grid/zwave":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;


varying float vDebugFlag;
varying float vNow;
varying float vScroller;
varying float vDownward;

bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }


void main() {
    vTileCount = tileCount;
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vDebugFlag = debugFlag;
    vScroller = scroller;
    vNow = now;
    vDownward = downward;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);
    
    
    float F = 0.25;
    float f0 = ease1( pos.y + now );  
    float f1 = 0.;
    
    
    pos.z += F * cos( pos.x + 2. * now         );
    pos.z += F * cos( pos.y + 3. * now + PI/2. );
    pos.z += F * f0;


    if (mod(tileCount.x, 2.0) != 1.0 ) { pos.x -= tileSize.x/2.; }
    if (mod(tileCount.y, 2.0) != 1.0 ) { pos.y -= tileSize.y/2.; }

    
    gl_Position = projection * view * model * pos;
}
`,




"mask/identity":`


uniform float debugFlag;
uniform float ratio;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vTexCoord;
varying float vRatio;
varying float vDebugFlag;

bool DEBUG = debugFlag > 0.0;


void main() {
    vTexCoord = texCoord;
    vRatio = ratio;
    vDebugFlag = debugFlag;
    
    
    gl_Position = vec4(vertex,1);
}
`,


}


var FragmentShader = map[string]string{


"color":`


varying float vDebugFlag;
varying vec4 vFragColor;

bool DEBUG = vDebugFlag > 0.0;

void main() {
    gl_FragColor = vFragColor;
}
`,




"identity":`
uniform sampler2D texture;

varying vec2 vFragCoord;
varying float vDebugFlag;

bool DEBUG = vDebugFlag > 0.0;


void main() {
    vec4 tex = texture2D(texture,vFragCoord);
    gl_FragColor = tex;
}
`,




"grid/debug":`

uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;
varying float vScroller;
varying float vDownward;

bool DEBUG    = vDebugFlag > 0.0;
bool downward = vDownward == 1.0;

bool firstLine() {
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        return 0.5*vTileCount.y       == vTileCoord.y ;
    } else {
        return  0.5*(vTileCount.y+1.) == vTileCoord.y + 1. ;
    }
}

bool lastLine() {
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        return -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
    } else {
        return -0.5*(vTileCount.y+1.) == vTileCoord.y - 1. ;
    }

}


void main() {

    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { 
        col = texture2D(texture, vTexCoord); 
    }
    
    col.b =  1.-  abs((2.*vTileCoord.x /vTileCount.x));
    
    if ( ! downward && firstLine() || downward && lastLine() ) {
        col.g = 0.0;
    }

    if ( ! downward && lastLine() || downward && firstLine() ) {
        col.r = 0.0;
        col.b = 0.0;
    }
    
    
    gl_FragColor = col;
}
`,




"grid/fader":`

uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;
varying float vScroller;
varying float vDownward;

bool DEBUG    = vDebugFlag > 0.0;
bool downward = vDownward == 1.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;


bool firstLine() {
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        return 0.5*vTileCount.y       == vTileCoord.y ;
    } else {
        return  0.5*(vTileCount.y+1.) == vTileCoord.y + 1. ;
    }
}

bool lastLine() {
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        return -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
    } else {
        return -0.5*(vTileCount.y+1.) == vTileCoord.y - 1. ;
    }

}



void main() {

    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { 
        col = texture2D(texture, vTexCoord); 
    }
    
    
    
    if (downward) {
        if ( firstLine() ) {
            col.rgb *= (-1.0 * vScroller);
        }
        if ( lastLine() ) {
            col.rgb *= (1.0 - -1.0 * vScroller);
        }
    } else { // ! downward
        if ( firstLine() ) {
            col.rgb *= (1.0 - vScroller);
        }
        if ( lastLine() ) {
            col.rgb *= vScroller; 
        }
    }
        
    gl_FragColor = vec4(col.rgb,1.0);
}
`,




"grid/identity":`

uniform sampler2D texture;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vTileCount;

varying float vDebugFlag;
varying float vNow;

bool DEBUG    = vDebugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;



void main() {

    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { col = 
        texture2D(texture, vTexCoord); 
    }
    
    bool firstLine,lastLine;
    if (mod(vTileCount.y, 2.0) != 1.0 ) { 
        firstLine =  0.5*vTileCount.y       == vTileCoord.y ;
        lastLine  = -0.5*vTileCount.y + 1.0 == vTileCoord.y ;
    } else {
        firstLine =  0.5*(vTileCount.y+1.) == vTileCoord.y + 1. ;
        lastLine  = -0.5*(vTileCount.y+1.) == vTileCoord.y - 1. ;
    }

    if (DEBUG && firstLine ) {
        col.rgb = vec3(1.,0.,1.);
    }

    if (DEBUG && lastLine ) {
        col.rgb = vec3(0.,1.,0.);
    }
    
    gl_FragColor = col;
}
`,




"mask/debug":`

varying vec2 vTexCoord;
varying float vDebugFlag;
varying float vRatio;


bool DEBUG = vDebugFlag > 0.0;
float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.5) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}


void main() {

    float MARGIN = 0.5;

    vec2 pos = vTexCoord;
    vec3 col = vec3(1.0,1.0,1.0);

//    if ( grid(pos) ) { col = vec4(1.,1.,1.,0.5); }
//    if ( pos.y > 0.0 && pos.y < 1.0 && abs(pos.x) <= w ) { col = vec4(0.,1.,0.); }
//    if ( pos.x > 0.0 && pos.x < 1.0 && abs(pos.y) <= w ) { col = vec3(1.,0.,0.); }
  
  

    if (pos.x > MARGIN || pos.x < -1. * MARGIN  ) {
        col -= 0.5;
    }
    
    if (pos.y > MARGIN || pos.y < -1. * MARGIN  ) {
        col -= 0.5;
    }
    
       
    gl_FragColor = vec4(col.rgb, 1.0);
}


`,




"mask/identity":`

varying vec2 vTexCoord;
varying float vDebugFlag;

bool DEBUG = vDebugFlag > 0.0;


void main() {

    vec2 pos = vTexCoord;
    vec4 col = vec4(0.0,0.0,0.0,0.0);

    gl_FragColor = col;
}
`,


}
