package gfx

type ShaderType string

const (
	VertType ShaderType = "vert"
	FragType ShaderType = "frag"
)

//ring buffer
type RB struct {
	buf      []float32
	idx, cnt uint
}

func NewRB(cnt uint) *RB       { return &RB{buf: make([]float32, cnt), idx: 0, cnt: cnt} }
func (rb *RB) Add(val float32) { rb.buf[rb.idx] = val; rb.idx += 1; rb.idx %= rb.cnt }
func (rb *RB) Last() float32   { return rb.buf[(rb.idx+rb.cnt-1)%rb.cnt] }
func (rb *RB) Average() float32 {
	ret := float32(0.)
	for i := rb.idx; i < rb.idx+rb.cnt; i++ {
		ret += rb.buf[i%rb.cnt]
	}
	return ret / float32(rb.cnt)
}
