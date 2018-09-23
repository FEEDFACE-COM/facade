uniform sampler2D texture;

varying vec2 fragcoord;

void main() {
    vec4 tex = texture2D(texture,fragcoord);
    gl_FragColor = tex;
}
