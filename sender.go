
package main

import (
    "fmt"
    "time"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/status"
    proto "./proto"
)

type FcdSender   struct {
    host string
    port uint
    timeout float64
    text string
    
}
func NewFcdSender(host string, port uint, timeout float64) (*FcdSender) { 
    return &FcdSender{host:host, port:port, timeout:timeout}
}

func (sender *FcdSender) send() { 
    Info("facade sender connecting to %s:%d",sender.host,sender.port) 
    opts := grpc.WithInsecure()
    
    connect, err := grpc.Dial(fmt.Sprintf("%s:%d",sender.host,sender.port), opts)
    //sender.host+":"+string(sender.port), opts)
    if err != nil {
        FATAL("fail dial to %s:%d because %s",sender.host,sender.port,err)
    }
    defer connect.Close()
    
    client := proto.NewFacadeClient(connect)
    
    if sender.text == "" {
        Info("piping from stdin to %s:%d",sender.host,sender.port)
    } else {
        err := sender.runBeamText(client)
        if err != nil {
            Error("fail to send to %s:%d because %s",sender.host,sender.port,err)   
        }
    }
}

func (sender *FcdSender) runBeamText(client proto.FacadeClient) error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(sender.timeout)*time.Second)
    defer cancel()
    response, err := client.BeamText(ctx, &proto.BeamRequest{Text: sender.text})
    if err != nil {
        status,_ := status.FromError(err)
        FATAL("fail to beam to %s:%d because %s",sender.host,sender.port,status.Message())
    }
    if ! response.Success {
        Error("received unsuccessful response")
        return fmt.Errorf("received unsuccessful response")
    }
    Debug(">> %s",sender.text)
    return nil
}

