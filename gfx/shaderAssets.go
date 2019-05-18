
package gfx
var VertexShader = map[string]string{


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




"grid/drain":`
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
    pos.z += scroller;


    
    vec3 axis = vec3(-1.,-1.,0.);
    mat4 rot = rotationMatrix(axis, PI/2.);
    pos = rot * pos;

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
uniform vec2 tileOffset;

uniform float now;
uniform float scroller;
uniform float debugFlag;
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vScroller;


bool DEBUG = debugFlag > 0.0;


bool oddColCount() { return mod(tileCount.x, 2.0) == 1.0 ; }
bool oddRowCount() { return mod(tileCount.y, 2.0) == 1.0 ; }


float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;


void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    pos.y -= scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);
    
    float F = 0.5;
    float t = -1.;
    if (downward == 1.0) {
        t = 1.;
    }
    float y  =  vTileCoord.y       / (tileCount.y/2.);
    float yy = (vTileCoord.y + t ) / (tileCount.y/2.);

    float freq = 2.;
    float f0 = cos( freq * y  * PI + now + PI/2. );
    float f1 = cos( freq * yy * PI + now + PI/2. );
    float d =  f0 + abs(scroller) * (f1 - f0);
    pos.z += F * d;


    
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

    gl_Position = projection * view * model * pos;
}
`,




"grid/disk":`
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

varying float vScroller;

bool DEBUG = debugFlag > 0.0;


float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;
float TAU= 6.2831853071795864769252867665590057683943387987502116419498891840;
float ease1(float x)          { return 0.5 * cos(     x + PI/2.0 ) + 0.5; }

bool oddColCount() { return mod(tileCount.x, 2.0) == 1.0 ; }
bool oddRowCount() { return mod(tileCount.y, 2.0) == 1.0 ; }


bool upperVertex(vec2 vec) {
    return vec.y >= 0.; 
}


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
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);


    float rx = tileCoord.x / tileCount.x;

    float r0 = 4.;
    float r1 = 1.;
    
    float phi = (tileCoord.x / tileCount.x) * TAU;
    
    if ( upperVertex(pos.xy) ) {
        pos.x += cos(phi) * r0;
    } else {
        pos.x += cos(phi) * r1;
    }
//    if (pos.y > 0.) {
//        pos.x = cos(phi) * r0;
//    } else {
//        pos.y = cos(phi) * r1;
//    }


    vec3 axis = vec3(-1.,-1.,0.);
    mat4 rot = rotationMatrix(axis, PI/4.);
//    pos = rot * pos;
    
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
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;

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
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);


//    float c = tileCount.x * tileSize.x;
    float c = tileCount.y * tileSize.y;
    

    // c = 2π * r <=> c/2π = r //
    float r = c / TAU;
    
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
uniform float downward;

attribute vec3 vertex;
attribute vec2 texCoord;
attribute vec2 tileCoord;


varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vScroller;

bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;


void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    pos.y -= scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);

    float F = 0.25;
    
    pos.z += F * cos( pos.x + 2. * now         );
    pos.z += F * cos( pos.y + 3. * now + PI/2. );
	

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


varying vec2 vTexCoord;
varying vec2 vTileCoord;

varying float vScroller;

bool DEBUG = debugFlag > 0.0;

float PI = 3.1415926535897932384626433832795028841971693993751058209749445920;


void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    vScroller = abs(scroller);
    
    vec4 pos = vec4(vertex,1);

    pos.y -= scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    pos.x += ( tileOffset.x * tileSize.x);
    pos.y += ( tileOffset.y * tileSize.y);

    float F = 0.75;
//    float F = 0.5;
//    float x = vTileCoord.x / (tileCount.x/2.);
//    float y = vTileCoord.y / (tileCount.y/2.);
    float x = pos.x;
    float y = pos.y;
    
    float freq = 1./24.;
//    float freq = 1./8.;
    pos.y += F * cos( 2. * freq * x * PI + now         );
    pos.x += F * cos( 3. * freq * y * PI + now + PI/2. );
	

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
    vScroller = abs(scroller);
    
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
    float zoom = 2.;
    a = zoom * radius;
    b = zoom * radius/2.;

    

    float phase = -1. * PI/8.; 
    
    
//    phase += PI/2. * (ease1( now/8. ) );
    
    float alpha = grad.x * TAU + phase;
 
//    alpha += now/10.;
 
    
    alpha = alpha * -1. + PI/4.;
        
    pos = rotationMatrix(vec3(0., 0., 1.), alpha) * pos;


    

//    a += 0.25 * cos(0.3*now * ease1((  mod(now,TAU) +scroller)/100.) * vTileCoord.y);

//    vec3 tmp = vec3(pos);
    

    pos.x -= cos(alpha) * a;
    pos.y -= sin(alpha) * b;

    pos.z += (tileCoord.y*tileSize.y);
    pos.z += scroller;

//    pos = rotationMatrix(vec3(-1.,1.,0.), PI/4.) * pos;
//    pos = rotationMatrix(vec3(1.,0.,0.), PI/2.) * pos;

    
    gl_Position = projection * view * model * pos;
}
`,


}


