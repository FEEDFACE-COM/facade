uniform sampler2D texture;

uniform vec2 tileSize;
uniform vec2 tileCount;

uniform vec2 glyphCount;
uniform vec2 glyphSize;

varying vec2 vTexCoord;
varying vec2 vTileCoord;
varying vec2 vGlyphCoord;
varying vec2 vGlyphSize;


float clamp(float x) { 
    if ( x < 0.0 ) return 0.0;
    if ( x > 1.0 ) return 1.0;
    return x;
}

void main() {

    float H = 9.;
    float W = 18.;
//
//    vec2 grid = tileCount;
    vec2 pos = vTileCoord;

    vec2 tex = vTexCoord;

    vec2 coord = vGlyphCoord;


    vec4 col = texture2D(texture, tex);

    
//    vec4 col = vec4(0.,1.,0., 1.);

//    if ( mod(pos.y,2.0) != 1. ) { col.r += 0.5; }
//    if ( mod(pos.x, 2.0) == 1. ) { col.g += 0.5; }
//        
//
//    if (pos.x == 0.0 && pos.y == 0.0 ) {
//        col.rgb = vec3(1,1,1);    
//    }
    

    gl_FragColor = vec4(col.rgb,1.0);

//    gl_FragColor = vec4( tex.x*10., tex.y*10., 0.0, 1.0 );
}
