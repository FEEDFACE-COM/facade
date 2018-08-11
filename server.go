
package main

import (
    "fmt"
    "net"    
    "time"
    "bufio"
    log "./log"
    render "./render"
)


const DEBUG_ACCEPT = true
const DEBUG_RECV =   true


type Server   struct {
    host string
    confPort uint
    textPort uint
}

func NewServer(host string, confPort uint, textPort uint) (*Server) { 
    return &Server{host:host, confPort: confPort, textPort: textPort} 
}

func (server *Server) Listen(textChan chan render.Text) { 
    confListenStr := fmt.Sprintf("%s:%d",server.host,server.confPort)
    log.Debug("listen for config on %s",confListenStr) 
    confListener, err := net.Listen("tcp",confListenStr)
    if err != nil {
        log.Fatal("fail to listen on %s: %s",confListenStr,err)
    }
    defer func() { log.Debug("stop listen for conf on %s",confListener.Addr().String()); confListener.Close() }()
    
    
    textListenStr := fmt.Sprintf("%s:%d",server.host,server.textPort)
    log.Debug("listen for text on %s",textListenStr) 
    textListener, err := net.Listen("tcp",textListenStr)
    if err != nil {
        log.Fatal("fail to listen on %s: %s",textListenStr,err)
    }
    defer func() { log.Debug("stop listen for text on %s",textListener.Addr().String()); textListener.Close() }()

    log.Info("listening for text on %s listening for conf on %s",textListener.Addr().String(),confListener.Addr().String()) 
    
    for {
        textConn, err := textListener.Accept()
        if err != nil {
            log.Error("fail to accept on %s: %s",textListenStr,err)    
            continue
        }
        if DEBUG_ACCEPT {
            log.Debug("accept %s",textConn.RemoteAddr().String())    
        }
        go server.ReceiveText(textConn, textChan)

    }
}

func (server *Server) ReceiveText(textConn net.Conn, textChan chan render.Text) {
    defer func() { log.Debug("close %s",textConn.RemoteAddr().String()); textConn.Close() }()
    scanner := bufio.NewScanner(textConn)
    for scanner.Scan() {
        textConn.SetReadDeadline(time.Now().Add( 5 * time.Second ) )
        text := scanner.Text()
        if DEBUG_RECV {
            log.Debug(">>>> %s",text)
        }
        textChan <- render.Text(text)
    }
    err := scanner.Err()
    if err != nil {
        log.Error("fail to scan %s: %s",textConn.RemoteAddr().String(),err)
    }
}

