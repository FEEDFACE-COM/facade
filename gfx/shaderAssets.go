
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




"def":`
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




"grid/adhoc":`
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

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    pos.y -= scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);


    float x = float( tileCoord.x );
    float y = float( tileCoord.y );
    bool a = mod(x,2.0) == 1.0 ;
    bool b = mod(y,2.0) == 1.0 ;
    
    float A = 0.5;
    if ( a == b) {
        pos.z +=  A * sin(now);
    } else {
//        pos.z -=  1. * sin(now + PI / 2. );
        pos.z -=  A * sin(now * 2. );
    }

    gl_Position = projection * view * model * pos;
}
`,




"grid/cylinder":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

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


//    float c = tileCount.x * tileSize.x;
    float c = tileCount.y * tileSize.y;
    

    // c = 2π * r <=> c/2π = r //
    float r = 1. * c / TAU;
    
    float a;
//    a = (tileCoord.x / (0.5*tileCount.x + 2.)) * PI - PI/8.;
    a = (tileCoord.y / (0.5*tileCount.y + 2.)) * PI - PI/8.;

    a += PI/4.;
//    a += now/10.;
//    a += ease1(now/2.);
    

    pos = rotationMatrix(vec3(1.,0.,0.), a*1.) * pos;
    

    r *= 1.1;
    pos.z +=  cos(a) * r;
    pos.y +=  sin(a) * r;
    pos.x += (tileCoord.x * tileSize.x);




    
    vec3 axis = vec3(-1.,0.,0.);
    mat4 rot = rotationMatrix(axis, -PI/2.);
    pos = rot * pos;


//    pos.xyz *= 2.5;
    
    gl_Position = projection * view * model * pos;
}
`,




"grid/def":`
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

float PI  = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU = 6.2831853071795864769252867665590057683943387987502116419498891840;

float Identity(float x) { return x; }
float EaseInEaseOut(float x) { return -0.5 * cos( x * PI ) + 0.5; }

float EaseOut(float x) { return cos(x*PI/2. + 3.*PI/2. ); }
float EaseIn(float x) { return  -1. * cos(x*PI/2. ) + 1.  ; }

void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vGridCoord = gridCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    float foo = Identity( scroller );

    pos.y += foo;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);

    gl_Position = projection * view * model * pos;
}

`,




"grid/disk":`
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


    float RADIUS = 10.;
    float R0 = 4.0;
    float rad = RADIUS / (tileCount.y + R0); 


    float delta = 0.0;
//    delta += now/10.;
    delta += ease1(now/2.) - 0.5;
    

    float ARC = TAU;
    float A0 = 2.0;
  
    float alpha,gamma;
    
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

    gl_Position = projection * view * model * pos;
}

`,




"grid/pipe":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;

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


    float c = tileCount.x * tileSize.x;

    // c = 2π * r <=> c/2π = r //
    float r = c / TAU;
    
    float a;
    a = (tileCoord.x / (0.5*tileCount.x + 2.)) * PI - PI/8.;

    a += now/10.;
    
    
    a += ease1(now/2.);

    pos = rotationMatrix(vec3(1.,0.,0.), PI/2.) * pos;
    pos = rotationMatrix(vec3(0.,0.,1.), -a-PI/2.) * pos;
    

    pos.x +=  cos(a) * r;
    pos.y +=  sin(a) * r;


    pos.z -= tileCoord.y;
    pos.z -= scroller;


    
    vec3 axis = vec3(-1.,-1.,0.);
    mat4 rot = rotationMatrix(axis, PI/2.);
    pos = rot * pos;

//    pos.x += (tileCoord.x * tileSize.x);
//    pos.y += (tileCoord.y * tileSize.y);
    
        
    

    
    gl_Position = projection * view * model * pos;
}
`,




"grid/plasma":`
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

uniform vec2 tileSize;
uniform vec2 tileCount;
uniform vec2 tileOffset;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

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
    vScroller = abs(scroller);
    vGridCoord = gridCoord;
    
    vec4 pos = vec4(vertex,1);

    pos.y += scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);

    float F = 0.75;
    float x = pos.x;
    float y = pos.y;
    
    float freq = -1./24.;
    pos.y += F * cos( 2. * freq * x * PI + now         );
    pos.x += F * cos( 3. * freq * y * PI + now + PI/2. );
	

    gl_Position = projection * view * model * pos;
}

