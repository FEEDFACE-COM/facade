
package main

import (
    "fmt"
    "net"    
    "time"
    "io"
//    "os"
    "bufio"
    "encoding/json"
//    "encoding/gob"
    log "./log"
    facade "./facade"
)


const DEBUG_ACCEPT = false
const DEBUG_RECV =   false


type Server   struct {
    host string
    confPort uint
    textPort uint
    timeout float64
}


func NewServer(host string, confPort uint, textPort uint, timeout float64) (*Server) { 
    return &Server{host:host, confPort: confPort, textPort: textPort, timeout: timeout} 
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
        if DEBUG_ACCEPT {log.Debug("accept conf from %s",confConn.RemoteAddr().String()) }
        go server.ReceiveConf(confConn, confChan)

    }
}

func (server *Server) ListenText(bufChan chan facade.BufferItem) { 
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
        if DEBUG_ACCEPT { log.Debug("accept text from %s",textConn.RemoteAddr().String()) }
        if server.timeout == 0.0 {
            textConn.SetReadDeadline( time.Time{} )
        } else {
            textConn.SetReadDeadline(time.Now().Add( 1 * time.Second ) )
        }
        go server.ReceiveText(textConn, bufChan)

    }
}


func (server *Server) ReceiveConf(confConn net.Conn, confChan chan facade.Config) {
    defer func() { 
        if DEBUG_ACCEPT { log.Debug("close conf %s",confConn.RemoteAddr().String()) }
        confConn.Close() 
    }()
    decoder := json.NewDecoder(confConn)
    config := make(facade.Config)
    err := decoder.Decode(&config)
    if err != nil {
        log.Error("fail to decode %s: %s",confConn.RemoteAddr().String(),err)
        return
    }
    if DEBUG_RECV { log.Debug("receive conf %s",config.Desc()) }
    confChan <- config
}

func (server *Server) ReceiveText(textConn net.Conn, bufChan chan facade.BufferItem) {
    defer func() { 
        if DEBUG_ACCEPT { log.Debug("close text %s",textConn.RemoteAddr().String()); }
        textConn.Close() 
    }()
    const BUFFER_SIZE = 1024
	var buf []byte = make([]byte, BUFFER_SIZE)
	var rem []byte = []byte{}
	var tmp []byte
	reader := bufio.NewReader( textConn )
	for {
        n,err := reader.Read(buf)
		if err == io.EOF { break }
		if err != nil {
			log.Error("read %s error: %s",textConn.RemoteAddr().String(),err)
			break
		}
        if DEBUG_RECV { log.Debug("recv %d byte:\n%s",n,log.Dump(buf,n,0)) }
		tmp = append(rem, buf[:n] ... )
		rem, err = facade.ProcessRaw(tmp, bufChan)
		if err != nil {
            log.Error("process error: %s",err)    		
        }
    }
    
//    scanner := bufio.NewScanner(textConn)
//    for scanner.Scan() {
//        if server.timeout == 0.0 {
//            textConn.SetReadDeadline( time.Time{} )
//        } else {
//            textConn.SetReadDeadline(time.Now().Add( 1 * time.Second ) )
//        }
//        
////        textConn.SetReadDeadline(time.Now().Add( 5 * time.Second ) )
//        text := scanner.Text()
//        if DEBUG_RECV {
//            log.Debug("receive text %s",text)
//        }
//        textChan <- facade.RawText(text)
//    }
//    err := scanner.Err()
//    if err != nil {
//        log.Error("fail to scan %s: %s",textConn.RemoteAddr().String(),err)
//    }
}

