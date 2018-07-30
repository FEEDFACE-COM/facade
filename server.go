
// +build !linux,!arm
// +build !darwin,!amd64


package main



const SERVER_MODE_AVAILABLE = false

type Server   struct {}
func NewServer(host string, port uint) (*Server) { return &Server{} }
func (beamer *Server) Serve() { FATAL("server mode not available for %s",BUILD_PLATFORM) }

