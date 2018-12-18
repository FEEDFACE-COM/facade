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


bool DEBUG = debugFlag > 0.0;


bool oddColCount() { return mod(tileCount.x, 2.0) == 1.0 ; }
bool oddRowCount() { return mod(tileCount.y, 2.0) == 1.0 ; }




void main() {
    vTexCoord = texCoord;
    vTileCoord = tileCoord;
    
    vec4 pos = vec4(vertex,1);

    pos.y -= scroller;
    pos.x += (tileCoord.x * tileSize.x);
    pos.y += (tileCoord.y * tileSize.y);

    if ( oddColCount() ) {
    	pos.x += (-1.0 * tileSize.x);
    } else {
    	pos.x += ( 0.5 * tileSize.x);
   	}

	if ( oddRowCount() ) {
		pos.y += (-1.0 * tileSize.y);
	} else {
		pos.y += (-0.5 * tileSize.y);
	}

    gl_Position = projection * view * model * pos;
}
