
package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
//    "bytes"
    "encoding/json"
//    "time"
    log "./log"
    facade "./facade"
)

const DEBUG_CLIENT = true
const DEBUG_CLIENT_DUMP = true

type Client   struct {
    host string
    confPort uint
    textPort uint
    textConn net.Conn
    textConnStr string
    timeout float64
    
}
func NewClient(host string, confPort uint, textPort uint, timeout float64) (*Client) { 
    return &Client{host:host, confPort:confPort, textPort: textPort, timeout:timeout}
}

func (client *Client) OpenText() error {
    var err error
    client.textConnStr = fmt.Sprintf("%s:%d",client.host,client.textPort)
    if DEBUG_CLIENT {
        log.Debug("dial %s",client.textConnStr) 
    }
    client.textConn, err = net.Dial("tcp", client.textConnStr)
    if err != nil {
        if DEBUG_CLIENT {
            log.Error("fail to dial %s: %s",client.textConnStr,err)
        }
        return log.NewError("fail to dial %s",client.textConnStr) 
    }
    return nil
}

func (client *Client) CloseText() {
    log.Debug("close %s",client.textConn.RemoteAddr().String());
    client.textConn.Close()    
}


func (client *Client) ScanAndSendText() error {
    var err error

    var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        err = client.SendText( []byte(text+"\n") )
        if err != nil {
            return err
        }
    }
    err = scanner.Err()
    if err != nil {
        log.Error("fail to scan: %s",err)
    }
    return nil
}

func (client *Client) SendText(text []byte) error {
    var err error
    _, err = client.textConn.Write(text) 
    if err != nil {
        if DEBUG_CLIENT { log.Error("fail to write to %s: %s",client.textConnStr,err) }
        return log.NewError("fail to write to %s",client.textConnStr)
    }
    if DEBUG_CLIENT_DUMP { log.Debug("sent %d byte text:\n%s",len(text),log.Dump(text,0,0)) 
    } else if DEBUG_CLIENT { log.Debug("sent %d byte text",len(text)) }
    return nil
}

func (client *Client) SendConf(config *facade.Config) error { 
    confConnStr := fmt.Sprintf("%s:%d",client.host,client.confPort)
    if DEBUG_CLIENT { log.Debug("dial %s",confConnStr) }
    conn, err := net.Dial("tcp", confConnStr)
    if err != nil {
        if DEBUG_CLIENT { log.Error("fail to dial %s: %s",confConnStr,err) }
        return log.NewError("fail to dial %s",confConnStr)
    }
    defer func() { 
        if DEBUG_CLIENT { log.Debug("close %s",conn.RemoteAddr().String()); }
        conn.Close()
    }()
    encoder := json.NewEncoder(conn)
    err = encoder.Encode( *config )
    if err != nil {
        if DEBUG_CLIENT { log.Error("fail to encode %s: %s",config.Desc(),err) }
        return log.NewError("fail to encode %s",config.Desc())
    }
    if DEBUG_CLIENT { log.Debug("sent conf %s",config.Desc()) }
    return nil
}

    
