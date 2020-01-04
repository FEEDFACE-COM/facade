package main

import (
	"fmt"
	//    "net"
	"bufio"
	"os"
	"time"

	//    "bytes"
	//    "encoding/json"
	facade "./facade"
	log "./log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	grpcstatus "google.golang.org/grpc/status"
)

const DEBUG_CLIENT = true
const DEBUG_CLIENT_DUMP = false

type Client struct {
	host    string
	port    uint
	connStr string
	timeout time.Duration

	connection *grpc.ClientConn
	client     facade.FacadeClient
	stream     facade.Facade_DisplayClient
	cancel     context.CancelFunc
}

func NewClient(host string, port uint, timeout float64) *Client {
	ret := &Client{host: host, port: port}
	ret.connStr = fmt.Sprintf("%s:%d", ret.host, ret.port)
	ret.timeout = time.Duration(1000.*timeout) * time.Millisecond
	return ret
}

func (client *Client) Close() {
	if client.cancel != nil {
		client.cancel()
		client.cancel = nil
		client.stream = nil
	}
	if client.connection != nil {
		client.connection.Close()
		client.connection = nil
	}
}

//
// reading bytes not scanning lines:
//
//func (client *Client) ScanAndSendText() error {
//    var ret error = nil
//
//
//    const BUFFER_SIZE = 4096
//	var buf []byte = make([]byte, BUFFER_SIZE)
//
//    reader := bufio.NewReader(os.Stdin)
//
//
//	for {
//        n,err := reader.Read(buf)
//		if err == io.EOF {
//    		if DEBUG_CLIENT { log.Debug("read stdin eof") }
//    		break
//        }
//		if err != nil {
//			return log.NewError("fail read stdin: %s",err)
//		}
//        if DEBUG_CLIENT_DUMP { log.Debug("read %d byte:\n%s",n,log.Dump(buf,n,0))
//        } else if DEBUG_CLIENT { log.Debug("read %d byte",n) }
//
//        err = client.SendText( buf[:n] )
//        if err != nil {
//            log.Error("fail to send: %s",err)
//            ret = log.NewError("fail to send: %s",err)
//        }
//    }
//
//    return ret
//
//}

func (client *Client) ScanAndSendText() error {
	var err error

	var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		err = client.SendText([]byte(text + "\n"))
		if err != nil {
			return log.NewError("fail to send: %s", err)
		}
	}
	err = scanner.Err()
	if err != nil {
		log.Error("fail to scan: %s", err)
	}

	return nil
}

func (client *Client) OpenTextStream() error {
	var err error
	var ctx context.Context

	if client.stream != nil || client.cancel != nil {
		return log.NewError("fail to open stream: existing stream/cancel")
	}

	ctx, client.cancel = context.WithCancel(context.Background())
	client.stream, err = client.client.Display(ctx)
	if err != nil {
		return log.NewError("fail to get display stream: %s", err)
	}
	return nil
}

func (client *Client) CloseTextStream() error {

	if client.stream == nil {
		return log.NewError("fail to close stream: no stream")
	}

	status, err := client.stream.CloseAndRecv()
	client.stream = nil
	client.cancel = nil

	if err != nil {
		stat, _ := grpcstatus.FromError(err)
		return log.NewError("fail to close: %s", stat.Message())
	} else if !status.GetSuccess() {
		return log.NewError("display error: %s", status.GetError())
	}

	return nil
}

func (client *Client) SendText(raw []byte) error {

	if client.stream == nil {
		return log.NewError("no stream")
	}

	rawText := facade.RawText{Raw: raw}
	ret := client.stream.Send(&rawText)
	if DEBUG_CLIENT_DUMP {
		log.Debug("sent %d byte text:\n%s", len(raw), log.Dump(raw, 0, 0))
	} else if DEBUG_CLIENT {
		log.Debug("sent %d byte text", len(raw))
	}
	return ret
}

func (client *Client) SendConf(config *facade.Config) error {

	if client.connection == nil {
		return log.NewError("no connection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	status, err := client.client.Configure(ctx, config)
	if err != nil {
		stat, _ := grpcstatus.FromError(err)
		return log.NewError("fail to send: %s", stat.Message())
	} else if !status.GetSuccess() {
		return log.NewError("conf error: %s", status.GetError())
	}
	if DEBUG_CLIENT {
		log.Debug("sent to %s %s", client.connStr, config.Desc())
	}
	return nil
}

func (client *Client) QueryInfo() (string, error) {
	if client.connection == nil {
		return "", log.NewError("no connection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	status, err := client.client.Query(ctx, &facade.Empty{})
	if err != nil {
		stat, _ := grpcstatus.FromError(err)
		return "", log.NewError("fail to send: %s", stat.Message())
	} else if !status.GetSuccess() {
		return "", log.NewError("query error: %s", status.GetError())
	}

	info := status.GetInfo()
	if info == "" {
		return "", log.NewError("empty info")
	}

	return info, nil

}

func (client *Client) Dial() error {
	var err error

	opts := grpc.WithInsecure()

	if DEBUG_CLIENT {
		log.Debug("dial %s timeout %.1fs", client.connStr, client.timeout.Seconds())
	}
	client.connection, err = grpc.Dial(client.connStr, opts)
	if err != nil {
		return log.NewError("fail to dial %s: %s", client.connStr, err)
	}
	client.client = facade.NewFacadeClient(client.connection)
	return nil
}
