
package main

import (
//    "fmt"
//    "time"
//    "golang.org/x/net/context"
//    "google.golang.org/grpc"
//    "google.golang.org/grpc/status"
    log "./log"
//    proto "./proto"
    
)

type Client   struct {
    host string
    confPort uint
    textPort uint
    timeout float64
    text string
    
}
func NewClient(host string, confPort uint, textPort uint, timeout float64) (*Client) { 
    return &Client{host:host, confPort:confPort, textPort: textPort, timeout:timeout}
}

func (client *Client) SendText() { 
    log.Info("connect to %s:%d",client.host,client.textPort) 

}

func (client *Client) SendConf() { 
    log.Info("connect to %s:%d",client.host,confPort) 

}

    
//    connect, err := grpc.Dial(fmt.Sprintf("%s:%d",client.host,client.port), opts)
//    if err != nil {
//        log.Fatal("fail dial to %s:%d because %s",client.host,client.port,err)
//    }
//    defer connect.Close()
//    
//    facadeClient := proto.NewFacadeClient(connect)
//    
//    if client.text == "" {
//        log.Info("piping from stdin to %s:%d",client.host,client.port)
//    } else {
//        err := client.BeamText(facadeClient)
//        if err != nil {
//            log.Error("fail to send to %s:%d because %s",client.host,client.port,err)   
//        }
//    }
//}
//
//func (client *Client) BeamText(facadeClient proto.FacadeClient) error {
//    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(client.timeout)*time.Second)
//    defer cancel()
//    response, err := facadeClient.BeamText(ctx, &proto.BeamRequest{Text: client.text})
//    if err != nil {
//        status,_ := status.FromError(err)
//        log.Fatal("fail to beam to %s:%d because %s",client.host,client.port,status.Message())
//    }
//    if ! response.Success {
//        return fmt.Errorf("received unsuccessful response")
//    }
//    log.Debug(">> %s",client.text)
//    return nil
//}
//
