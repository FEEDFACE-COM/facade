
package main

import (
    "fmt"
    "net"    
    log "./log"
//    render "./render"
//    "golang.org/x/net/context"
//    "google.golang.org/grpc"
//    proto "./proto"
)


const SERVER_MODE_AVAILABLE = true


type Server   struct {
    host string
    confPort uint
    textPort uint
}

func NewServer(host string, confPort uint, textPort uint) (*Server) { 
    return &Server{host:host, confPort: confPort, textPort: textPort} 
}

func (server *Server) Serve() { 
    
    log.Debug("listen for config on %s:%d",server.host,server.confPort) 
    confListener, err := net.Listen("tcp",fmt.Sprintf("%s:%d",server.host,server.confPort))
    if err != nil {
        log.Fatal("fail to listen on %s:%d: %s",server.host,server.confPort,err)
    }
    log.Debug("listen for text on %s:%d",server.host,server.textPort) 
    textListener, err := net.Listen("tcp",fmt.Sprintf("%s:%d",server.host,server.textPort))
    if err != nil {
        log.Fatal("fail to listen on %s:%d: %s",server.host,server.textPort,err)
    }



    log.Info("listening for text on %s listening for conf on %s",textListener.Addr().String(),confListener.Addr().String()) 



    
    
//    facadeServer := grpc.NewServer()
//    proto.RegisterFacadeServer(facadeServer, &Server{} )
//    err = facadeServer.Serve(listen)
//    if err != nil {
//        FATAL("fail to serve on %s:%d: %s",server.host,server.port,err)
//    }
}


//func (server *Server) BeamText(ctx context.Context, request *proto.BeamRequest) (*proto.Response, error) {
//    Debug(">> %s",request.Text)
//    return &proto.Response{Success: true}, nil
//}


