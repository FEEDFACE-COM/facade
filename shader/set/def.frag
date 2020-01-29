
uniform float now;
uniform float debugFlag;
uniform sampler2D texture;

uniform float tagCount;
uniform float tagMaxWidth;
uniform float tagFader;
uniform float tagIndex;


varying vec2 vTexCoord;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if ( col.a < 1. ) {
        col.a = 1.0;
        col.rgb = 0.25  * vec3(1.,1.,1.);
    }


    if (DEBUG) { 
        col.r = 1.;
//        col.rgb = vec3(1.,1.,1.);
//        col.g = tagIndex / tagCount;
        col.a = 1.0;
    } 

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
