
package main

import (
    "fmt"
    "net"    
    "time"
    "bufio"
    "encoding/gob"
    log "./log"
    facade "./facade"
)


const DEBUG_ACCEPT = false
const DEBUG_RECV =   false


type Server   struct {
    host string
    confPort uint
    textPort uint
}


func NewServer(host string, confPort uint, textPort uint) (*Server) { 
    return &Server{host:host, confPort: confPort, textPort: textPort} 
}

func (server *Server) ListenConf(confChan chan facade.Config) {
    confListenStr := fmt.Sprintf("%s:%d",server.host,server.confPort)
    log.Debug("listen for config on %s",confListenStr) 
    confListener, err := net.Listen("tcp",confListenStr)
    if err != nil {
        log.PANIC("fail listen on %s: %s",confListenStr,err)
    }
    defer func() { log.Debug("stop listen conf on %s",confListener.Addr().String()); confListener.Close() }()
    log.Info("listening for conf on %s",confListener.Addr().String()) 
    for {
        confConn, err := confListener.Accept()
        if err != nil {
            log.Error("fail accept on %s: %s",confListenStr,err)    
            continue
        }
        if DEBUG_ACCEPT {
            log.Debug("accept conf from %s",confConn.RemoteAddr().String())    
        }
        go server.ReceiveConf(confConn, confChan)

    }
}

func (server *Server) ListenText(textChan chan facade.RawText) { 
    textListenStr := fmt.Sprintf("%s:%d",server.host,server.textPort)
    log.Debug("listen for text on %s",textListenStr) 
    textListener, err := net.Listen("tcp",textListenStr)
    if err != nil {
        log.PANIC("fail listen on %s: %s",textListenStr,err)
    }
    defer func() { /*log.Debug("stop listen text on %s",textListener.Addr().String());*/ textListener.Close() }()
    log.Info("listening for text on %s",textListener.Addr().String()) 
    
    for {
        textConn, err := textListener.Accept()
        if err != nil {
            log.Error("fail accept on %s: %s",textListenStr,err)    
            continue
        }
        if DEBUG_ACCEPT {
            log.Debug("accept text from %s",textConn.RemoteAddr().String())    
        }
        go server.ReceiveText(textConn, textChan)

    }
}


func (server *Server) ReceiveConf(confConn net.Conn, confChan chan facade.Config) {
    defer func() { 
        if DEBUG_ACCEPT {
            log.Debug("close conf %s",confConn.RemoteAddr().String());
        }
            confConn.Close() 
    }()
    decoder := gob.NewDecoder(confConn)
    config := &facade.Config{}
    confConn.SetReadDeadline(time.Now().Add( 5 * time.Second ) )
    err := decoder.Decode(config)
    if err != nil {
        log.Error("fail to decode %s: %s",confConn.RemoteAddr().String(),err)
        return
    }
    if DEBUG_RECV {
        log.Debug("receive conf %s",config.Desc())
    }
    confChan <- *config
}

func (server *Server) ReceiveText(textConn net.Conn, textChan chan facade.RawText) {
    defer func() { 
        if DEBUG_ACCEPT {
            log.Debug("close text %s",textConn.RemoteAddr().String());
        }
        textConn.Close() 
    }()
    scanner := bufio.NewScanner(textConn)
    for scanner.Scan() {
        textConn.SetReadDeadline(time.Now().Add( 5 * time.Second ) )
        text := scanner.Text()
        if DEBUG_RECV {
            log.Debug("receive text %s",text)
        }
        textChan <- facade.RawText(text)
    }
    err := scanner.Err()
    if err != nil {
        log.Error("fail to scan %s: %s",textConn.RemoteAddr().String(),err)
    }
}

