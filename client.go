
package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
//    "time"
    log "./log"
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

func (client *Client) SendText(text string) { 
    textConnStr := fmt.Sprintf("%s:%d",client.host,client.textPort)
    log.Info("connect text to %s",textConnStr) 
    conn, err := net.Dial("tcp", textConnStr)
    if err != nil {
        log.Error("fail to dial %s: %s",textConnStr,err)
        return
    }
    _, err = conn.Write( []byte(text) )
    if err != nil {
        log.Error("fail to write to %s: %s",textConnStr,err)
    }
    if DEBUG_SEND {
        log.Debug("> %s",text)
    }
    conn.Close()

}

func (client *Client) ScanAndSendText() {
    textConnStr := fmt.Sprintf("%s:%d",client.host,client.textPort)
    log.Info("connect text to %s",textConnStr) 
    conn, err := net.Dial("tcp", textConnStr)
    if err != nil {
        log.Error("fail to dial %s: %s",textConnStr,err)
        return
    }
    defer func() { log.Debug("close %s",conn.RemoteAddr().String()); conn.Close() }()

    var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        _, err = conn.Write( []byte(text+"\n") )
        if err != nil {
            log.Error("fail to write to %s: %s",textConnStr,err)
            return    
        }
        if DEBUG_SEND {
            log.Debug(">>>> %s",text)
        }
    }
    err = scanner.Err()
    if err != nil {
        log.Error("fail to scan: %s",err)
    }
    
}


func (client *Client) SendConf(conf string) { 
    confConnStr := fmt.Sprintf("%s:%d",client.host,client.confPort)
    log.Info("connect conf to %s",confConnStr) 
    conn, err := net.Dial("tcp", confConnStr)
    if err != nil {
        log.Error("fail to dial %s: %s",confConnStr,err)
        return
    }
    _, err = conn.Write( []byte(conf) )
    if err != nil {
        log.Error("fail to write to %s: %s",confConnStr,err)
    }
    if DEBUG_SEND {
        log.Debug("==== %s",conf)
    }
    conn.Close()

}

    
