// +build linux,arm
package facade

var ShaderAsset = map[string]string{

	"def.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCmF0dHJpYnV0ZSB2ZWMzIHZlcnRleDsK
YXR0cmlidXRlIHZlYzIgdGV4Q29vcmQ7Cgp2YXJ5aW5nIHZlYzIgdkZyYWdDb29yZDsKdmFyeWlu
ZyBmbG9hdCB2RGVidWdGbGFnOwoKYm9vbCBERUJVRyA9IGRlYnVnRmxhZyA+IDAuMDsKCnZvaWQg
bWFpbigpIHsKICAgIHZGcmFnQ29vcmQgPSB0ZXhDb29yZDsKICAgIGdsX1Bvc2l0aW9uID0gcHJv
amVjdGlvbiAqIHZpZXcgKiBtb2RlbCAqIHZlYzQodmVydGV4LCAxKTsKfQo=
`,

	"def.frag": `
dW5pZm9ybSBzYW1wbGVyMkQgdGV4dHVyZTsKCnZhcnlpbmcgdmVjMiB2RnJhZ0Nvb3JkOwp2YXJ5
aW5nIGZsb2F0IHZEZWJ1Z0ZsYWc7Cgpib29sIERFQlVHID0gdkRlYnVnRmxhZyA+IDAuMDsKCgp2
b2lkIG1haW4oKSB7CiAgICB2ZWM0IHRleCA9IHRleHR1cmUyRCh0ZXh0dXJlLHZGcmFnQ29vcmQp
OwogICAgZ2xfRnJhZ0NvbG9yID0gdGV4Owp9Cg==
`,

	"color.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCmF0dHJpYnV0ZSB2ZWMzIHZlcnRleDsK
YXR0cmlidXRlIHZlYzQgY29sb3I7Cgp2YXJ5aW5nIHZlYzQgIHZGcmFnQ29sb3I7CnZhcnlpbmcg
ZmxvYXQgdkRlYnVnRmxhZzsKCgp2b2lkIG1haW4oKSB7CiAgICB2RnJhZ0NvbG9yID0gY29sb3I7
CiAgICB2RGVidWdGbGFnID0gZGVidWdGbGFnOwogICAgZ2xfUG9zaXRpb24gPSBwcm9qZWN0aW9u
ICogdmlldyAqIG1vZGVsICogdmVjNCh2ZXJ0ZXgsIDEpOwp9Cg==
`,

	"color.frag": `
Cgp2YXJ5aW5nIGZsb2F0IHZEZWJ1Z0ZsYWc7CnZhcnlpbmcgdmVjNCB2RnJhZ0NvbG9yOwoKYm9v
bCBERUJVRyA9IHZEZWJ1Z0ZsYWcgPiAwLjA7Cgp2b2lkIG1haW4oKSB7CiAgICBnbF9GcmFnQ29s
b3IgPSB2RnJhZ0NvbG9yOwp9Cg==
`,

	"grid/def.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIHZlYzIgdGlsZVNpemU7CnVuaWZvcm0gdmVjMiB0aWxlQ291bnQ7CnVu
aWZvcm0gdmVjMiB0aWxlT2Zmc2V0OwoKdW5pZm9ybSBmbG9hdCBub3c7CnVuaWZvcm0gZmxvYXQg
c2Nyb2xsZXI7CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwoKYXR0cmlidXRlIHZlYzMgdmVydGV4
OwoKYXR0cmlidXRlIHZlYzIgdGV4Q29vcmQ7CmF0dHJpYnV0ZSB2ZWMyIHRpbGVDb29yZDsKYXR0
cmlidXRlIHZlYzIgZ3JpZENvb3JkOwoKCnZhcnlpbmcgdmVjMiB2VGV4Q29vcmQ7CnZhcnlpbmcg
dmVjMiB2VGlsZUNvb3JkOwp2YXJ5aW5nIHZlYzIgdkdyaWRDb29yZDsKdmFyeWluZyBmbG9hdCB2
U2Nyb2xsZXI7Cgpib29sIERFQlVHID0gZGVidWdGbGFnID4gMC4wOwoKCnZvaWQgbWFpbigpIHsK
ICAgIHZUZXhDb29yZCA9IHRleENvb3JkOwogICAgdlRpbGVDb29yZCA9IHRpbGVDb29yZDsKICAg
IHZHcmlkQ29vcmQgPSBncmlkQ29vcmQ7CiAgICB2U2Nyb2xsZXIgPSBhYnMoc2Nyb2xsZXIpOwog
ICAgCiAgICB2ZWM0IHBvcyA9IHZlYzQodmVydGV4LDEpOwoKICAgIHBvcy55ICs9IHNjcm9sbGVy
OwogICAgcG9zLnggKz0gKHRpbGVDb29yZC54ICogdGlsZVNpemUueCk7CiAgICBwb3MueSArPSAo
dGlsZUNvb3JkLnkgKiB0aWxlU2l6ZS55KTsKCiAgICBwb3MueCArPSAoIHRpbGVPZmZzZXQueCAq
IHRpbGVTaXplLngpOwogICAgcG9zLnkgKz0gKCB0aWxlT2Zmc2V0LnkgKiB0aWxlU2l6ZS55KTsK
CiAgICBnbF9Qb3NpdGlvbiA9IHByb2plY3Rpb24gKiB2aWV3ICogbW9kZWwgKiBwb3M7Cn0KCg==
`,

	"grid/def.frag": `
CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwp1bmlmb3JtIGZsb2F0IGRvd253YXJkOwp1bmlmb3Jt
IHZlYzIgdGlsZUNvdW50Owp1bmlmb3JtIHNhbXBsZXIyRCB0ZXh0dXJlOwp1bmlmb3JtIGZsb2F0
IHNjcm9sbGVyOwp1bmlmb3JtIHZlYzIgY3Vyc29yUG9zOwoKdmFyeWluZyB2ZWMyIHZUZXhDb29y
ZDsKdmFyeWluZyB2ZWMyIHZUaWxlQ29vcmQ7CnZhcnlpbmcgdmVjMiB2R3JpZENvb3JkOwp2YXJ5
aW5nIGZsb2F0IHZTY3JvbGxlcjsKCmJvb2wgREVCVUcgICAgPSBkZWJ1Z0ZsYWcgPiAwLjA7CgoK
CnZvaWQgbWFpbigpIHsKICAgIHZlYzQgY29sOwogICAgY29sID0gdGV4dHVyZTJEKHRleHR1cmUs
IHZUZXhDb29yZCk7IAoKICAgIGlmIChERUJVRykgeyAKICAgICAgICBjb2wucmdiID0gdmVjMygx
LiwxLiwxLik7CiAgICAgICAgY29sLmEgPSAxLjA7CiAgICB9IAoKICAgIGlmICggdkdyaWRDb29y
ZC55ID09IDAuMCApIHsgLy8gb2xkZXN0IGxpbmUKICAgICAgICBjb2wuYSAqPSAoMS4tdlNjcm9s
bGVyKTsKICAgIH0KICAgIAogICAgaWYgKCB2R3JpZENvb3JkLnkgPT0gdGlsZUNvdW50LnkgKSB7
IC8vIG5ld2VzdCBsaW5lCiAgICAgICAgY29sLmEgKj0gdlNjcm9sbGVyOyAgICAKICAgIH0KICAg
IAogICAgaWYgKCBjdXJzb3JQb3MueCA9PSB2R3JpZENvb3JkLnggJiYgY3Vyc29yUG9zLnkgPT0g
dkdyaWRDb29yZC55ICkgeyAvLyBpbnZlcnQgY3Vyc29yCiAgICAgICAgY29sLnJnYmEgPSAxLiAt
IGNvbC5yZ2JhOwogICAgfQoKICAgIGlmICghZ2xfRnJvbnRGYWNpbmcpIHsgY29sLmEgLz0gNC47
IH0KCiAgICBnbF9GcmFnQ29sb3IgPSBjb2w7CiAgICAKfQo=
`,

	"grid/debug.frag": `
CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwp1bmlmb3JtIGZsb2F0IGRvd253YXJkOwp1bmlmb3Jt
IHZlYzIgdGlsZUNvdW50Owp1bmlmb3JtIHNhbXBsZXIyRCB0ZXh0dXJlOwp1bmlmb3JtIGZsb2F0
IHNjcm9sbGVyOwp1bmlmb3JtIHZlYzIgY3Vyc29yUG9zOwoKdmFyeWluZyB2ZWMyIHZUZXhDb29y
ZDsKdmFyeWluZyB2ZWMyIHZUaWxlQ29vcmQ7CnZhcnlpbmcgdmVjMiB2R3JpZENvb3JkOwp2YXJ5
aW5nIGZsb2F0IHZTY3JvbGxlcjsKCmJvb2wgREVCVUcgICAgPSBkZWJ1Z0ZsYWcgPiAwLjA7CgoK
CnZvaWQgbWFpbigpIHsKICAgIHZlYzQgY29sOwogICAgY29sID0gdGV4dHVyZTJEKHRleHR1cmUs
IHZUZXhDb29yZCk7IAoKICAgIGlmICh0cnVlKSB7IAogICAgICAgIGNvbC5yZ2IgPSB2ZWMzKDEu
LDEuLDEuKTsKICAgICAgICBjb2wuYSA9IDEuMDsKICAgIH0gCgogICAgZmxvYXQgRiA9IDEuOwoK
ICAgIGZsb2F0IHggPSB2R3JpZENvb3JkLnggLyB0aWxlQ291bnQueDsKICAgIGZsb2F0IHkgPSB2
R3JpZENvb3JkLnkgLyB0aWxlQ291bnQueTsKICAgIAogICAgY29sLnIgKj0gRiAqICgxLiAtIHgp
OwogICAgY29sLmcgKj0gRiAqICgxLiAtIHkpOwoKICAgIGlmICggYWJzKHZHcmlkQ29vcmQueSkg
PT0gdGlsZUNvdW50LnkgICkgewogICAgICAgIGNvbC5yID0gMS4wOwogICAgICAgIGNvbC5nID0g
MS4wOwogICAgICAgIGNvbC5iID0gMC47CiAgICB9CiAgICAKICAgIGlmICggY3Vyc29yUG9zLngg
PT0gdkdyaWRDb29yZC54ICYmIGN1cnNvclBvcy55ID09IHZHcmlkQ29vcmQueSApIHsKICAgICAg
ICBjb2wucmdiYSA9IDEuIC0gY29sLnJnYmE7CiAgICB9CgogICAgaWYgKGdsX0Zyb250RmFjaW5n
KSB7IAogICAgICAgIHZlYzMgdG1wID0gdmVjMyhjb2wucmdiKTsKICAgICAgICBjb2wuciA9IHRt
cC5nOwogICAgICAgIGNvbC5nID0gdG1wLmI7CiAgICAgICAgY29sLmIgPSB0bXAucjsKICAgIH0K
CiAgICBnbF9GcmFnQ29sb3IgPSBjb2w7CiAgICAKfQo=
`,

	"grid/debug2.frag": `
CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwp1bmlmb3JtIGZsb2F0IGRvd253YXJkOwp1bmlmb3Jt
IHZlYzIgdGlsZUNvdW50Owp1bmlmb3JtIHNhbXBsZXIyRCB0ZXh0dXJlOwp1bmlmb3JtIGZsb2F0
IHNjcm9sbGVyOwp1bmlmb3JtIHZlYzIgY3Vyc29yUG9zOwoKCgp2YXJ5aW5nIHZlYzIgdlRleENv
b3JkOwp2YXJ5aW5nIHZlYzIgdlRpbGVDb29yZDsKdmFyeWluZyB2ZWMyIHZHcmlkQ29vcmQ7CnZh
cnlpbmcgZmxvYXQgdlNjcm9sbGVyOwoKYm9vbCBERUJVRyAgICA9IGRlYnVnRmxhZyA+IDAuMDsK
CgoKdm9pZCBtYWluKCkgewogICAgdmVjNCBjb2w7CiAgICBjb2wgPSB0ZXh0dXJlMkQodGV4dHVy
ZSwgdlRleENvb3JkKTsgCgogICAgaWYgKERFQlVHKSB7IAogICAgICAgIGNvbC5yZ2IgPSB2ZWMz
KDEuLDEuLDEuKTsKICAgICAgICBjb2wuYSA9IDEuMDsKICAgIH0gCgogICAgY29sLnJnYiA9IHZl
YzMoMC4sMC4sMC4pOwogICAgY29sLmEgPSAxLjsKICAgIAogICAgaWYgKCB2VGlsZUNvb3JkLngg
PT0gMC4wIHx8IHZUaWxlQ29vcmQueSA9PSAwLjAgKSB7CiAgICAgICAgCiAgICAgICAgY29sLnJn
YiArPSB2ZWMzKDEuLDEuLDEuKTsKICAgIAogICAgfSAKICAgIGlmICggdlRpbGVDb29yZC54ICsg
dGlsZUNvdW50LngvMi4gPD0gLjUgKSB7CiAgICAKICAgICAgICBjb2wucmdiICs9IHZlYzMoMC4s
MS4sMC4pOwogICAgCiAgICB9CiAgICBpZiAoIHZUaWxlQ29vcmQueCArIHRpbGVDb3VudC54LzIu
ID49IHRpbGVDb3VudC54LTIuICkgewogICAgCiAgICAgICAgY29sLnJnYiArPSB2ZWMzKDEuLDAu
LDEuKTsKICAgIAogICAgfQogICAgaWYgKCB2VGlsZUNvb3JkLnkgKyB0aWxlQ291bnQueS8yLiA8
PSAuNSApIHsKICAgIAogICAgICAgIGNvbC5yZ2IgKz0gdmVjMygxLiwwLiwwLik7CgogICAgfQog
ICAgaWYgKCB2VGlsZUNvb3JkLnkgKyB0aWxlQ291bnQueS8yLiA+PSB0aWxlQ291bnQueS0xLiAp
IHsKICAgIAogICAgICAgIGNvbC5yZ2IgKz0gdmVjMygwLiwxLiwxLik7CiAgICAKICAgIH0KICAg
IGlmICggbW9kKC1hYnModlRpbGVDb29yZC54KSAsIDIuKSA9PSAwLjAgIF5eIG1vZCgtYWJzKHZU
aWxlQ29vcmQueSkgLCAyLikgPT0gMC4wICkgewogICAgCiAgICAgICAgY29sLnJnYiArPSAwLjUg
KiB2ZWMzKDEuLDEuLDEuKTsKICAgIAogICAgfQogICAgCiAgICAKICAgIGlmICggY3Vyc29yUG9z
LnggPT0gdkdyaWRDb29yZC54ICYmIGN1cnNvclBvcy55ID09IHZHcmlkQ29vcmQueSApIHsKICAg
ICAgICBjb2wucmdiYSA9IDEuIC0gY29sLnJnYmE7CiAgICB9CgogICAgaWYgKGdsX0Zyb250RmFj
aW5nKSB7IAogICAgICAgIGNvbC5yZ2IgLz0gMi47CiAgICB9CgogICAgZ2xfRnJhZ0NvbG9yID0g
Y29sOwogICAgCn0K
`,

	"grid/wave.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIHZlYzIgdGlsZVNpemU7CnVuaWZvcm0gdmVjMiB0aWxlQ291bnQ7CnVu
aWZvcm0gdmVjMiB0aWxlT2Zmc2V0OwoKdW5pZm9ybSBmbG9hdCBub3c7CnVuaWZvcm0gZmxvYXQg
c2Nyb2xsZXI7CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwp1bmlmb3JtIGZsb2F0IGRvd253YXJk
OwoKYXR0cmlidXRlIHZlYzMgdmVydGV4OwphdHRyaWJ1dGUgdmVjMiB0ZXhDb29yZDsKYXR0cmli
dXRlIHZlYzIgdGlsZUNvb3JkOwphdHRyaWJ1dGUgdmVjMiBncmlkQ29vcmQ7CgoKdmFyeWluZyB2
ZWMyIHZUZXhDb29yZDsKdmFyeWluZyB2ZWMyIHZUaWxlQ29vcmQ7CnZhcnlpbmcgdmVjMiB2R3Jp
ZENvb3JkOwoKdmFyeWluZyBmbG9hdCB2U2Nyb2xsZXI7Cgpib29sIERFQlVHID0gZGVidWdGbGFn
ID4gMC4wOwoKZmxvYXQgUEkgPSAzLjE0MTU5MjY1MzU4OTc5MzIzODQ2MjY0MzM4MzI3OTUwMjg4
NDE5NzE2OTM5OTM3NTEwNTgyMDk3NDk0NDU5MjA7CgoKdm9pZCBtYWluKCkgewogICAgdlRleENv
b3JkID0gdGV4Q29vcmQ7CiAgICB2VGlsZUNvb3JkID0gdGlsZUNvb3JkOwogICAgdlNjcm9sbGVy
ID0gYWJzKHNjcm9sbGVyKTsKICAgIHZHcmlkQ29vcmQgPSBncmlkQ29vcmQ7CiAgICAKICAgIHZl
YzQgcG9zID0gdmVjNCh2ZXJ0ZXgsMSk7CgogICAgcG9zLnkgKz0gc2Nyb2xsZXI7CiAgICBwb3Mu
eCArPSAodGlsZUNvb3JkLnggKiB0aWxlU2l6ZS54KTsKICAgIHBvcy55ICs9ICh0aWxlQ29vcmQu
eSAqIHRpbGVTaXplLnkpOwoKICAgIHBvcy54ICs9ICggdGlsZU9mZnNldC54ICogdGlsZVNpemUu
eCk7CiAgICBwb3MueSArPSAoIHRpbGVPZmZzZXQueSAqIHRpbGVTaXplLnkpOwoKICAgIGZsb2F0
IEYgPSAwLjU7CiAgICBmbG9hdCB4ID0gcG9zLng7CiAgICBmbG9hdCB5ID0gcG9zLnk7CiAgICAK
ICAgIGZsb2F0IGZyZXEgPSAtMS4vMjQuOwogICAgcG9zLnkgKz0gRiAqIGNvcyggMi4gKiBmcmVx
ICogeCAqIFBJICsgbm93ICAgICAgICAgKTsKICAgIHBvcy54ICs9IEYgKiBjb3MoIDMuICogZnJl
cSAqIHkgKiBQSSArIG5vdyArIFBJLzIuICk7CiAgICAKICAgIHBvcy56ICs9IEYgKiBjb3MoIDUu
ICogZnJlcSAqICh4K3kpICogUEkgKyBub3cgKyBQSS8yLiApOwoKICAgIGdsX1Bvc2l0aW9uID0g
cHJvamVjdGlvbiAqIHZpZXcgKiBtb2RlbCAqIHBvczsKfQoK
`,

	"grid/roll.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIHZlYzIgdGlsZVNpemU7CnVuaWZvcm0gdmVjMiB0aWxlQ291bnQ7CnVu
aWZvcm0gdmVjMiB0aWxlT2Zmc2V0OwoKdW5pZm9ybSBmbG9hdCBub3c7CnVuaWZvcm0gZmxvYXQg
c2Nyb2xsZXI7CnVuaWZvcm0gZmxvYXQgc2NyZWVuUmF0aW87CnVuaWZvcm0gZmxvYXQgZm9udFJh
dGlvOwp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCmF0dHJpYnV0ZSB2ZWMzIHZlcnRleDsKCmF0
dHJpYnV0ZSB2ZWMyIHRleENvb3JkOwphdHRyaWJ1dGUgdmVjMiB0aWxlQ29vcmQ7CmF0dHJpYnV0
ZSB2ZWMyIGdyaWRDb29yZDsKCgp2YXJ5aW5nIHZlYzIgdlRleENvb3JkOwp2YXJ5aW5nIHZlYzIg
dlRpbGVDb29yZDsKdmFyeWluZyB2ZWMyIHZHcmlkQ29vcmQ7CnZhcnlpbmcgZmxvYXQgdlNjcm9s
bGVyOwoKYm9vbCBERUJVRyA9IGRlYnVnRmxhZyA+IDAuMDsKCgpmbG9hdCBQSSA9IDMuMTQxNTky
NjUzNTg5NzkzMjM4NDYyNjQzMzgzMjc5NTAyODg0MTk3MTY5Mzk5Mzc1MTA1ODIwOTc0OTQ0NTky
MDsKZmxvYXQgVEFVPSA2LjI4MzE4NTMwNzE3OTU4NjQ3NjkyNTI4Njc2NjU1OTAwNTc2ODM5NDMz
ODc5ODc1MDIxMTY0MTk0OTg4OTE4NDA7CmZsb2F0IGVhc2UxKGZsb2F0IHgpICAgICAgICAgIHsg
cmV0dXJuIDAuNSAqIGNvcyggICAgIHggKyBQSS8yLjAgKSArIDAuNTsgfQoKCgoKdm9pZCBtYWlu
KCkgewogICAgdlRleENvb3JkID0gdGV4Q29vcmQ7CiAgICB2VGlsZUNvb3JkID0gdGlsZUNvb3Jk
OwogICAgdkdyaWRDb29yZCA9IGdyaWRDb29yZDsKICAgIHZTY3JvbGxlciA9IGFicyhzY3JvbGxl
cik7CiAgICAKICAgIHZlYzQgcG9zID0gdmVjNCh2ZXJ0ZXgsMSk7CgogICAgZmxvYXQgb2Zmc2V0
ID0gUEkvMTYuOwogICAgb2Zmc2V0ID0gMC4wOwogICAgZmxvYXQgQVJDID0gUEkvMi4gLSBvZmZz
ZXQ7CiAgICBmbG9hdCBSQURJVVMgPSB0aWxlQ291bnQueS8yLiAqIHRpbGVTaXplLnkgLzIuIDsK
CgogICAgCgogICAgZmxvYXQgZGVsdGEgPSAwLjA7CiAgICBmbG9hdCBhbHBoYSxiZXRhOwoKICAg
IGZsb2F0IHggPSB0aWxlQ29vcmQueC0wLjUrdGlsZU9mZnNldC54OwogICAgZmxvYXQgeSA9IHRp
bGVDb29yZC55LTEuNSt0aWxlT2Zmc2V0Lnk7CiAgICAKCiAgICBhbHBoYSA9IC0xLiAqIEFSQyAv
ICh0aWxlQ291bnQueSk7CiAgICBkZWx0YSA9IFBJLzIuIC0gb2Zmc2V0ICsgYWxwaGE7CiAgICBi
ZXRhID0gZGVsdGEgKyAoIGFscGhhICogKHNjcm9sbGVyK3kpICkgOwoKCiAgICBmbG9hdCByID0g
UkFESVVTICogMi47CiAgICAKICAgIAogICAgdmVjMyBBID0gdmVjMyggKHgrMS4pKnRpbGVTaXpl
LngsIGNvcyhhbHBoYStiZXRhKSpyLCAtdGlsZUNvdW50LnkvMi4rc2luKGFscGhhK2JldGEpKnIp
OwogICAgdmVjMyBCID0gdmVjMyggKHgrMS4pKnRpbGVTaXplLngsIGNvcyhiZXRhKSpyLCAgICAg
ICAtdGlsZUNvdW50LnkvMi4rc2luKGJldGEpKnIpOwogICAgdmVjMyBDID0gdmVjMyggKHgpKnRp
bGVTaXplLngsICAgIGNvcyhhbHBoYStiZXRhKSpyLCAtdGlsZUNvdW50LnkvMi4rc2luKGFscGhh
K2JldGEpKnIpOwogICAgdmVjMyBEID0gdmVjMyggKHgpKnRpbGVTaXplLngsICAgIGNvcyhiZXRh
KSpyLCAgICAgICAtdGlsZUNvdW50LnkvMi4rc2luKGJldGEpKnIpOwogICAgCiAgIAogICAgaWYg
KCBwb3MueCA+IDAuICYmIHBvcy55ID4gMC4gKSB7CiAgICAgICAgcG9zLnh5eiA9IEE7CiAgICB9
IGVsc2UgaWYgKCBwb3MueCA+IDAuICYmIHBvcy55IDwgMC4gKSB7CiAgICAgICAgcG9zLnh5eiA9
IEI7CiAgICB9IGVsc2UgaWYgKCBwb3MueCA8IDAuICYmIHBvcy55ID4gMC4gKSB7CiAgICAgICAg
cG9zLnh5eiA9IEM7CiAgICB9IGVsc2UgaWYgKCBwb3MueCA8IDAuICYmIHBvcy55IDwgMC4gKSB7
CiAgICAgICAgcG9zLnh5eiA9IEQ7CiAgICB9CgoKICAgIGZsb2F0IHJhdGlvID0gc2NyZWVuUmF0
aW8gLyBmb250UmF0aW87CiAgICBmbG9hdCB6b29tID0gcmF0aW8gKiAyLiAvICggdGlsZUNvdW50
LnggKTsKICAgIAogICAgbWF0NCBtZGwgPSBtYXQ0KDEuMCk7CiAgICBtZGxbMF1bMF0gPSB6b29t
OwogICAgbWRsWzFdWzFdID0gem9vbTsKICAgIG1kbFsyXVsyXSA9IHpvb207CgogICAgZ2xfUG9z
aXRpb24gPSBwcm9qZWN0aW9uICogdmlldyAqIG1kbCAqIHBvczsKfQoK
`,

	"grid/rows.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIHZlYzIgdGlsZVNpemU7CnVuaWZvcm0gdmVjMiB0aWxlQ291bnQ7CnVu
aWZvcm0gdmVjMiB0aWxlT2Zmc2V0OwoKdW5pZm9ybSBmbG9hdCBub3c7CnVuaWZvcm0gZmxvYXQg
c2Nyb2xsZXI7CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwoKYXR0cmlidXRlIHZlYzMgdmVydGV4
OwoKYXR0cmlidXRlIHZlYzIgdGV4Q29vcmQ7CmF0dHJpYnV0ZSB2ZWMyIHRpbGVDb29yZDsKYXR0
cmlidXRlIHZlYzIgZ3JpZENvb3JkOwoKCnZhcnlpbmcgdmVjMiB2VGV4Q29vcmQ7CnZhcnlpbmcg
dmVjMiB2VGlsZUNvb3JkOwp2YXJ5aW5nIHZlYzIgdkdyaWRDb29yZDsKdmFyeWluZyBmbG9hdCB2
U2Nyb2xsZXI7Cgpib29sIERFQlVHID0gZGVidWdGbGFnID4gMC4wOwoKCmZsb2F0IFBJID0gMy4x
NDE1OTI2NTM1ODk3OTMyMzg0NjI2NDMzODMyNzk1MDI4ODQxOTcxNjkzOTkzNzUxMDU4MjA5NzQ5
NDQ1OTIwOwoKCnZvaWQgbWFpbigpIHsKICAgIHZUZXhDb29yZCA9IHRleENvb3JkOwogICAgdlRp
bGVDb29yZCA9IHRpbGVDb29yZDsKICAgIHZHcmlkQ29vcmQgPSBncmlkQ29vcmQ7CiAgICB2U2Ny
b2xsZXIgPSBhYnMoc2Nyb2xsZXIpOwogICAgCiAgICB2ZWM0IHBvcyA9IHZlYzQodmVydGV4LDEp
OwoKICAgIHBvcy55ICs9IHNjcm9sbGVyOwogICAgcG9zLnggKz0gKHRpbGVDb29yZC54ICogdGls
ZVNpemUueCk7CiAgICBwb3MueSArPSAodGlsZUNvb3JkLnkgKiB0aWxlU2l6ZS55KTsKCiAgICBw
b3MueCArPSAoIHRpbGVPZmZzZXQueCAqIHRpbGVTaXplLngpOwogICAgcG9zLnkgKz0gKCB0aWxl
T2Zmc2V0LnkgKiB0aWxlU2l6ZS55KTsKCiAgICBmbG9hdCBGID0gMS4gKyB0aWxlQ291bnQueS84
LjsKCiAgICBmbG9hdCB5ICA9ICB2VGlsZUNvb3JkLnkgICAgICAgLyAodGlsZUNvdW50LnkvMi4p
OwogICAgZmxvYXQgeXkgPSAodlRpbGVDb29yZC55ICsgKChzY3JvbGxlcikpICkgLyAodGlsZUNv
dW50LnkvMi4pOwoKCiAgICBmbG9hdCBmcmVxID0gLTEuOwogICAgZmxvYXQgZjAgPSBjb3MoIGZy
ZXEgKiB5ICAqIFBJICsgbm93ICsgUEkvMi4gKTsKICAgIGZsb2F0IGYxID0gY29zKCBmcmVxICog
eXkgKiBQSSArIG5vdyArIFBJLzIuICk7CiAgICBmbG9hdCBkID0gIGYwICsgLyooc2Nyb2xsZXIp
ICogKi8oZjEgLSBmMCk7CiAgICBwb3MueiArPSBGICogZDsKICAgIHBvcy56IC09IEY7CgoKICAg
IAogICAgZ2xfUG9zaXRpb24gPSBwcm9qZWN0aW9uICogdmlldyAqIG1vZGVsICogcG9zOwp9Cg==
`,

	"grid/crawl.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIHZlYzIgdGlsZVNpemU7CnVuaWZvcm0gdmVjMiB0aWxlQ291bnQ7CnVu
aWZvcm0gdmVjMiB0aWxlT2Zmc2V0OwoKdW5pZm9ybSBmbG9hdCBub3c7CnVuaWZvcm0gZmxvYXQg
c2Nyb2xsZXI7CnVuaWZvcm0gZmxvYXQgc2NyZWVuUmF0aW87CnVuaWZvcm0gZmxvYXQgZm9udFJh
dGlvOwp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCmF0dHJpYnV0ZSB2ZWMzIHZlcnRleDsKCmF0
dHJpYnV0ZSB2ZWMyIHRleENvb3JkOwphdHRyaWJ1dGUgdmVjMiB0aWxlQ29vcmQ7CmF0dHJpYnV0
ZSB2ZWMyIGdyaWRDb29yZDsKCgp2YXJ5aW5nIHZlYzIgdlRleENvb3JkOwp2YXJ5aW5nIHZlYzIg
dlRpbGVDb29yZDsKdmFyeWluZyB2ZWMyIHZHcmlkQ29vcmQ7CnZhcnlpbmcgZmxvYXQgdlNjcm9s
bGVyOwoKYm9vbCBERUJVRyA9IGRlYnVnRmxhZyA+IDAuMDsKCmZsb2F0IGR5KGZsb2F0IHkpIHsK
ICAgIHJldHVybiAxLiAtICggIDAuMjUgKyAwLjUgKiAoKHkrdGlsZUNvdW50LnkvMi4pIC8gdGls
ZUNvdW50LnkgKSApOwp9Cgp2b2lkIG1haW4oKSB7CiAgICB2VGV4Q29vcmQgPSB0ZXhDb29yZDsK
ICAgIHZUaWxlQ29vcmQgPSB0aWxlQ29vcmQ7CiAgICB2R3JpZENvb3JkID0gZ3JpZENvb3JkOwog
ICAgdlNjcm9sbGVyID0gYWJzKHNjcm9sbGVyKTsKICAgIAogICAgCiAgICBmbG9hdCB3MCA9IGR5
KHRpbGVDb29yZC55KzEuKSAqIHRpbGVTaXplLng7CiAgICBmbG9hdCB3MSA9IGR5KHRpbGVDb29y
ZC55KSAgICAqIHRpbGVTaXplLng7CiAgICBmbG9hdCBoICA9IHRpbGVTaXplLnk7CiAgICAKCiAg
ICB2ZWM0IHBvcyA9IHZlYzQoMC4sMC4sMC4sMS4pOwoKICAgIHZlYzIgQSA9IHZlYzIoKHcwKih0
aWxlQ29vcmQueC0xLikpICwgKGggKiAodGlsZUNvb3JkLnkgICApKSk7CiAgICB2ZWMyIEIgPSB2
ZWMyKCh3MSoodGlsZUNvb3JkLngtMS4pKSAsIChoICogKHRpbGVDb29yZC55LTEuKSkpOwogICAg
dmVjMiBDID0gdmVjMigodzEqKHRpbGVDb29yZC54ICAgKSkgLCAoaCAqICh0aWxlQ29vcmQueS0x
LikpKTsKICAgIHZlYzIgRCA9IHZlYzIoKHcwKih0aWxlQ29vcmQueCAgICkpICwgKGggKiAodGls
ZUNvb3JkLnkgICApKSk7CiAgICAKICAgIAogICAKICAgIGlmICAgICAgICAoIHZlcnRleC54IDwg
MC4gJiYgdmVydGV4LnkgPiAwLiApIHsKICAgICAgICBwb3MueHkgPSBBOwogICAgfSBlbHNlIGlm
ICggdmVydGV4LnggPCAwLiAmJiB2ZXJ0ZXgueSA8IDAuICkgewogICAgICAgIHBvcy54eSA9IEI7
CiAgICB9IGVsc2UgaWYgKCB2ZXJ0ZXgueCA+IDAuICYmIHZlcnRleC55IDwgMC4gKSB7CiAgICAg
ICAgcG9zLnh5ID0gQzsKICAgIH0gZWxzZSBpZiAoIHZlcnRleC54ID4gMC4gJiYgdmVydGV4Lnkg
PiAwLiApIHsKICAgICAgICBwb3MueHkgPSBEOwogICAgfQogICAgCiAgICBwb3MueHkgKz0gKHNj
cm9sbGVyICogKEEtQikpOwoKCiAgICBmbG9hdCByYXRpbyA9IHNjcmVlblJhdGlvIC8gZm9udFJh
dGlvOwogICAgZmxvYXQgem9vbV9jb2xzID0gcmF0aW8gKiAyLiAvIHRpbGVDb3VudC54ICogMS41
OwogICAgZmxvYXQgem9vbV9yb3dzID0gMi4vKHRpbGVDb3VudC55KzEuKTsKICAgIGZsb2F0IHpv
b20gPSBtaW4oem9vbV9jb2xzLHpvb21fcm93cyk7CiAgICAKCiAgICBtYXQ0IG1kbDsKICAgIG1k
bCA9IG1hdDQoMS4wKTsKICAgIG1kbFswXVswXSA9IHpvb207CiAgICBtZGxbMV1bMV0gPSB6b29t
OwoKICAgIGdsX1Bvc2l0aW9uID0gcHJvamVjdGlvbiAqIHZpZXcgKiBtZGwgKiBwb3M7Cn0KCg==
`,

	"grid/disk.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIHZlYzIgdGlsZVNpemU7CnVuaWZvcm0gdmVjMiB0aWxlQ291bnQ7CnVu
aWZvcm0gdmVjMiB0aWxlT2Zmc2V0OwoKdW5pZm9ybSBmbG9hdCBub3c7CnVuaWZvcm0gZmxvYXQg
c2Nyb2xsZXI7CnVuaWZvcm0gZmxvYXQgc2NyZWVuUmF0aW87CnVuaWZvcm0gZmxvYXQgZm9udFJh
dGlvOwp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCmF0dHJpYnV0ZSB2ZWMzIHZlcnRleDsKCmF0
dHJpYnV0ZSB2ZWMyIHRleENvb3JkOwphdHRyaWJ1dGUgdmVjMiB0aWxlQ29vcmQ7CmF0dHJpYnV0
ZSB2ZWMyIGdyaWRDb29yZDsKCgp2YXJ5aW5nIHZlYzIgdlRleENvb3JkOwp2YXJ5aW5nIHZlYzIg
dlRpbGVDb29yZDsKdmFyeWluZyB2ZWMyIHZHcmlkQ29vcmQ7CnZhcnlpbmcgZmxvYXQgdlNjcm9s
bGVyOwoKYm9vbCBERUJVRyA9IGRlYnVnRmxhZyA+IDAuMDsKCgpmbG9hdCBQSSA9IDMuMTQxNTky
NjUzNTg5NzkzMjM4NDYyNjQzMzgzMjc5NTAyODg0MTk3MTY5Mzk5Mzc1MTA1ODIwOTc0OTQ0NTky
MDsKZmxvYXQgVEFVPSA2LjI4MzE4NTMwNzE3OTU4NjQ3NjkyNTI4Njc2NjU1OTAwNTc2ODM5NDMz
ODc5ODc1MDIxMTY0MTk0OTg4OTE4NDA7CmZsb2F0IGVhc2UxKGZsb2F0IHgpICAgICAgICAgIHsg
cmV0dXJuIDAuNSAqIGNvcyggICAgIHggKyBQSS8yLjAgKSArIDAuNTsgfQoKCgoKdm9pZCBtYWlu
KCkgewogICAgdlRleENvb3JkID0gdGV4Q29vcmQ7CiAgICB2VGlsZUNvb3JkID0gdGlsZUNvb3Jk
OwogICAgdkdyaWRDb29yZCA9IGdyaWRDb29yZDsKICAgIHZTY3JvbGxlciA9IGFicyhzY3JvbGxl
cik7CiAgICAKICAgIHZlYzQgcG9zID0gdmVjNCh2ZXJ0ZXgsMSk7CgoKICAgIGZsb2F0IFJBRElV
UyA9IHRpbGVDb3VudC55LzIuIDsKICAgIGZsb2F0IFIwID0gUkFESVVTLzQuOwogICAgZmxvYXQg
cmFkID0gUkFESVVTIC8gKHRpbGVDb3VudC55ICsgUjApOyAKCgogICAgZmxvYXQgZGVsdGEgPSAw
LjA7Ci8vICAgIGRlbHRhICs9IG5vdy8xMC47CiAgICBkZWx0YSArPSBlYXNlMShub3cvMi4pIC0g
MC41OwogICAgCgogICAgZmxvYXQgQVJDID0gVEFVOwogICAgZmxvYXQgQTAgPSAwLjA7CiAgCiAg
ICBmbG9hdCBhbHBoYSxnYW1tYTsKICAgIAogICAgZmxvYXQgcm93ID0gKC10aWxlQ29vcmQueSt0
aWxlQ291bnQueS8yLik7CgoKICAgIGFscGhhID0gQVJDIC8gKEEwICsgdGlsZUNvdW50LngpOwog
ICAgZ2FtbWEgKz0gZGVsdGE7CiAgICBnYW1tYSArPSAoIEFSQyAvICh0aWxlQ291bnQueCtBMCkp
ICogdGlsZUNvb3JkLng7CgoKICAgIAogICAgZmxvYXQgcjAgPSBSMCArIChyYWQgKiByb3cgKSA7
CiAgICBmbG9hdCByMSA9IHIwICsgcmFkOwoKICAgIHIwIC09IChzY3JvbGxlcipyYWQpOwogICAg
cjEgLT0gKHNjcm9sbGVyKnJhZCk7CgogICAgCiAgICB2ZWMyIEEgPSB2ZWMyKCBjb3MoZ2FtbWEr
YWxwaGEpKnIwLCBzaW4oZ2FtbWErYWxwaGEpKnIwKTsKICAgIHZlYzIgQiA9IHZlYzIoIGNvcyhn
YW1tYSthbHBoYSkqcjEsIHNpbihnYW1tYSthbHBoYSkqcjEpOwogICAgdmVjMiBDID0gdmVjMigg
Y29zKGdhbW1hICAgICAgKSpyMSwgc2luKGdhbW1hICAgICAgKSpyMSk7CiAgICB2ZWMyIEQgPSB2
ZWMyKCBjb3MoZ2FtbWEgICAgICApKnIwLCBzaW4oZ2FtbWEgICAgICApKnIwKTsKICAgIAogICAK
ICAgCiAgICBpZiAgICAgICAgKCBwb3MueCA+IDAuICYmIHBvcy55ID4gMC4gKSB7CiAgICAgICAg
cG9zLnh5ID0gQTsKICAgIH0gZWxzZSBpZiAoIHBvcy54ID4gMC4gJiYgcG9zLnkgPCAwLiApIHsK
ICAgICAgICBwb3MueHkgPSBCOwogICAgfSBlbHNlIGlmICggcG9zLnggPCAwLiAmJiBwb3MueSA+
IDAuICkgewogICAgICAgIHBvcy54eSA9IEQ7CiAgICB9IGVsc2UgaWYgKCBwb3MueCA8IDAuICYm
IHBvcy55IDwgMC4gKSB7CiAgICAgICAgcG9zLnh5ID0gQzsKICAgIH0KCiAgICBmbG9hdCB6b29t
ID0gIDcuLyAoNC4qdGlsZUNvdW50LnkpOwogICAgbWF0NCBtZGwgPSBtYXQ0KDEuMCk7CiAgICBt
ZGxbMF1bMF0gPSB6b29tOwogICAgbWRsWzFdWzFdID0gem9vbTsKICAgIAoKICAgIGdsX1Bvc2l0
aW9uID0gcHJvamVjdGlvbiAqIHZpZXcgKiBtZGwgKiBwb3M7Cgp9Cgo=
`,

	"grid/drop.vert": `
dW5pZm9ybSBtYXQ0IHByb2plY3Rpb247CnVuaWZvcm0gbWF0NCB2aWV3Owp1bmlmb3JtIG1hdDQg
bW9kZWw7Cgp1bmlmb3JtIHZlYzIgdGlsZVNpemU7CnVuaWZvcm0gdmVjMiB0aWxlQ291bnQ7CnVu
aWZvcm0gdmVjMiB0aWxlT2Zmc2V0OwoKdW5pZm9ybSBmbG9hdCBub3c7CnVuaWZvcm0gZmxvYXQg
c2Nyb2xsZXI7CnVuaWZvcm0gZmxvYXQgc2NyZWVuUmF0aW87CnVuaWZvcm0gZmxvYXQgZm9udFJh
dGlvOwp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCmF0dHJpYnV0ZSB2ZWMzIHZlcnRleDsKCmF0
dHJpYnV0ZSB2ZWMyIHRleENvb3JkOwphdHRyaWJ1dGUgdmVjMiB0aWxlQ29vcmQ7CmF0dHJpYnV0
ZSB2ZWMyIGdyaWRDb29yZDsKCgp2YXJ5aW5nIHZlYzIgdlRleENvb3JkOwp2YXJ5aW5nIHZlYzIg
dlRpbGVDb29yZDsKdmFyeWluZyB2ZWMyIHZHcmlkQ29vcmQ7CnZhcnlpbmcgZmxvYXQgdlNjcm9s
bGVyOwoKYm9vbCBERUJVRyA9IGRlYnVnRmxhZyA+IDAuMDsKCgpmbG9hdCBQSSA9IDMuMTQxNTky
NjUzNTg5NzkzMjM4NDYyNjQzMzgzMjc5NTAyODg0MTk3MTY5Mzk5Mzc1MTA1ODIwOTc0OTQ0NTky
MDsKZmxvYXQgVEFVPSA2LjI4MzE4NTMwNzE3OTU4NjQ3NjkyNTI4Njc2NjU1OTAwNTc2ODM5NDMz
ODc5ODc1MDIxMTY0MTk0OTg4OTE4NDA7CmZsb2F0IGVhc2UxKGZsb2F0IHgpICAgICAgICAgIHsg
cmV0dXJuIDAuNSAqIGNvcyggICAgIHggKyBQSS8yLjAgKSArIDAuNTsgfQoKCgp2b2lkIG1haW4o
KSB7CiAgICB2VGV4Q29vcmQgPSB0ZXhDb29yZDsKICAgIHZUaWxlQ29vcmQgPSB0aWxlQ29vcmQ7
CiAgICB2R3JpZENvb3JkID0gZ3JpZENvb3JkOwogICAgdlNjcm9sbGVyID0gYWJzKHNjcm9sbGVy
KTsKICAgIAogICAgdmVjNCBwb3MgPSB2ZWM0KHZlcnRleCwxKTsKCiAgICBmbG9hdCBSQURJVVMg
PSA0Li8yLiAqIHRpbGVDb3VudC55OwogICAgZmxvYXQgUjAgPSBSQURJVVMvNC47CiAgICBmbG9h
dCByYWQgPSBSQURJVVMgLyAodGlsZUNvdW50LnkgKyBSMCk7IAoKCiAgICBmbG9hdCBkZWx0YSA9
IDMuKlBJLzIuOwogICAgZGVsdGEgKz0gZWFzZTEobm93LzIuKSAtIDAuNTsKCiAgICBmbG9hdCBB
UkMgPSBUQVU7CiAgICBmbG9hdCBBMCA9IDAuOwogIAogICAgZmxvYXQgYWxwaGEsZ2FtbWE7CiAg
ICAKICAgIGZsb2F0IHJvdyA9ICgtdGlsZUNvb3JkLnkrdGlsZUNvdW50LnkvMi4pOwoKCiAgICBh
bHBoYSA9IEFSQyAvIChBMCArIHRpbGVDb3VudC54KTsKICAgIGdhbW1hICs9IGRlbHRhOwogICAg
Z2FtbWEgKz0gKCBBUkMgLyAodGlsZUNvdW50LngrQTApKSAqIHRpbGVDb29yZC54OwoKICAgIGZs
b2F0IGYgPSAgKCh0aWxlQ29vcmQueSt0aWxlQ291bnQueS8yLitzY3JvbGxlcikgLyB0aWxlQ291
bnQueCkgOwogICAgZmxvYXQgZyA9ICAoKHRpbGVDb29yZC55K3RpbGVDb3VudC55LzIuKSAvIHRp
bGVDb3VudC54KSA7CiAgICAKICAgIAogICAgCiAgICBnYW1tYSArPSBmICogMS4gKiBjb3MoLjUq
bm93KTsKICAgIGdhbW1hICs9IGYgKiAyLiAqIHNpbiguNSpub3cpOwogICAgZ2FtbWEgKz0gbm93
LzEyLjsKICAgIAogICAgZmxvYXQgcjAgPSBSMCArIChyYWQgKiByb3cgKSA7CiAgICBmbG9hdCBy
MSA9IHIwICsgcmFkOwoKCiAgICByMCArPSA0LiAqICgwLjUgKiBjb3MoMS4wKm5vdyArIGYqUEkg
KSArIDAuNSk7CiAgICByMSArPSA0LiAqICgwLjUgKiBjb3MoMS4wKm5vdyArIGYqUEkgKSArIDAu
NSk7CgogICAgcjAgKz0gMi4gKiAoMC41ICogY29zKDIuKm5vdyArIC0xLi8oZisuMikqUEkgKyBQ
SS80LikgKyAwLjUpOwogICAgcjEgKz0gMi4gKiAoMC41ICogY29zKDIuKm5vdyArIC0xLi8oZisu
MikqUEkgKyBQSS80LikgKyAwLjUpOwogICAgCgogICAgcjAgLT0gKHNjcm9sbGVyKnJhZCk7CiAg
ICByMSAtPSAoc2Nyb2xsZXIqcmFkKTsKICAgIAogICAgdmVjMiBvZmYgPSB2ZWMyKDAuKTsKICAg
IG9mZi54ID0gZio0LiAqIGNvcygwLjI1Km5vdyk7CiAgICBvZmYueSA9IGYqNC4gKiBzaW4oMC41
Km5vdyk7CiAgICAKCiAgICB2ZWMyIEEgPSB2ZWMyKCBvZmYueCtjb3MoZ2FtbWErYWxwaGEpKnIw
LCBvZmYueStzaW4oZ2FtbWErYWxwaGEpKnIwKTsKICAgIHZlYzIgQiA9IHZlYzIoIG9mZi54K2Nv
cyhnYW1tYSthbHBoYSkqcjEsIG9mZi55K3NpbihnYW1tYSthbHBoYSkqcjEpOwogICAgdmVjMiBD
ID0gdmVjMiggb2ZmLngrY29zKGdhbW1hICAgICAgKSpyMSwgb2ZmLnkrc2luKGdhbW1hICAgICAg
KSpyMSk7CiAgICB2ZWMyIEQgPSB2ZWMyKCBvZmYueCtjb3MoZ2FtbWEgICAgICApKnIwLCBvZmYu
eStzaW4oZ2FtbWEgICAgICApKnIwKTsKICAgIAogICAKICAgCiAgICBpZiAgICAgICAgKCBwb3Mu
eCA+IDAuICYmIHBvcy55ID4gMC4gKSB7CiAgICAgICAgcG9zLnh5ID0gQTsKICAgIH0gZWxzZSBp
ZiAoIHBvcy54ID4gMC4gJiYgcG9zLnkgPCAwLiApIHsKICAgICAgICBwb3MueHkgPSBCOwogICAg
fSBlbHNlIGlmICggcG9zLnggPCAwLiAmJiBwb3MueSA+IDAuICkgewogICAgICAgIHBvcy54eSA9
IEQ7CiAgICB9IGVsc2UgaWYgKCBwb3MueCA8IDAuICYmIHBvcy55IDwgMC4gKSB7CiAgICAgICAg
cG9zLnh5ID0gQzsKICAgIH0KICAgIAogICAgZmxvYXQgem9vbSA9IDEuLygyLip0aWxlQ291bnQu
eSk7CgogICAgbWF0NCBtZGwgPSBtYXQ0KDEuMCk7CiAgICBtZGxbMF1bMF0gPSB6b29tOwogICAg
bWRsWzFdWzFdID0gem9vbTsKICAgIAoKICAgIGdsX1Bvc2l0aW9uID0gcHJvamVjdGlvbiAqIHZp
ZXcgKiBtZGwgKiBwb3M7Cn0KCg==
`,

	"set/def.vert": `
CnVuaWZvcm0gbWF0NCBwcm9qZWN0aW9uOwp1bmlmb3JtIG1hdDQgdmlldzsKdW5pZm9ybSBtYXQ0
IG1vZGVsOwoKCnVuaWZvcm0gZmxvYXQgbm93Owp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCnVu
aWZvcm0gZmxvYXQgd29yZENvdW50Owp1bmlmb3JtIGZsb2F0IHdvcmRJbmRleDsKCnVuaWZvcm0g
ZmxvYXQgd29yZFdpZHRoOwp1bmlmb3JtIGZsb2F0IHdvcmRGYWRlcjsKdW5pZm9ybSBmbG9hdCB3
b3JkVmFsdWU7Cgp1bmlmb3JtIGZsb2F0IGNoYXJDb3VudDsKCnVuaWZvcm0gZmxvYXQgc2NyZWVu
UmF0aW87CgphdHRyaWJ1dGUgdmVjMyB2ZXJ0ZXg7CmF0dHJpYnV0ZSB2ZWMyIHRleENvb3JkOwph
dHRyaWJ1dGUgZmxvYXQgY2hhckluZGV4OwphdHRyaWJ1dGUgZmxvYXQgY2hhck9mZnNldDsKCnZh
cnlpbmcgdmVjNCB2VGV4Q29vcmQ7CnZhcnlpbmcgdmVjNCB2UG9zaXRpb247CnZhcnlpbmcgZmxv
YXQgdkNoYXJJbmRleDsKCmJvb2wgREVCVUcgPSBkZWJ1Z0ZsYWcgPiAwLjA7CgpmbG9hdCBQSSAg
PSAzLjE0MTU5MjY1MzU4OTc5MzIzODQ2MjY0MzM4MzI3OTUwMjg4NDE5NzE2OTM5OTM3NTEwNTgy
MDk3NDk0NDU5MjA7CmZsb2F0IFRBVSA9IDYuMjgzMTg1MzA3MTc5NTg2NDc2OTI1Mjg2NzY2NTU5
MDA1NzY4Mzk0MzM4Nzk4NzUwMjExNjQxOTQ5ODg5MTg0MDsKCmZsb2F0IElkZW50aXR5KGZsb2F0
IHgpIHsgcmV0dXJuIHg7IH0KZmxvYXQgRWFzZUluRWFzZU91dChmbG9hdCB4KSB7IHJldHVybiAt
MC41ICogY29zKCB4ICogUEkgKSArIDAuNTsgfQoKZmxvYXQgRWFzZU91dChmbG9hdCB4KSB7IHJl
dHVybiBjb3MoeCpQSS8yLiArIDMuKlBJLzIuICk7IH0KZmxvYXQgRWFzZUluKGZsb2F0IHgpIHsg
cmV0dXJuICAtMS4gKiBjb3MoeCpQSS8yLiApICsgMS4gIDsgfQoKCnZvaWQgbWFpbigpIHsKICAg
IGZsb2F0IGZhZGVyID0gd29yZEZhZGVyOwogICAgCiAgICB2ZWM0IHBvcyA9IHZlYzQodmVydGV4
LDEpOwoKICAgIHZQb3NpdGlvbiA9ICAgIHBvczsKICAgIHZUZXhDb29yZCA9IHZlYzQodGV4Q29v
cmQueHksMS4sMS4pOwogICAgdkNoYXJJbmRleCA9IGNoYXJJbmRleDsKCiAgICBib29sIG9kZCA9
IG1vZCh3b3JkSW5kZXgsMi4wKSA9PSAxLjA7CiAgICAKICAgIGlmIChvZGQpIHsKICAgICAgICBw
b3MueCArPSB3b3JkV2lkdGgvMi47CiAgICB9IGVsc2UgewogICAgICAgIHBvcy54IC09IHdvcmRX
aWR0aC8yLjsKICAgIH0KICAgIHBvcy55IC09ICh3b3JkSW5kZXggLSB3b3JkQ291bnQvMi4pLzIu
OwogICAgCgogICAgICAgIAogICAgZ2xfUG9zaXRpb24gPSBwcm9qZWN0aW9uICogdmlldyAqIG1v
ZGVsICogcG9zOwp9Cgo=
`,

	"set/def.frag": `
CnVuaWZvcm0gZmxvYXQgbm93Owp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKdW5pZm9ybSBzYW1w
bGVyMkQgdGV4dHVyZTsKCnVuaWZvcm0gZmxvYXQgd29yZE1heDsKdW5pZm9ybSBmbG9hdCB3b3Jk
TWF4V2lkdGg7CnVuaWZvcm0gZmxvYXQgd29yZEZhZGVyOwp1bmlmb3JtIGZsb2F0IHdvcmRJbmRl
eDsKdW5pZm9ybSBmbG9hdCB3b3JkQ291bnQ7Cgp2YXJ5aW5nIHZlYzQgdlBvc2l0aW9uOwp2YXJ5
aW5nIHZlYzQgdlRleENvb3JkOwoKYm9vbCBERUJVRyAgICA9IGRlYnVnRmxhZyA+IDAuMDsKCmZs
b2F0IFBJICA9IDMuMTQxNTkyNjUzNTg5NzkzMjM4NDYyNjQzMzgzMjc5NTAyODg0MTk3MTY5Mzk5
Mzc1MTA1ODIwOTc0OTQ0NTkyMDsKZmxvYXQgVEFVID0gNi4yODMxODUzMDcxNzk1ODY0NzY5MjUy
ODY3NjY1NTkwMDU3NjgzOTQzMzg3OTg3NTAyMTE2NDE5NDk4ODkxODQwOwoKZmxvYXQgSWRlbnRp
dHkoZmxvYXQgeCkgeyByZXR1cm4geDsgfQpmbG9hdCBFYXNlSW5FYXNlT3V0KGZsb2F0IHgpIHsg
cmV0dXJuIC0wLjUgKiBjb3MoIHggKiBQSSApICsgMC41OyB9CgpmbG9hdCBFYXNlT3V0KGZsb2F0
IHgpIHsgcmV0dXJuICAgICAgICBjb3MoeCpQSS8yLiArIDMuKlBJLzIuICk7ICAgICAgICB9CmZs
b2F0IEVhc2VJbihmbG9hdCB4KSB7ICByZXR1cm4gIC0xLiAqIGNvcyh4KlBJLzIuICAgICAgICAg
ICAgKSArIDEuICA7IH0KCgp2b2lkIG1haW4oKSB7CiAgICB2ZWM0IGNvbDsKICAgIHZlYzQgcG9z
ID0gdlBvc2l0aW9uOwogICAgdmVjNCB0ZXggPSB2VGV4Q29vcmQ7CiAgICAKICAgIAogICAgaWYg
KDA9PTApIHsKICAgICAgICB2ZWMyIFRleHR1cmVTaXplID0gdmVjMihhYnModlBvc2l0aW9uLngp
LGFicyh2UG9zaXRpb24ueSkpOwogICAgICAgIHRleC54eSA9IHZUZXhDb29yZC54eSArIChnbF9Q
b2ludENvb3JkICogVGV4dHVyZVNpemUpOwogICAgICAgIGNvbCA9IHRleHR1cmUyRCh0ZXh0dXJl
LHRleC54eSk7CiAgICAKICAgIH0gZWxzZSB7CiAgICAKICAgICAgICBjb2wgPSB0ZXh0dXJlMkRQ
cm9qKHRleHR1cmUsIHRleCk7CiAgICB9CiAgICAKICAgIAogICAgZmxvYXQgeiA9IDEuOwogICAg
aWYgKGNvbC5hID4gMC4wICkgewogICAgICAgIGlmICggd29yZEZhZGVyIDwgLjEyNSApIHsKICAg
ICAgICAgICAgeiA9IEVhc2VPdXQod29yZEZhZGVyICogOC4pOwogICAgICAgIH0gZWxzZSBpZiAo
d29yZEZhZGVyIDwgLjc1ICkgewogICAgICAgICAgICB6ID0gMS4wOwogICAgICAgIH0gZWxzZSB7
CiAgICAgICAgICAgIHogPSAxLiAtIEVhc2VPdXQoMy4vNC4gKyB3b3JkRmFkZXIgKiA0Lik7CiAg
ICAgICAgfQogICAgfQogICAgCiAgICBjb2wuYSA9IGNvbC5hICogejsKICAgIAoKICAgIAovLyAg
ICBjb2wucmdiICo9IGNvbC5hOwogICAgCiAgICBpZiAoREVCVUcpIHsgCiAgICAgICAgY29sLnJn
YiA9IHZlYzMoMC4sMC4sMS4pOwogICAgICAgIGlmICggd29yZEZhZGVyID4gLjc1ICkgewogICAg
ICAgICAgICBjb2wuZyA9IDEuOwogICAgICAgIH0gZWxzZSBpZiAod29yZEZhZGVyID4gLjEyNSAp
IHsKICAgICAgICAgICAgY29sLnIgPSAxLjsKICAgICAgICB9IGVsc2UgewogICAgICAgIH0KICAg
ICAgICBjb2wuYSA9IDEuMDsKICAgIH0gCgoKICAgIGlmICghZ2xfRnJvbnRGYWNpbmcpIHsgY29s
LmEgLz0gNC47IH0KCiAgICBnbF9GcmFnQ29sb3IgPSBjb2w7CiAgICAKfQo=
`,

	"set/scroll.vert": `
CnVuaWZvcm0gbWF0NCBwcm9qZWN0aW9uOwp1bmlmb3JtIG1hdDQgdmlldzsKdW5pZm9ybSBtYXQ0
IG1vZGVsOwoKCnVuaWZvcm0gZmxvYXQgbm93Owp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCnVu
aWZvcm0gZmxvYXQgd29yZENvdW50Owp1bmlmb3JtIGZsb2F0IHdvcmRJbmRleDsKCnVuaWZvcm0g
ZmxvYXQgd29yZFdpZHRoOwp1bmlmb3JtIGZsb2F0IHdvcmRGYWRlcjsKdW5pZm9ybSBmbG9hdCB3
b3JkVmFsdWU7Cgp1bmlmb3JtIGZsb2F0IGNoYXJDb3VudDsKCnVuaWZvcm0gZmxvYXQgc2NyZWVu
UmF0aW87CgphdHRyaWJ1dGUgdmVjMyB2ZXJ0ZXg7CmF0dHJpYnV0ZSB2ZWMyIHRleENvb3JkOwph
dHRyaWJ1dGUgZmxvYXQgY2hhckluZGV4OwphdHRyaWJ1dGUgZmxvYXQgY2hhck9mZnNldDsKCnZh
cnlpbmcgdmVjNCB2VGV4Q29vcmQ7CnZhcnlpbmcgdmVjNCB2UG9zaXRpb247CnZhcnlpbmcgZmxv
YXQgdkNoYXJJbmRleDsKCgoKYm9vbCBERUJVRyA9IGRlYnVnRmxhZyA+IDAuMDsKCgoKZmxvYXQg
UEkgID0gMy4xNDE1OTI2NTM1ODk3OTMyMzg0NjI2NDMzODMyNzk1MDI4ODQxOTcxNjkzOTkzNzUx
MDU4MjA5NzQ5NDQ1OTIwOwpmbG9hdCBUQVUgPSA2LjI4MzE4NTMwNzE3OTU4NjQ3NjkyNTI4Njc2
NjU1OTAwNTc2ODM5NDMzODc5ODc1MDIxMTY0MTk0OTg4OTE4NDA7CgpmbG9hdCBJZGVudGl0eShm
bG9hdCB4KSB7IHJldHVybiB4OyB9CmZsb2F0IEVhc2VJbkVhc2VPdXQoZmxvYXQgeCkgeyByZXR1
cm4gLTAuNSAqIGNvcyggeCAqIFBJICkgKyAwLjU7IH0KCmZsb2F0IEVhc2VPdXQoZmxvYXQgeCkg
eyByZXR1cm4gY29zKHgqUEkvMi4gKyAzLipQSS8yLiApOyB9CmZsb2F0IEVhc2VJbihmbG9hdCB4
KSB7IHJldHVybiAgLTEuICogY29zKHgqUEkvMi4gKSArIDEuICA7IH0KCgp2b2lkIG1haW4oKSB7
CiAgICB2VGV4Q29vcmQueHkgPSB0ZXhDb29yZC54eTsKICAgIAogICAgdmVjNCBwb3M7CiAgICBw
b3MgPSB2ZWM0KHZlcnRleCwxKTsKICAgIAogICAgdmVjNCB0ZXg7CiAgICB0ZXggPSB2ZWM0KHRl
eENvb3JkLnh5LCAxLiwgMS4pOwoKICAgIHBvcy55IC09IHdvcmRDb3VudC8yLjsKICAgIHBvcy55
ICs9ICh3b3JkSW5kZXgvd29yZENvdW50KSAqIHdvcmRDb3VudDsKICAgIHBvcy55ICs9IDAuNTsK
CgogICAgcG9zLnggKz0gKHdvcmRDb3VudC8yLikgKiBzY3JlZW5SYXRpbyA7CiAgICBwb3MueCAt
PSB3b3JkRmFkZXIgKiAoIHdvcmRDb3VudCAqIHNjcmVlblJhdGlvKTsKICAgIAogICAgdlBvc2l0
aW9uID0gcG9zOwogICAgdlRleENvb3JkID0gdGV4OwogICAgdkNoYXJJbmRleCA9IGNoYXJJbmRl
eDsKICAgIAogICAgZ2xfUG9zaXRpb24gPSBwcm9qZWN0aW9uICogdmlldyAqIG1vZGVsICogcG9z
Owp9Cgo=
`,

	"set/field.vert": `
CnVuaWZvcm0gbWF0NCBwcm9qZWN0aW9uOwp1bmlmb3JtIG1hdDQgdmlldzsKdW5pZm9ybSBtYXQ0
IG1vZGVsOwoKCnVuaWZvcm0gZmxvYXQgbm93Owp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKCnVu
aWZvcm0gZmxvYXQgd29yZENvdW50Owp1bmlmb3JtIGZsb2F0IHdvcmRJbmRleDsKCnVuaWZvcm0g
ZmxvYXQgd29yZFdpZHRoOwp1bmlmb3JtIGZsb2F0IHdvcmRGYWRlcjsKdW5pZm9ybSBmbG9hdCB3
b3JkVmFsdWU7Cgp1bmlmb3JtIGZsb2F0IGNoYXJDb3VudDsKCnVuaWZvcm0gZmxvYXQgc2NyZWVu
UmF0aW87CgphdHRyaWJ1dGUgdmVjMyB2ZXJ0ZXg7CmF0dHJpYnV0ZSB2ZWMyIHRleENvb3JkOwph
dHRyaWJ1dGUgZmxvYXQgY2hhckluZGV4OwphdHRyaWJ1dGUgZmxvYXQgY2hhck9mZnNldDsKCnZh
cnlpbmcgdmVjNCB2VGV4Q29vcmQ7CnZhcnlpbmcgdmVjNCB2UG9zaXRpb247CnZhcnlpbmcgZmxv
YXQgdkNoYXJJbmRleDsKCgpib29sIERFQlVHID0gZGVidWdGbGFnID4gMC4wOwoKbWF0NCBzY2Fs
ZU1hdHJpeChmbG9hdCBhKSB7CiAgICByZXR1cm4gbWF0NCgKICAgICAgICAgYSwgMC4sIDAuLCAw
LiwKICAgICAgICAwLiwgIGEsIDAuLCAwLiwKICAgICAgICAwLiwgMC4sIDEuLCAwLiwKICAgICAg
ICAwLiwgMC4sIDAuLCAxLgogICAgKTsKfQoKCmZsb2F0IFBJICA9IDMuMTQxNTkyNjUzNTg5Nzkz
MjM4NDYyNjQzMzgzMjc5NTAyODg0MTk3MTY5Mzk5Mzc1MTA1ODIwOTc0OTQ0NTkyMDsKZmxvYXQg
VEFVID0gNi4yODMxODUzMDcxNzk1ODY0NzY5MjUyODY3NjY1NTkwMDU3NjgzOTQzMzg3OTg3NTAy
MTE2NDE5NDk4ODkxODQwOwoKZmxvYXQgSWRlbnRpdHkoZmxvYXQgeCkgeyByZXR1cm4geDsgfQpm
bG9hdCBFYXNlSW5FYXNlT3V0KGZsb2F0IHgpIHsgcmV0dXJuIC0wLjUgKiBjb3MoIHggKiBQSSAp
ICsgMC41OyB9CgpmbG9hdCBFYXNlT3V0KGZsb2F0IHgpIHsgcmV0dXJuIGNvcyh4KlBJLzIuICsg
My4qUEkvMi4gKTsgfQpmbG9hdCBFYXNlSW4oZmxvYXQgeCkgeyByZXR1cm4gIC0xLiAqIGNvcyh4
KlBJLzIuICkgKyAxLiAgOyB9CgoKdm9pZCBtYWluKCkgewogICAgZmxvYXQgZmFkZXIgPSB3b3Jk
RmFkZXI7CiAgICAKICAgIHZlYzQgcG9zOwogICAgcG9zID0gdmVjNCh2ZXJ0ZXgsMSk7CgogICAg
ZmxvYXQgYTsKCiAgICBwb3MueCAqPSB3b3JkQ291bnQvMTAuOwogICAgcG9zLnkgKj0gd29yZENv
dW50LzEwLjsKICAgIAogICAgZmxvYXQgcnggPSAwLjUgKiB3b3JkQ291bnQvMi47CiAgICBmbG9h
dCByeSA9IHdvcmRDb3VudC8yLjsKICAgIAogICAgZmxvYXQgaWR4ID0gd29yZEluZGV4L3dvcmRD
b3VudDsKICAgIAogICAgZmxvYXQgdyA9IDIuKlRBVTsKICAgIGZsb2F0IG8gPSBQSS8yLiArIFBJ
LzQuOwogICAgYSA9IGlkeCAqIHcgKyBvOwogICAgcG9zLnggKz0gcnggKiBjb3MoYSk7CiAgICBw
b3MueSArPSByeSAqIHNpbihhKTsKICAgICAgICAKCiAgICBpZiAoaWR4ID49IDAuNSApIHsKICAg
ICAgICBwb3MueCArPSB3b3JkQ291bnQvMi47CiAgICB9IGVsc2UgewogICAgICAgIHBvcy54IC09
IHdvcmRDb3VudC8yLjsKICAgIH0KCgogICAgcG9zLnogLT0gd29yZENvdW50OyAKICAgIHBvcy56
ICs9IGZhZGVyICogMS41ICogd29yZENvdW50OwoKICAgIG1hdDQgdiA9IHZpZXc7Ci8vICAgIHYg
PSB2aWV3ICogc2NhbGVNYXRyaXgoIDEuICApOwoKICAgIHZQb3NpdGlvbiA9IHBvczsKICAgIHZU
ZXhDb29yZCA9IHZlYzQodGV4Q29vcmQueHksMS4sMS4pOwogICAgdkNoYXJJbmRleCA9IGNoYXJJ
bmRleDsKICAgIAogICAgICAgIAogICAgZ2xfUG9zaXRpb24gPSBwcm9qZWN0aW9uICogdiAqIG1v
ZGVsICogcG9zOwp9Cgo=
`,

	"set/debug.frag": `
CnVuaWZvcm0gZmxvYXQgbm93Owp1bmlmb3JtIGZsb2F0IGRlYnVnRmxhZzsKdW5pZm9ybSBzYW1w
bGVyMkQgdGV4dHVyZTsKCnVuaWZvcm0gZmxvYXQgd29yZENvdW50Owp1bmlmb3JtIGZsb2F0IHdv
cmRNYXhXaWR0aDsKdW5pZm9ybSBmbG9hdCB3b3JkRmFkZXI7CnVuaWZvcm0gZmxvYXQgd29yZElu
ZGV4Owp1bmlmb3JtIGZsb2F0IHdvcmRWYWx1ZTsKCnZhcnlpbmcgdmVjNCB2UG9zaXRpb247CnZh
cnlpbmcgdmVjNCB2VGV4Q29vcmQ7CgoKCmJvb2wgREVCVUcgICAgPSBkZWJ1Z0ZsYWcgPiAwLjA7
CgpmbG9hdCBQSSAgPSAzLjE0MTU5MjY1MzU4OTc5MzIzODQ2MjY0MzM4MzI3OTUwMjg4NDE5NzE2
OTM5OTM3NTEwNTgyMDk3NDk0NDU5MjA7CmZsb2F0IFRBVSA9IDYuMjgzMTg1MzA3MTc5NTg2NDc2
OTI1Mjg2NzY2NTU5MDA1NzY4Mzk0MzM4Nzk4NzUwMjExNjQxOTQ5ODg5MTg0MDsKCmZsb2F0IElk
ZW50aXR5KGZsb2F0IHgpIHsgcmV0dXJuIHg7IH0KZmxvYXQgRWFzZUluRWFzZU91dChmbG9hdCB4
KSB7IHJldHVybiAtMC41ICogY29zKCB4ICogUEkgKSArIDAuNTsgfQoKZmxvYXQgRWFzZU91dChm
bG9hdCB4KSB7IHJldHVybiAgICAgICAgY29zKHgqUEkvMi4gKyAzLipQSS8yLiApOyAgICAgICAg
fQpmbG9hdCBFYXNlSW4oZmxvYXQgeCkgeyAgcmV0dXJuICAtMS4gKiBjb3MoeCpQSS8yLiAgICAg
ICAgICAgICkgKyAxLiAgOyB9CgoKdm9pZCBtYWluKCkgewogICAgdmVjNCBjb2w7CiAgICB2ZWM0
IHBvcyA9IHZQb3NpdGlvbjsKICAgIHZlYzQgdGV4ID0gdlRleENvb3JkOwogICAgCiAgICAKLy8g
ICAgY29sID0gdGV4dHVyZTJEUHJvaih0ZXh0dXJlLCB0ZXgpOwogICAgCiAgICBpZiAoIURFQlVH
KSB7CiAgICAgICAgY29sLmEgPSAwLjU7CiAgICAgICAgY29sLnJnID0gdGV4Lnh5OyAgICAKICAg
ICAgICBjb2wuYiA9IDEuOy8vdGV4Lno7CiAgICB9CiAgICAKICAgIGlmIChERUJVRykgewogICAg
ICAgIGNvbC5hID0gMS4wOwogICAgICAgIGNvbC5yZyA9IHRleC54eTsgICAgCiAgICAgICAgY29s
LmIgPSAxLjsvL3RleC56OwogICAgfQogICAgCiAgICAKCgogICAgaWYgKCFnbF9Gcm9udEZhY2lu
ZykgeyBjb2wuYSAvPSA0LjsgfQoKICAgIGdsX0ZyYWdDb2xvciA9IGNvbDsKICAgIAp9Cg==
`,

	"mask/def.vert": `
CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwoKYXR0cmlidXRlIHZlYzMgdmVydGV4OwphdHRyaWJ1
dGUgdmVjMiB0ZXhDb29yZDsKCnZhcnlpbmcgdmVjMiB2VGV4Q29vcmQ7Cgpib29sIERFQlVHID0g
ZGVidWdGbGFnID4gMC4wOwoKCnZvaWQgbWFpbigpIHsKICAgIHZUZXhDb29yZCA9IHRleENvb3Jk
OwogICAgZ2xfUG9zaXRpb24gPSB2ZWM0KHZlcnRleCwxLik7Cn0K
`,

	"mask/def.frag": `
CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwp1bmlmb3JtIGZsb2F0IHNjcmVlblJhdGlvOwoKdmFy
eWluZyB2ZWMyIHZUZXhDb29yZDsKCgpib29sIERFQlVHID0gZGVidWdGbGFnID4gMC4wOwpmbG9h
dCB3ID0gMC4wMDU7Cgpib29sIGdyaWQodmVjMiBwb3MpIHsKCiAgICBmb3IgKGZsb2F0IGQgPSAt
Mi4wOyBkPD0yLjA7IGQrPTAuNSkgewogICAgICAgIGlmIChhYnMocG9zLnkgLSBkKSAtIHcgPD0g
MC4wICkgeyByZXR1cm4gdHJ1ZTsgfQogICAgICAgIGlmIChhYnMocG9zLnggLSBkKSAtIHcgPD0g
MC4wICkgeyByZXR1cm4gdHJ1ZTsgfQogICAgfQogICAgCiAgICByZXR1cm4gZmFsc2U7Cn0KCgp2
b2lkIG1haW4oKSB7CiAgICB2ZWMzIGNvbCA9IHZlYzMoMC4wLDAuMCwwLjApOwogICAgZmxvYXQg
YSA9IDEuMDsKCgogICAgdmVjMiBwb3MgPSB2VGV4Q29vcmQ7CgoKICAgIGlmICggdHJ1ZSAmJiBn
cmlkKHBvcykgKSB7CiAgICAgICAgZmxvYXQgZ3JheSA9IDAuNTsKICAgICAgICBjb2wgPSBncmF5
ICogdmVjMygxLiwxLiwxLik7CiAgICB9CgogICAgCiAgICAKICAgIGdsX0ZyYWdDb2xvciA9IHZl
YzQoY29sLnJnYiwgYSk7Cn0KCgo=
`,

	"mask/mask.frag": `
CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwp1bmlmb3JtIGZsb2F0IHNjcmVlblJhdGlvOwoKdmFy
eWluZyB2ZWMyIHZUZXhDb29yZDsKCmJvb2wgREVCVUcgPSBkZWJ1Z0ZsYWcgPiAwLjA7CgpmbG9h
dCBNQVgoZmxvYXQgYSwgZmxvYXQgYikgeyAKICAgIGlmIChhPj1iKSAKICAgICAgICByZXR1cm4g
YTsgCiAgICBlbHNlIAogICAgICAgIHJldHVybiBiOyAKfQoKZmxvYXQgbWFzayh2ZWMyIHApIHsg
CiAgICBmbG9hdCB4ID0gcC54OyBmbG9hdCB5ID0gcC55OwogICAgZmxvYXQgZmYgPSAxLjsKICAg
IGZsb2F0IEEgPSA3LiAqIGZmOwogICAgZmxvYXQgQiA9IDguICogZmY7IAogICAgaWYgKGFicyh4
KSA+PSBBL0IgJiYgYWJzKHkpID49IEEvQikKICAgICAgICByZXR1cm4gMS4wIC0gKCBCICogTUFY
KGFicyh4KSxhYnMoeSkpIC0gQSApOwogICAgZWxzZSBpZiAoYWJzKHgpID49IEEvQikgCiAgICAg
ICAgcmV0dXJuIDEuMCAtICggQiAqIGFicyh4KSAtIEEpOwogICAgZWxzZSBpZiAoYWJzKHkpID49
IEEvQikKICAgICAgICByZXR1cm4gMS4wIC0gKCBCICogYWJzKHkpIC0gQSk7CiAgICByZXR1cm4g
MS4wOwp9CgoKCnZvaWQgbWFpbigpIHsKICAgIHZlYzMgY29sID0gdmVjMygwLjAsMC4wLDAuMCk7
CiAgICBmbG9hdCBhID0gMS47CgoKICAgIHZlYzIgcG9zID0gdlRleENvb3JkOwoKICAgIGEgPSBt
YXNrKHZlYzIocG9zLngvc2NyZWVuUmF0aW8scG9zLnkpKTsKCiAgICBpZiAoREVCVUcpIHsKICAg
ICAgICBjb2wucmdiID0gIHZlYzMoICAxLiAtIG1hc2sodmVjMihwb3MueC9zY3JlZW5SYXRpbyxw
b3MueSkpICk7CiAgICB9CiAgICAKICAgIAogICAgZ2xfRnJhZ0NvbG9yID0gdmVjNChjb2wucmdi
LCBhKTsKfQo=
`,

	"mask/debug.frag": `
CnVuaWZvcm0gZmxvYXQgZGVidWdGbGFnOwp1bmlmb3JtIGZsb2F0IHNjcmVlblJhdGlvOwoKdmFy
eWluZyB2ZWMyIHZUZXhDb29yZDsKCgpib29sIERFQlVHID0gZGVidWdGbGFnID4gMC4wOwpmbG9h
dCB3ID0gMC4wMDI7Cgpib29sIG1ham9yKHZlYzIgcG9zKSB7CgogICAgZm9yIChmbG9hdCBkID0g
LTIuMDsgZDw9Mi4wOyBkKz0xLjApIHsKICAgICAgICBpZiAoYWJzKHBvcy55IC0gZCkgLSB3IDw9
IDAuMCApIHsgcmV0dXJuIHRydWU7IH0KICAgICAgICBpZiAoYWJzKHBvcy54IC0gZCkgLSB3IDw9
IDAuMCApIHsgcmV0dXJuIHRydWU7IH0KICAgIH0KICAgIAogICAgcmV0dXJuIGZhbHNlOwp9Cgpi
b29sIG1pbm9yKHZlYzIgcG9zKSB7CgogICAgZm9yIChmbG9hdCBkID0gLTIuMDsgZDw9Mi4wOyBk
Kz0wLjI1KSB7CiAgICAgICAgaWYgKGFicyhwb3MueSAtIGQpIC0gdyA8PSAwLjAgKSB7IHJldHVy
biB0cnVlOyB9CiAgICAgICAgaWYgKGFicyhwb3MueCAtIGQpIC0gdyA8PSAwLjAgKSB7IHJldHVy
biB0cnVlOyB9CiAgICB9CiAgICAKICAgIHJldHVybiBmYWxzZTsKfQoKCgp2b2lkIG1haW4oKSB7
CiAgICB2ZWMzIGNvbCA9IHZlYzMoMC4wLDAuMCwwLjApOwogICAgZmxvYXQgYSA9IDEuMDsKCgog
ICAgdmVjMiBwb3MgPSB2VGV4Q29vcmQ7CgoKICAgIGlmICggdHJ1ZSAmJiBtYWpvcihwb3MpICkg
ewogICAgICAgIGZsb2F0IGdyYXkgPSAxLjA7CiAgICAgICAgY29sID0gZ3JheSAqIHZlYzMoMS4s
MS4sMS4pOwogICAgfSBlbHNlIGlmICggdHJ1ZSAmJiBtaW5vcihwb3MpICkgewogICAgICAgIGZs
b2F0IGdyYXkgPSAwLjU7CiAgICAgICAgY29sID0gZ3JheSAqIHZlYzMoMS4sMS4sMS4pOwogICAg
fQoKICAgIGdsX0ZyYWdDb2xvciA9IHZlYzQoY29sLnJnYiwgYSk7Cn0KCgoK
`,
}
