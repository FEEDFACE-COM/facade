
package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
//    "bytes"
    "encoding/gob"
//    "time"
    log "./log"
    facade "./facade"
)

const DEBUG_SEND = true

type Client   struct {
    host string
    confPort uint
    textPort uint
    timeout float64
    
}
func NewClient(host string, confPort uint, textPort uint, timeout float64) (*Client) { 
    return &Client{host:host, confPort:confPort, textPort: textPort, timeout:timeout}
}


func (client *Client) ScanAndSendText() {
    textConnStr := fmt.Sprintf("%s:%d",client.host,client.textPort)
    log.Info("connect %s",textConnStr) 
    conn, err := net.Dial("tcp", textConnStr)
    if err != nil {
        log.Error("fail to dial %s: %s",textConnStr,err)
        return
    }
    defer func() { /*log.Debug("close %s",conn.RemoteAddr().String());*/ conn.Close() }()

    var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        _, err = conn.Write( []byte(text+"\n") )
        if err != nil {
            log.Error("fail to write to %s: %s",textConnStr,err)
            return    
        }
        if DEBUG_SEND {
            log.Debug("send text: %s",text)
        }
    }
    err = scanner.Err()
    if err != nil {
        log.Error("fail to scan: %s",err)
    }
    
}


func (client *Client) SendConf(config *facade.Config) { 
    confConnStr := fmt.Sprintf("%s:%d",client.host,client.confPort)
    log.Info("config %s",config.Desc())
    log.Info("connect to %s",confConnStr) 
    conn, err := net.Dial("tcp", confConnStr)
    if err != nil {
        log.Error("fail to dial %s: %s",confConnStr,err)
        return
    }
    defer func() { /*log.Debug("close %s",conn.RemoteAddr().String());*/ conn.Close() }()
    enc := gob.NewEncoder(conn)
    err = enc.Encode( *config )
    if err != nil {
        log.Error("fail to encode %s: %s",config.Desc(),err)
        return
    }
    if DEBUG_SEND {
        log.Debug("send conf: %s",config.Desc())
    }
}

    