var FragmentShader = map[string]string{


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




"color":`


varying float vDebugFlag;
varying vec4 vFragColor;

bool DEBUG = vDebugFlag > 0.0;

void main() {
    gl_FragColor = vFragColor;
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
    float A = 3. * ff;
    float B = 4. * ff; 
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
    float a = 1.0;


    vec2 pos = vTexCoord;

    a = mask(vec2(pos.x/ratio,pos.y));

    if (DEBUG) {
        col.b =    1. - mask(vec2(pos.x/ratio,pos.y)) ;
    }
    
    
    gl_FragColor = vec4(col.rgb, a);
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




"grid/fader":`

uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;



varying vec2 vTexCoord;
varying vec2 vTileCoord;



bool DEBUG    = debugFlag > 0.0;
bool DOWNWARD = downward == 1.0;

bool firstLine() {
    float t = 0.0;
    if (DOWNWARD) {
        t = 1.;
    }    
    return (  tileCount.y + vTileCoord.y - 1.) * 2. + t <= tileCount.y;
}


bool lastLine() { 
    float t = 0.0;
    if (DOWNWARD) {
        t = 1.0;
    } 
    return  vTileCoord.y*2.  > tileCount.y - t;
}

bool newestLine() {
	return  ( !DOWNWARD && firstLine() ) || ( DOWNWARD && lastLine() ) ;
}

bool oldestLine() {
	return ( !DOWNWARD && lastLine() ) || ( DOWNWARD && firstLine() );
}


void main() {
    vec4 col;
    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    } else { 
        col = texture2D(texture, vTexCoord);
    }

	if ( newestLine() ) {
		col.a *= (1.0 - scroller);
	} else if ( oldestLine() ) {
		col.a *= scroller;
	}


    if (!gl_FrontFacing) {
		col.a /= 2.;
	}

    gl_FragColor = col;

}

`,




"grid/debug":`

uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;

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

    if ( vGridCoord.y == 0.0 ) {
        col.rb *= 0.0;
    }
    
    if ( vGridCoord.y == tileCount.y ) {
        col.gb *= 0.0;    
    }

/*    
//    if (firstLine()) { col.r *= 0.; } //cyan
//    if (lastLine())  { col.g *= 0.; } //magenta
    
    
    if ( newestLine() )   { col.rb *= 0.; col.a *= (1.-vScroller); } //green
    if ( oldestLine() ) { col.gb *= 0.; col.a *= vScroller; }      //red
*/
    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
`,




"grid/def":`

uniform float debugFlag;
uniform float downward;
uniform vec2 tileCount;
uniform sampler2D texture;
uniform float scroller;



varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying float vScroller;



bool DEBUG    = debugFlag > 0.0;
bool DOWNWARD = downward == 1.0;
bool ODDLINES = mod(tileCount.y, 2.0) == 1.0;


bool firstLine() {
    float s = 0.0;
    float t = 0.0;
    if ( ! DOWNWARD && ! ODDLINES) { s = 0.0; t = 0.0; }
    if ( ! DOWNWARD &&   ODDLINES) { s = 1.0; t = 0.0; }
    if (   DOWNWARD && ! ODDLINES) { s = 0.0; t = 1.0; }
    if (   DOWNWARD &&   ODDLINES) { s = 0.0; t = 1.0; }
    return (  tileCount.y + (vTileCoord.y-s) - 1.) * 2. + t <= tileCount.y;
}


bool lastLine() { 
    float s = 0.0;
    float t = 0.0;
    if ( ! DOWNWARD && ! ODDLINES) { s = 0.0; t = 0.0; }
    if ( ! DOWNWARD &&   ODDLINES) { s = 1.0; t = 0.0; }
    if (   DOWNWARD && ! ODDLINES) { s = 0.0; t = 1.0; }
    if (   DOWNWARD &&   ODDLINES) { s = 0.0; t = 1.0; }
    return  (vTileCoord.y-s)*2.  > tileCount.y - t;
}

bool newestLine() {
    return  ( !DOWNWARD && firstLine() ) || ( DOWNWARD && lastLine() ) ;
}

bool oldestLine() {
    return ( !DOWNWARD && lastLine() ) || ( DOWNWARD && firstLine() );
}


void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord);

    if (DEBUG) { 
        col = vec4(1.,1.,1.,1.); 
    }

	if ( newestLine() ) { col.a *= (1.0 - vScroller); }
	if ( oldestLine() ) { col.a *= vScroller; }

    if (!gl_FrontFacing) { col.a /= 2.; }

    gl_FragColor = col;

}

`,


}
