
// +build linux,arm
package gfx
var FragmentShader = map[string]string{


"color":`

varying vec4 fragcolor;

void main() {
    gl_FragColor = fragcolor;
}
`,




"grid":`
uniform sampler2D texture;


varying vec2 vTexCoord;
varying vec2 vTileCoord;


bool debug = false;

void main() {

    vec2 pos = vTileCoord;
    vec2 tex = vTexCoord;


    vec4 col = texture2D(texture, tex);

    
    if (false && pos.x == 0.0 && pos.y == 0.0 ) {
        col.r += 0.5;
    }
    
    gl_FragColor = vec4(col.rgb,1.0);
}
`,




"ident":`
uniform sampler2D texture;

varying vec2 fragcoord;

void main() {
    vec4 tex = texture2D(texture,fragcoord);
    gl_FragColor = tex;
}
`,




"mask":`
varying vec4 fragcolor;
varying vec2 fragcoord;

float w = 0.005;

bool grid(vec2 pos) {

    for (float d = -2.0; d<=2.0; d+=0.5) {
        if (abs(pos.y - d) - w <= 0.0 ) { return true; }
        if (abs(pos.x - d) - w <= 0.0 ) { return true; }
    }
    
    return false;
}

void main() {
    vec4 col = vec4(0.0,0.0,0.0,0.0);
    vec2 pos = fragcoord;
    
    if ( grid(pos) ) { col = vec4(1.,1.,1.,0.5); }
    
    if ( pos.y > 0.0 && pos.y < 1.0 && abs(pos.x) <= w ) { col = vec4(0.,1.,0.,1.); }
    if ( pos.x > 0.0 && pos.x < 1.0 && abs(pos.y) <= w ) { col = vec4(1.,0.,0.,1.); }

    gl_FragColor = col;

}
`,


}
