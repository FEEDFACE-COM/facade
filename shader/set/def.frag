uniform float debugFlag;
uniform sampler2D texture;

varying vec2 vTexCoord;

bool DEBUG    = debugFlag > 0.0;



void main() {
    vec4 col;
    col = texture2D(texture, vTexCoord); 

    if (DEBUG) { 
        col.rgb = vec3(1.,1.,1.);
        col.a = 1.0;
    } 

    if (!gl_FrontFacing) { col.a /= 4.; }

    gl_FragColor = col;
    
}
