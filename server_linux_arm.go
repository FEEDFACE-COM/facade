
package main

import (
    "fmt"
    "net"    
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    proto "./proto"
)


const SERVER_MODE_AVAILABLE = true


type Server   struct {
    host string
    port uint
}

func NewServer(host string, port uint) (*Server) { 
    return &Server{host:host, port: port} 
}

func (server *Server) Serve() { 
    Debug("facade server listening at %s:%d",server.host,server.port) 
    listen, err := net.Listen("tcp",fmt.Sprintf("%s:%d",server.host,server.port))
    if err != nil {
        FATAL("fail to listen on %s:%d: %s",server.host,server.port,err)
    }
    facadeServer := grpc.NewServer()
    proto.RegisterFacadeServer(facadeServer, &Server{} )
    err = facadeServer.Serve(listen)
    if err != nil {
        FATAL("fail to serve on %s:%d: %s",server.host,server.port,err)
    }
}


func (server *Server) BeamText(ctx context.Context, request *proto.BeamRequest) (*proto.Response, error) {
    Debug(">> %s",request.Text)
    return &proto.Response{Success: true}, nil
}


