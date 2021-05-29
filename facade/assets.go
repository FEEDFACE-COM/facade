package facade

var Asset = map[string]string{

	"README": `
IyBGQUNBREUgYnkgRkVFREZBQ0UuQ09NCiAgICAKX0ZBQ0FERV8gaXMgYSBjcmVhdGl2ZSBjb2RpbmcgdG9vbCB0aGF0IGFsbG93cyB5b3UgdG8gcGlwZSBhbnkgdGV4dCBvbiBzdGRvdXQgZGlyZWN0bHkgb250byB0aGUgd2FsbCBvZiB5b3VyIGhvbWUgLyBvZmZpY2UgLyBoYWNrZXJzcGFjZS4gCgpZb3Ugd2lsbCBuZWVkOgoKLSBBIFJhc3BiZXJyeSBQaSBydW5uaW5nIGBmYWNhZGUgc2VydmVgLCBhbmQgcmVhY2hhYmxlIGJ5IG5ldHdvcmsKLSBBIHByb2plY3RvciBjb25uZWN0ZWQgdG8gdGhlIFJhc3BiZXJyeSBQaSB2aWEgSERNSQotIEFuIGFsaWFzIGBhbGlhcyBmY2Q9J25jIC1OIHJhc3BiZXJyeXBpIDQwNDUnYCBpbiB5b3VyIHNoZWxsCgpUaGVuIGp1c3QgcnVuIGBlY2hvIEZPT0JBUiB8IGZjZGAgb24gYW55IG1hY2hpbmUgdGhhdCBjYW4gcmVhY2ggdGhlIFJhc3BiZXJyeSBQaS4gVGhpcyB3aWxsIHJlbmRlciB0aGUgdGV4dCBfRk9PQkFSXyBvbnRvIHRoZSB3YWxsIGZvciBldmVyeW9uZSB0byBzZWUuIFRoZSB0ZXh0IGlzIGRpc3BsYXllZCBpbiBwbGFpbiBjb25zb2xlIHN0eWxlIGJ5IGRlZmF1bHQsIEZBQ0FERSBhbHNvIHN1cHBvcnRzIGN1c3RvbSBzdHlsZXMgaW4gdGhlIGZvcm0gb2YgT3BlbkdMRVMgc2hhZGVycyB0aGF0IHlvdSBjYW4gY3JlYXRlIGFuZCBhbHRlciBvbiB0aGUgZmx5LgoKCi0tLS0KCgpUaGUgbW90aXZhdGlvbiBmb3IgY3JlYXRpbmcgRkFDQURFIGlzIHR3b2ZvbGQ6CgoqIFdlIHNwZW5kIGEgbG90IG9mIHRpbWUgaW4gdmFyaW91cyBjb25zb2xlcywgdGVybWluYWxzIGFuZCBzaGVsbHMsIGJ1dCBtb3N0IG9mIHRoZSBwZW9wbGUgYXJvdW5kIHVzIG5ldmVyIGdldCB0byBzZWUgdGhlIHRleHQgd2UgaW50ZXJhY3Qgd2l0aC4KCiogVGhlcmUgaXMgYSBsb3Qgb2YgYXdlc29tZSBwbGFpbiB0ZXh0IGluIHRoZSBmb3JtIG9mIFJGQ3MsIC5uZm8gZmlsZXMsIHRleHQgYWR2ZW50dXJlcyBldGMgdGhhdCBsb29rcyBldmVuIGJldHRlciB3aGVuIHRha2VuIG91dCBvZiB0aGUgY29uZmluZXMgb2YgYW4geHRlcm0gd2luZG93IGFuZCBwcm9qZWN0ZWQgb250byBhIHdhbGwuCgoKCgojIyBSZXF1aXJlbWVudHMgCgpGQUNBREUgc2VydmVyIHdvcmtzIG9uCgotIFJhc3BiaWFuIExpdGUgb24gUmFzcGJlcnJ5IFBpIDIKLSBSYXNwYmlhbiBMaXRlIG9uIFJhc3BiZXJyeSBQaSAzCgpGQUNBREUgY2xpZW50IHdvcmtzIG9uCgotIEdOVS9MaW51eCBvbiB4ODYtNjQKLSBPcGVuQlNEIG9uIHg4Ni02NAotIEFwcGxlIE1hY09TIG9uIHg4Ni02NAoKCgoKIyMgU2V0dXAKCiMjIyBTZXR1cCBGQUNBREUgU2VydmVyIG9uIFJhc3BiZXJyeSBQaQoKMS4gUHJlcGFyZSB5b3VyIFJhc3BiZXJyeSBQaSBjb25maWd1cmF0aW9uOgogICAgLSBVc2UgYHJhc3BpLWNvbmZpZ2AgdG8gc2V0IHRoZSBtZW1vcnkgYXZhaWxhYmxlIHRvIHRoZSBHUFUgdG8gXzI1NiBNQl8KICAgIC0gVXNlIGByYXNwaS1jb25maWdgIHRvIHNlbGVjdCB0aGUgX0xlZ2FjeSBub24tR0wgZGVza3RvcCBkcml2ZXJfCiAgICAtIE1ha2Ugc3VyZSB0aGUgWCBXaW5kb3cgU3lzdGVtIGlzIG5vdCBydW5uaW5nCgoyLiBEb3dubG9hZCB0aGUgbGF0ZXN0IHJlbGVhc2UgcGFja2FnZToKCS0gUmFzcGJpYW4gTGl0ZTogYGZhY2FkZS14Lnkuei1saW51eC1hcm0udGd6YAoKMy4gRXh0cmFjdCB0aGUgX2ZhY2FkZV8gYmluYXJ5IGZyb20gdGhlIHJlbGVhc2UgcGFja2FnZToKCS0gYHRhciB4ZnogZmFjYWRlLXgueS56LWxpbnV4LWFybS50Z3pgIAoKNC4gUnVuIEZBQ0FERSBzZXJ2ZXIgZnJvbSBhIGNvbnNvbGUgb3Igc3NoIHNoZWxsOgoJLSBgLi9mYWNhZGUgc2VydmUgLWRgCgkKCVlvdSBzaG91bGQgbm93IHNlZSB0aGUgRkFDQURFIHRpdGxlIHNjcmVlbiBvbiB0aGUgSERNSSBkaXNwbGF5IG9mIHRoZSBSYXNwYmVycnkgUGkuCgo1LiBPbiB5b3VyIHdvcmtzdGF0aW9uLCBjcmVhdGUgdGhlIF9mY2RfIHNoZWxsIGFsaWFzOiAgCiAgIChyZXBsYWNlIF9yYXNwYmVycnlwaV8gd2l0aCB0aGUgaG9zdG5hbWUgb3IgSVAgYWRkcmVzcyBvZiB5b3VyIFJhc3BiZXJyeSBQaSkKCS0gUmFzcGJpYW4sIExpbnV4OiBgYWxpYXMgZmNkPSduYyAtTiByYXNwYmVycnlwaSA0MDQ1J2AKCS0gT3BlbkJTRCwgTWFjT1M6IGBhbGlhcyBmY2Q9J25jIHJhc3BiZXJyeXBpIDQwNDUnYAoJCjYuIFRlc3Qgd2hldGhlciB5b3UgY2FuIHNlbmQgcmF3IHRleHQgdG8gdGhlIEZBQ0FERSBzZXJ2ZXI6CgktIGB3aG9hbWkgfCBmY2RgCgogICBZb3Ugc2hvdWxkIG5vdyBzZWUgeW91ciB1c2VybmFtZSBvbiB0aGUgSERNSSBkaXNwbGF5IG9mIHRoZSBSYXNwYmVycnkgUGkuCgoKIyMjIFNldHVwIEZBQ0FERSBDbGllbnQgb24geW91ciB3b3Jrc3RhdGlvbgoKMS4gRG93bmxvYWQgdGhlIGxhdGVzdCByZWxlYXNlIHBhY2thZ2UgZm9yIHlvdXIgcGxhdGZvcm06CgktIEdOVS9MaW51eDogYGZhY2FkZS14Lnkuei1saW51eC1hbWQ2NC50Z3pgCgktIE9wZW5CU0Q6IGBmYWNhZGUteC55Lnotb3BlbmJzZC1hbWQ2NC50Z3pgCgktIEFwcGxlIE1hY09TOiBgZmFjYWRlLXgueS56LWRhcndpbi1hbWQ2NC50Z3pgCgoyLiBFeHRyYWN0IHRoZSBfZmFjYWRlXyBiaW5hcnkgZnJvbSB0aGUgcmVsZWFzZSBwYWNrYWdlOgoJLSBgdGFyIHhmeiBmYWNhZGUteC55Lnotb3MtYXJjaC50Z3pgIAoKMy4gUnVuIEZBQ0FERSBjbGllbnQ6ICAKICAgKHJlcGxhY2UgX3Jhc3BiZXJyeXBpXyB3aXRoIHRoZSBob3N0bmFtZSBvciBJUCBhZGRyZXNzIG9mIHlvdXIgUmFzcGJlcnJ5IFBpKQoJLSBgZGF0ZSAtLWlzby04NjAxIHwgLi9mYWNhZGUgcGlwZSAtaG9zdCByYXNwYmVycnlwaSBsaW5lcyAtZm9udCBhZG9yZTY0YAoKCVlvdSBzaG91bGQgbm93IHNlZSB5b3VyIHVzZXJuYW1lIGFuZCB0aGUgY3VycmVudCBkYXRlIGluIGEgcGl4ZWwgZm9udC4KCjQuIEV4cGxvcmUgRkFDQURFIG9wdGlvbnM6CgoJLSBgLi9mYWNhZGUgLWhgCgktIGAuL2ZhY2FkZSBjb25mIC1oYAoJLSBgLi9mYWNhZGUgY29uZiBsaW5lcyAtaGAKCgoKCiMjIFNlY3VyaXR5CgotIEZBQ0FERSBwcm92aWRlcyBubyBlbmNyeXB0aW9uLCB0cmFuc3BvcnQgbGF5ZXIgb3Igb3RoZXJ3aXNlLgotIEZBQ0FERSBwcm92aWRlcyBubyBhdXRoZW50aWNhdGlvbiBtZWNoYW5pc20gd2hhdHNvZXZlci4KClRoZSByZWFzb25pbmcgaXMgdGhhdCBhbnlvbmUgYWJsZSB0byByZWFjaCB0aGUgc2VydmljZSBhbmQgcmVhZCB0aGUgcGFja2V0cyBwcm9iYWJseSBjYW4gc2VlIHRoZSBvdXRwdXQgb24gdGhlIHdhbGwgYW55d2F5LCBoZW5jZSB0aGVyZSBpcyBubyBmb2N1cyBvbiBzZWN1cml0eSBhdCB0aGlzIHRpbWUuICoqUGxlYXNlIG1ha2Ugc3VyZSB0byBzZXR1cCBwYWNrZXQgZmlsdGVycyBhbmQgbmV0d29yayB0dW5uZWxzIGJlZm9yZSBzZW5kaW5nIGFueSBzZW5zaXRpdmUgZGF0YSB0byBGQUNBREUhKioKCgoKCiMjIEN1c3RvbSBTaGFkZXJzCgpGQUNBREUgc3VwcG9ydHMgY3VzdG9tIHZlcnRleCBhbmQgZnJhZ21lbnQgc2hhZGVyczoKCjEuIE9uIHlvdXIgUmFzcGJlcnJ5IFBpLCBjcmVhdGUgYSBfLmZhY2FkZS9fIGRpcmVjdG9yeSBpbiB5b3VyIF8kSE9NRV86CgktIGBta2RpciAtcCB+Ly5mYWNhZGUvYAogICAgCjIuIERvd25sb2FkIHRoZSBkZWZhdWx0IHNoYWRlciBmcm9tIHRoZSBGQUNBREUgc291cmNlIHJlcG9zaXRvcnk6CgktIGBnaXQgYXJjaGl2ZSAtLXJlbW90ZT1odHRwczovL2dpdGh1Yi5jb20vRkVFREZBQ0UtQ09NL2ZhY2FkZSBIRUFEIHwgdGFyIHggLUMgfi8uZmFjYWRlLyBzaGFkZXIvZ3JpZC9kZWYudmVydGAKCjMuIENvcHkgdGhlIGRlZmF1bHQgc2hhZGVyIHRvIGEgbmV3IGZpbGUgX2Zvb2Jhci52ZXJ0XzoKCS0gYGNwIH4vLmZhY2FkZS9zaGFkZXIvZ3JpZC9kZWYudmVydCB+Ly5mYWNhZGUvc2hhZGVyL2dyaWQvZm9vYmFyLnZlcnRgCgo0LiBJbnN0cnVjdCBGQUNBREUgdG8gdXNlIHRoZSBuZXcgc2hhZGVyOgoJLSBgLi9mYWNhZGUgLWQgc2VydmUgbGluZXMgLXZlcnQgZm9vYmFyYAoKNS4gRWRpdCB0aGUgX34vLmZhY2FkZS9zaGFkZXIvZ3JpZC9mb29iYXIudmVydF8gZmlsZS4KCS0gVHJ5IGFkZGluZyBgcG9zLngqPXNpbihub3cpOyBwb3MueSo9Y29zKG5vdyk7YCBqdXN0IGJlZm9yZSB0aGUgbGluZSBzdGFydGluZyB3aXRoIGBnbFBvc2l0aW9uYCAJIAoKNi4gU2F2ZSB0aGUgX34vLmZhY2FkZS9zaGFkZXIvZ3JpZC9mb29iYXIudmVydF8gZmlsZS4gIAogICAtIFlvdSBzaG91bGQgc2VlIHRoZSBlZmZlY3Qgb2YgeW91ciBjaGFuZ2VzLiAKClNlZSA8aHR0cHM6Ly93d3cua2hyb25vcy5vcmcvb3BlbmdsZXMvc2RrL2RvY3MvcmVmZXJlbmNlX2NhcmRzL09wZW5HTC1FUy0yXzAtUmVmZXJlbmNlLWNhcmQucGRmPiBmb3Igc2hhZGVyIHN5bnRheCBhbmQgYXZhaWxhYmxlIGZ1bmN0aW9ucy4KCgoKIyMgRXhhbXBsZXMKCiMjIyBJbmZvcm1hdGl2ZSBVc2UKCkZBQ0FERSBjYW4gc2hvdyB5b3UgdGhlIGxpdmUgc3RhdHVzIG9mIHlvdXIgbWFjaGluZXMsIHNlcnZpY2VzIGFuZCBuZXR3b3JrcywgZWc6CgoKIyMjIyBgdG9wYCAtIHN5c3RlbSBzdGF0dXMKfn5+CmZhY2FkZSBzZXJ2ZSB0ZXJtIC12ZXJ0PXdhdmUgLW1hc2s9bWFzawpmYWNhZGUgZXhlYyAtaG9zdCByYXNwYmVycnlwaSB0ZXJtIC13PTgwIHRvcCAtMSAgICAgICAgICMgcnVuIG9uIHNlcnZlcgp+fn4KCgojIyMjIGB0Y3BkdW1wYCAtIGxpdmUgRE5TIHF1ZXJpZXMKfn5+CmZhY2FkZSBzZXJ2ZSBsaW5lcyAtdz04MCAtYnVmZmVyPTMyIC12ZXJ0PXJvbGwgLWRvd24Kc3VkbyB0Y3BkdW1wIC1pIHdsYW4wIC1uIC10IC1sIGRzdCBwb3J0IDUzIHwgZmNkICAgICAgICAjIHJ1biBvbiByb3V0ZXIKfn5+CgoKIyMjIyBgYWNjZXNzLmxvZ2AgLSBsaXZlIHdlYiByZXF1ZXN0cwp+fn4KZmFjYWRlIHNlcnZlIGxpbmVzIC13PTEyMCAtaD0xMiAtYnVmZmVyPTMyIC12ZXJ0PXZvcnRleCAtbWFzaz1tYXNrCnRhaWwgLWYgL3Zhci9sb2cvbmdpbngvYWNjZXNzLmxvZyB8IGZjZCAgICAgICAgICAgICAgICAgIyBydW4gb24gd2Vic2VydmVyCn5+fgoKCiMjIyMgYG10cmAgLSBjb250aW51b3VzIHRyYWNlIHJvdXRlCn5+fgpmYWNhZGUgLXEgc2VydmUgdGVybSAtdmVydD1kaXNrICYKZmFjYWRlIGV4ZWMgdGVybSAtdz02NCAtaD0xNiBzdWRvIG10ciAtbSAxMCAtLWRpc3BsYXltb2RlIDEgOC44LjguOAp+fn4KCgojIyMjIGBkYXRlYCAtIGN1cnJlbnQgZGF0ZSBhbmQgdGltZQp+fn4KZmFjYWRlIC1xIHNlcnZlIGxpbmVzIC12ZXJ0PXdhdmUgLWg9MiAtdz0xMCAtbWFzaz1tYXNrIC1kb3duIC1zbW9vdGg9ZiAtZm9udD1vY3JhZXh0IC16b29tPS44ICYKd2hpbGUgdHJ1ZTsgZG8gZGF0ZSAiKyVZLSVtLSVkIjsgc2xlZXAgMTsgZGF0ZSAiKyAlSDolTTolUyI7IHNsZWVwIDE7IGRvbmUgfCBmY2QKfn5+CgoKCiMjIyBDb2xsYWJvcmF0aXZlIFVzZQoKWW91IGNhbiB1c2UgRkFDQURFIHRvIGxvb2sgYXQgdGV4dCBvdXRwdXQgdG9nZXRoZXIsIGllIG9uZSBwZXJzb24gZGlyZWN0bHkgaW50ZXJhY3RzIHdpdGggYSBwcm9ncmFtIHdoaWxlIHRoZSBvdGhlciBwZW9wbGUgaW4gdGhlIHJvb20gY2FuIG9ic2VydmUgYW5kIGNvbW1lbnQ6CgoKIyMjIyBgYmFzaGAgLSBzaG93IHlvdXIgdGVhbSB3aGF0IGV4YWN0bHkgeW91IGFyZSBkb2luZyBpbiB5b3VyIHNoZWxsCn5+fgpmYWNhZGUgc2VydmUgdGVybQpmYWNhZGUgZXhlYyAtaG9zdCByYXNwYmVycnlwaSB0ZXJtIC13PTgwIC1oPTI1IGJhc2ggICAgICMgcnVuIG9uIHdvcmtzdGF0aW9uCn5+fgoKCiMjIyMgYGZyb3R6YCAtIHBsYXkgdGV4dCBhZHZlbnR1cmVzIGFzIGEgZ3JvdXAKfn5+CmZhY2FkZSBzZXJ2ZSB0ZXJtIC1mb250PXNwYWNlbW9ubwpmYWNhZGUgZXhlYyB0ZXJtIC13PTY0IC1oPTE2IGZyb3R6IC9wYXRoL3RvL2hpdGNoaGlrZXJzX2d1aWRlLno1ICMgcnVuIG9uIHdvcmtzdGF0aW9uCn5+fgoKCgojIyMgRGVjb3JhdGl2ZSBVc2UKCkZBQ0FERSB3b3JrcyB2ZXJ5IHdlbGwgaWYgeW91IGp1c3Qgd2FudCB0byBoYXZlIHNvbWUgc3R5bGlzaCB0ZXh0IHNjcm9sbGluZyBhY3Jvc3MgeW91ciB3YWxsOgoKCiMjIyMgYG1hbmAgLSBzb21lIG1hbnBhZ2VzIGFyZSBxdWl0ZSBwcmV0dHkgOikKfn5+CmZhY2FkZSAtcSBzZXJ2ZSBsaW5lcyAtdz01MCAtdmVydD1jcmF3bCAmCk1BTldJRFRIPTUwIE1BTlBBR0VSPWNhdCBtYW4gZ2l0LXJlYmFzZSBcCnwgd2hpbGUgcmVhZCAtciBsaW5lOyBkbyBlY2hvICIkbGluZSIgfCBmY2Q7IHNsZWVwIC45OyBkb25lCn5+fgoKCiMjIyMgYFJGQ3NgIC0gaW50ZXJuZXR3b3JraW5nIHNwZWNpZmljYXRpb25zIGluIHBsYWluIHRleHQgZm9ybWF0Cn5+fgpmYWNhZGUgLXEgc2VydmUgbGluZXMgLXc9NzIgLXZlcnQ9cm93cyAmCmN1cmwgLUwgaHR0cHM6Ly90b29scy5pZXRmLm9yZy9yZmMvcmZjMjQ2MC50eHQgXAp8IHdoaWxlIHJlYWQgLXIgbGluZTsgZG8gZWNobyAiJGxpbmUiIHwgZmNkOyBzbGVlcCAuOTsgZG9uZQp+fn4KCgojIyMjIGBQSFJBQ0tgIC0geW91ciBmYXZvdXJpdGUgaGFja2luZyB6aW5lIGFydGljbGVzCn5+fgpmYWNhZGUgLXEgc2VydmUgbGluZXMgLXc9ODAgLXZlcnQ9cm9sbCAmCmN1cmwgLUwgaHR0cDovL3BocmFjay5vcmcvYXJjaGl2ZXMvdGd6L3BocmFjazQ5LnRhci5neiBcCnwgdGFyIHhmeiAvZGV2L3N0ZGluIC4vMTQudHh0IC0tdG8tc3Rkb3V0IFwKfCB3aGlsZSByZWFkIC1yIGxpbmU7IGRvIGVjaG8gIiRsaW5lIiB8IGZjZDsgc2xlZXAgLjk7IGRvbmUKfn5+CgoKIyMjIyBgcHIwbmAgLSBudWRlcyBvbGRlciB0aGFuIHRoZSA8SU1HPiB0YWcKfn5+CmZhY2FkZSAtcSBzZXJ2ZSBsaW5lcyAtdz04MCAtdmVydD13YXZlIC1tYXNrPW1hc2sgJgpjdXJsIC1MIGh0dHBzOi8vd3d3LmFzY2lpcHIwbi5jb20vcHIwbi9odW5reS9odW5reTAwLnR4dCBcCnwgd2hpbGUgcmVhZCAtciBsaW5lOyBkbyBlY2hvICIkbGluZSIgfCBmY2Q7IHNsZWVwIC45OyBkb25lCgp+fn4KCgojIyMjIGAubmZvYCAtIGRlbW8gc2NlbmUgcmVsZWFzZSBub3RlcyB3aXRoIDEzMzcgYXNjaWkgYXJ0Cn5+fgpmYWNhZGUgLXEgc2VydmUgbGluZXMgLXc9ODAgLXZlcnQ9d2F2ZSAtbWFzaz1tYXNrIC1mb250IGFkb3JlNjQgJgpjdXJsIC1MIGh0dHBzOi8vY29udGVudC5wb3VldC5uZXQvZmlsZXMvbmZvcy8wMDAxMi8wMDAxMjAzMS50eHQgXAp8IHdoaWxlIHJlYWQgLXIgbGluZTsgZG8gZWNobyAiJGxpbmUiIHwgZmNkOyBzbGVlcCAuOTsgZG9uZQp+fn4KCgotLS0tCgpJZiB5b3UgZW5qb3kgRkFDQURFLCB0ZWxsIHVzIGhvdyB5b3UgYXJlIHVzaW5nIGl0IGF0IDxmYWNhZGVAZmVlZGZhY2UuY29tPiEKCg==
`,
}
