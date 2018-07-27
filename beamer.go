
package main

import (
    "fmt"
    "net"    
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    proto "./proto"
)

type FcdBeamer   struct {
    host string
    port uint
}

func NewFcdBeamer(host string, port uint) (*FcdBeamer) { 
    return &FcdBeamer{host:host, port: port} 
}

func (beamer *FcdBeamer) beam() { 
    Debug("facade beamer listening at %s:%d",beamer.host,beamer.port) 
    listen, err := net.Listen("tcp",fmt.Sprintf("%s:%d",beamer.host,beamer.port))
    //beamer.host+":"+string(beamer.port) )
    if err != nil {
        FATAL("fail to listen on %s:%d: %s",beamer.host,beamer.port,err)
    }
    server := grpc.NewServer()
    proto.RegisterFacadeServer(server, &FcdServer{} )
    err = server.Serve(listen)
    if err != nil {
        FATAL("fail to serve on %s:%d: %s",beamer.host,beamer.port,err)
    }
}

type FcdServer struct {}

func (server *FcdServer) BeamText(ctx context.Context, request *proto.BeamRequest) (*proto.Response, error) {
    Debug(">> %s",request.Text)
    return &proto.Response{Success: true}, nil
}




