package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/log"
	"bufio"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	grpcstatus "google.golang.org/grpc/status"
	"net"
	"os"
	"time"
)

const DEBUG_CLIENT = false
const DEBUG_CLIENT_DUMP = false

type Client struct {
	connStr string
	timeout time.Duration

	connection *grpc.ClientConn
	client     facade.FacadeClient
	stream     facade.Facade_PipeClient
	cancel     context.CancelFunc
}

func NewClient(host string, port uint, timeout float64, inet bool, inet6 bool) *Client {

	ret := &Client{connStr: ""}

	var address string = host

	if !inet && !inet6 {
		log.PANIC("all network families disabled")
	}

	if inet || inet6 {

		// force network protocol by resolving hostname explicitly

		ip := net.ParseIP(host)
		if ip.To4() != nil && inet {
			if DEBUG_CLIENT {
				log.Debug("%s use given ipv4 address %s", ret.Desc(), ip.String())
			}
			address = ip.String()
		} else if ip.To16() != nil && inet6 {
			if DEBUG_CLIENT {
				log.Debug("%s use given ipv6 address %s", ret.Desc(), ip.String())
			}
			address = ip.String()
		} else {
			if DEBUG_CLIENT {
				log.Debug("%s lookup address for %s", ret.Desc(), host)
			}
			names, err := net.LookupHost(host)
			if err != nil {
				log.PANIC("fail to lookup address for %s: %s", host, err)
			}
			for _, name := range names {
				ip := net.ParseIP(name)
				if ip.To16() != nil && inet6 {
					if DEBUG_CLIENT {
						log.Debug("%s use resolved ipv6 address %s", ret.Desc(), ip.String())
					}
					address = ip.String()
					break
				}
				if ip.To4() != nil && inet {
					if DEBUG_CLIENT {
						log.Debug("%s use resolved ipv4 address %s", ret.Desc(), ip.String())
					}
					address = ip.String()
					break
				}
			}
			if address == host {
				log.PANIC("%s fail to find address for %s", ret.Desc(), host)
			}
		}

	}

	ret.connStr = net.JoinHostPort(address, fmt.Sprintf("%d", port))
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
		log.Error("%s fail to scan: %s", client.Desc(), err)
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
	client.stream, err = client.client.Pipe(ctx)
	if err != nil {
		return log.NewError("fail to get pipe stream: %s", err)
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
		return log.NewError("pipe error: %s", status.GetError())
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
		log.Debug("%s sent %d byte text:\n%s", client.Desc(), len(raw), log.Dump(raw, 0, 0))
	} else if DEBUG_CLIENT {
		log.Debug("%s sent %d byte text", client.Desc(), len(raw))
	}
	return ret
}

func (client *Client) SendConf(config *facade.Config) error {

	if client.connection == nil {
		return log.NewError("no connection")
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	log.Info("%s send config %s", client.Desc(), config.Desc())
	status, err := client.client.Conf(ctx, config)
	if err != nil {
		stat, _ := grpcstatus.FromError(err)
		return log.NewError("fail to send: %s", stat.Message())
	} else if !status.GetSuccess() {
		return log.NewError("conf error: %s", status.GetError())
	}
	return nil
}

func (client *Client) Dial() error {
	var err error

	opts := grpc.WithInsecure()

	if DEBUG_CLIENT {
		log.Debug("%s dial %s timeout %.1fs", client.Desc(), client.connStr, client.timeout.Seconds())
	}
	log.Info("%s connect %s", client.Desc(), client.connStr)

	client.connection, err = grpc.Dial(client.connStr, opts)
	if err != nil {
		return log.NewError("fail to dial %s: %s", client.connStr, err)
	}
	client.client = facade.NewFacadeClient(client.connection)
	return nil
}

func (client *Client) Desc() string {
	ret := "client["
	ret += client.connStr
	ret += "]"
	return ret

}
