package facade

var Asset = map[string]string{

	"README": `
CiMgRkFDQURFIGJ5IEZFRURGQUNFLkNPTQoKICAgIApGQUNBREUgCgoKCiMjIEV4YW1wbGVzCgoK
IyMjIFNldHVwIEFsaWFzCgogICAgCiAgICBhbGlhcyBmY2Q9J25jIC1OIGxvY2FsaG9zdCA0MDQ1
JyAjIGZvciBsaW51eAogICAgYWxpYXMgZmNkPSduYyBsb2NhbGhvc3QgNDA0NScgICAgIyBmb3Ig
bWFjL2JzZAoKCgoKCgoKIyMgRGVjb3JhdGl2ZSAvIFdhbGxwYXBlcgoKCiMjIyMgUGhyYWNrCgog
ICAgZmFjYWRlIGNvbmYgbGluZXMgLXc9ODAgLWg9MjUgLXZlcnQ9cm9sbAogICAgY3VybCAtTCBo
dHRwOi8vcGhyYWNrLm9yZy9hcmNoaXZlcy90Z3ovcGhyYWNrNDkudGFyLmd6IFwKICAgICB8IHRh
ciB4ZnogL2Rldi9zdGRpbiAuLzE0LnR4dCAtTyBcCiAgICAgfCB3aGlsZSByZWFkIC1yIGxpbmU7
IGRvIGVjaG8gIiRsaW5lIjsgc2xlZXAgLjk7IGRvbmUgXAogICAgIHwgZmNkCgoKIyMjIyBSRkMK
CiAgICBmYWNhZGUgY29uZiBsaW5lcyAtdz03MiAtaD0xNiAtdmVydD1yb3dzCiAgICBjdXJsIC1M
IGh0dHBzOi8vdG9vbHMuaWV0Zi5vcmcvcmZjL3JmYzI0NjAudHh0IFwKICAgICB8IHdoaWxlIHJl
YWQgLXIgbGluZTsgZG8gZWNobyAiJGxpbmUiOyBzbGVlcCAuOTsgZG9uZSBcCiAgICAgfCBmY2QK
CiMjIyMgLm5mbyAKCiAgICBmYWNhZGUgY29uZiBsaW5lcyAtdz04MCAtaD0yNSAtdmVydD13YXZl
IC1mb250IGFkb3JlNjQKICAgIGN1cmwgLUwgaHR0cHM6Ly9jb250ZW50LnBvdWV0Lm5ldC9maWxl
cy9uZm9zLzAwMDEyLzAwMDEyMDMxLnR4dCBcCiAgICAgfCB3aGlsZSByZWFkIC1yIGxpbmU7IGRv
IGVjaG8gIiRsaW5lIjsgc2xlZXAgLjk7IGRvbmUgXAogICAgIHwgZmNkCiAgICAKIyBNYW5wYWdl
cyAgICAKICAgIAogICAgZmFjYWRlIGNvbmYgbGluZXMgLXc9NTAgLWg9MjAgLXZlcnQ9Y3Jhd2wK
ICAgIGV4cG9ydCBNQU5XSURUSD01MCBNQU5QQUdFUj1jYXQKICAgIG1hbiBzc2ggfCB3aGlsZSBy
ZWFkIC1yIGxpbmU7IGRvIGVjaG8gIiRsaW5lIjsgc2xlZXAgLjk7IGRvbmUgfCBmY2QKCgoKCgoj
IyBJbmZvcm1hdGl2ZSAvIFN0YXR1cwoKCgoKIyMjIyBBY2Nlc3MgTG9ncwoKICAgIGZhY2FkZSBz
ZXJ2ZSBsaW5lcyAtdz0xMjAgLWg9OCAtdmVydD1kaXNrIC1tYXNrPW1hc2sgJgogICAgdGFpbCAt
ZiAvdmFyL2xvZy9uZ2lueC9hY2Nlc3MubG9nIHwgZmNkCgoKIyMjIyBDbG9jawoKICAgIGZhY2Fk
ZSBzZXJ2ZSBsaW5lcyAtdmVydD13YXZlIC1oPTIgLXc9MTAgLW1hc2s9bWFzayAtZG93biAtc21v
b3RoPWYgLWZvbnQgb2NyYWV4dCAtem9vbSAuOAogICAgd2hpbGUgdHJ1ZTsgZG8gZGF0ZSAiKyVZ
LSVtLSVkIjsgc2xlZXAgMTsgZGF0ZSAiKyAlSDolTTolUyI7IHNsZWVwIDE7IGRvbmUgfCBmY2QK
CiAgICAKIyMjIyBtdHIKICAgIAogICAgZmFjYWRlIHNlcnZlIHRlcm0gLXZlcnQ9ZGlzawogICAg
ZmFjYWRlIGV4ZWMgdGVybSAgLXc9NjQgLWg9MTYgc3VkbyBtdHIgLW0gMTAgLS1kaXNwbGF5bW9k
ZSAyIGV4YW1wbGUuY29tCgoKIyMjIyB0Y3BkdW1wCgogICAgZmFjYWRlIHNlcnZlIGxpbmVzIC12
ZXJ0IGRyb3AgLWRvd24gLXNwZWVkIC4yIC13IDEyMCAtaCA4IC1tYXNrPW1hc2sKICAgIHRjcGR1
bXAgLWkgdmxhbjUgLW4gLXQgLWwgLXYgZHN0IHBvcnQgNTMgIHwgZmNkCgojIyMjIHRvcAoKICAg
IGZhY2FkZSBleGVjIHRlcm0gLXc9NjQgLWg9MTYgLXZlcnQ9ZGlzayAvdXNyL2Jpbi90b3AKICAg
IAogICAgCgojIyBDb2xsYWJvcmF0aXZlIC8gSW50ZXJhY3RpdmUKCgojIE1hbnBhZ2VzCgogICAg
ZmFjYWRlIGV4ZWMgdGVybSAtdz01MCAtaD0yMCBtYW4gc3NoCgoKCgojIGZyb3R6CgogICAgZmFj
YWRlIGV4ZWMgdGVybSAtdz02NCAtaD0xNiAvcGF0aC90by9mcm90eiAvcGF0aC90by9oaXRjaGhp
a2Vyc19ndWlkZS5ibGIKCiAgICAKICAgIAogICAgCiMgY2xlYXIgZGlzcGxheQoKIyAgICBwcmlu
dGYgJ1wwMzNbODsxNjs2NHQnICMgcmVzaXplIHRlcm1pbmFsCgoKICAgIGNsZWFyIHwgZmNkCiAg
ICAKCg==
`,
}