`,




"grid/zstep":`
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

    float F = 1.;

    float y  =  vTileCoord.y       / (tileCount.y/2.);
    float yy = (vTileCoord.y + ((scroller)) ) / (tileCount.y/2.);


    float freq = -1.;
    float f0 = cos( freq * y  * PI + now + PI/2. );
    float f1 = cos( freq * yy * PI + now + PI/2. );
    float d =  f0 + /*(scroller) * */(f1 - f0);
    pos.z += F * d;
    pos.z -= F;


    
    gl_Position = projection * view * model * pos;
}
`,




"grid/zwave":`
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

    float F = 0.5;
    
    pos.z += F * cos( pos.x + 2. * now         );
    pos.z += F * cos( pos.y + 3. * now + PI/2. );
    pos.z -= 2. * F;
    
    gl_Position = projection * view * model * pos;
}

`,




"mask/def":`

uniform float debugFlag;

attribute vec3 vertex;
attribute vec2 texCoord;

varying vec2 vTexCoord;

bool DEBUG = debugFlag > 0.0;


void main() {
    vTexCoord = texCoord;
    gl_Position = vec4(vertex,1.);
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




"def":`
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

uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;
uniform vec2 cursorPos;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying float vScroller;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    float F = 1.;

    float x = vGridCoord.x / tileCount.x;
    float y = vGridCoord.y / tileCount.y;
    
    col.r *= F * (1. - x);
    col.g *= F * (1. - y);

    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) {
        col.rgba = 1. - col.rgba;
    }

    if (gl_FrontFacing) { 
        vec3 tmp = vec3(col.rgb);
        col.r = tmp.g;
        col.g = tmp.b;
        col.b = tmp.r;
    }

    gl_FragColor = col;
    
}
`,




"grid/def":`

uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;
uniform vec2 cursorPos;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying float vScroller;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    if ( vGridCoord.y == 0.0 ) { // oldest line
        col.a *= (1.-vScroller);
    }
    
    if ( vGridCoord.y == tileCount.y ) { // newest line
        col.a *= vScroller;    
    }
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) { // invert cursor
        col.rgba = 1. - col.rgba;
    }

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
`,




"grid/fader":`

uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;
uniform vec2 cursorPos;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGridCoord;
varying float vScroller;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    if ( vGridCoord.y == 0.0 ) { // oldest line
        col.a *= (1.-vScroller);
    }
    
    if ( vGridCoord.y == tileCount.y ) { // newest line
        col.a *= vScroller;    
    }
    
    if ( cursorPos.x == vGridCoord.x && cursorPos.y == vGridCoord.y ) {
        col.rgba = 1. - col.rgba;
    }

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
`,




"mask/debug":`

uniform float debugFlag;
uniform float ratio;

varying vec2 vTexCoord;


bool DEBUG = debugFlag > 0.0;
float w = 0.002;

bool major(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.5) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}

bool minor(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.25) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}



void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    float a = 1.0;


    vec2 pos = vTexCoord;


    if ( true && major(pos) ) {
        float gray = 0.9;
        col = gray * vec3(1.,1.,1.);
    } else if ( true && minor(pos) ) {
        float gray = 0.5;
        col = gray * vec3(1.,1.,1.);
    }

    gl_FragColor = vec4(col.rgb, a);
}



`,




"mask/def":`

uniform float debugFlag;
uniform float ratio;

varying vec2 vTexCoord;


bool DEBUG = debugFlag > 0.0;
float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.5) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}


void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    float a = 1.0;


    vec2 pos = vTexCoord;


    if ( true && grid(pos) ) {
        float gray = 0.5;
        col = gray * vec3(1.,1.,1.);
    }

    
    
    gl_FragColor = vec4(col.rgb, a);
}


`,




"mask/mask":`

uniform float debugFlag;
uniform float ratio;

varying vec2 vTexCoord;

bool DEBUG = debugFlag > 0.0;

float MAX(float a, float b) { 
    if (a>=b) 
        return a; 
    else 
        return b; 
}

float mask(vec2 p) { 
    float x = p.x; float y = p.y;
    float ff = 1.;
    float A = 7. * ff;
    float B = 8. * ff; 
    if (abs(x) >= A/B && abs(y) >= A/B)
        return 1.0 - ( B * MAX(abs(x),abs(y)) - A );
    else if (abs(x) >= A/B) 
        return 1.0 - ( B * abs(x) - A);
    else if (abs(y) >= A/B)
        return 1.0 - ( B * abs(y) - A);
    return 1.0;
}



void main() {
    vec3 col = vec3(0.0,0.0,0.0);
    float a = 1.;


    vec2 pos = vTexCoord;

    a = mask(vec2(pos.x/ratio,pos.y));

    if (DEBUG) {
        col.rgb =  vec3(  1. - mask(vec2(pos.x/ratio,pos.y)) );
    }
    
    
    gl_FragColor = vec4(col.rgb, a);
}
`,


}
