uniform sampler2D texture;

uniform vec2 tilesize;

varying vec2 fragcoord;
varying vec2 fraggridcoord;
varying vec2 offset;

void main() {

    float H = 8.;
    float W = 32.;

    vec2 tex = fragcoord;

    tex.x = fragcoord.x / W;
//    tex.x += 2. * fraggridcoord.x/W;
    tex.x += offset.x / W;    
    
    tex.y = fragcoord.y / H;
//    tex.y += 1. * fraggridcoord.y/H;
    tex.y += offset.y / H;    

    gl_FragColor = texture2D(texture, tex);

//    gl_FragColor = vec4( tex.x*10., tex.y*10., 0.0, 1.0 );
}
