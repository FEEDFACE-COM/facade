uniform sampler2D texture;

varying vec2 fragcoord;

uniform vec2 debugFlag;

void main() {
    vec4 tex = texture2D(texture,fragcoord);
    gl_FragColor = tex;
}
