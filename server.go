package main

import (
	"fmt"
	"io"
	"net"
	"time"
	//    "os"
	facade "./facade"
	log "./log"
	"bufio"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const DEBUG_SERVER = true
const DEBUG_SERVER_DUMP = false

type Server struct {
	host     string
	confPort uint
	textPort uint
	timeout  float64

	connStr    string
	connection *grpc.ClientConn

	bufferChan chan facade.TextSeq
	confChan   chan facade.Config
	queryChan  chan (chan string)
}

func NewServer(host string, confPort uint, textPort uint, timeout float64) *Server {
	return &Server{host: host, confPort: confPort, textPort: textPort, timeout: timeout}
}

func (server *Server) ListenText(bufChan chan facade.TextSeq) {
	textListenStr := fmt.Sprintf("%s:%d", server.host, server.textPort)
	log.Debug("listen for text on %s", textListenStr)
	textListener, err := net.Listen("tcp", textListenStr)
	if err != nil {
		log.PANIC("fail listen on %s: %s", textListenStr, err)
	}
	defer func() { /*log.Debug("stop listen text on %s",textListener.Addr().String());*/ textListener.Close() }()
	log.Info("listening for text on %s", textListener.Addr().String())

	for {
		textConn, err := textListener.Accept()
		if err != nil {
			log.Error("fail accept on %s: %s", textListenStr, err)
			continue
		}
		if DEBUG_SERVER {
			log.Debug("accept text from %s", textConn.RemoteAddr().String())
		}
		if server.timeout == 0.0 {
			textConn.SetReadDeadline(time.Time{})
		} else {
			textConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		}
		go server.ReceiveText(textConn, bufChan)

	}
}

func (server *Server) Query(ctx context.Context, empty *facade.Empty) (*facade.Status, error) {
	if DEBUG_SERVER {
		log.Debug("received info request")
	}
	ret := &facade.Status{}

	chn := make(chan string)
	server.queryChan <- chn

	info := ""
	select {
	case info = <-chn:
		ret.Info = info

	case <-time.After(5. * time.Second):
		if DEBUG_SERVER {
			log.Debug("query channel time out")
		}
		return &facade.Status{Success: false, Error: "timeout"}, log.NewError("timeout")

	}

	if DEBUG_SERVER {
		log.Debug("respond query info: %s", ret.Info)
	}

	ret.Success = true
	return ret, nil

}

func (server *Server) Configure(ctx context.Context, config *facade.Config) (*facade.Status, error) {
	if DEBUG_SERVER {
		log.Debug("receive conf %s", config.Desc())
	}

	server.confChan <- *config
	return &facade.Status{Success: true}, nil
}

func (server *Server) Display(stream facade.Facade_DisplayServer) error {
	if DEBUG_SERVER {
		log.Debug("receive text stream")
	}
	var rem []byte = []byte{}
	var tmp []byte
	for {
		msg, err := stream.Recv()
		if err != nil && err != io.EOF {
			if DEBUG_SERVER {
				log.Debug("fail to receive: %s", err)
			}
			return log.NewError("fail to receive: %s", err)
		}
		raw := msg.GetRaw()
		if DEBUG_SERVER_DUMP {
			log.Debug("recv %d byte raw:\n%s", len(raw), log.Dump(raw, len(raw), 0))
		} else if DEBUG_SERVER {
			log.Debug("recv %d byte raw", len(raw))
		}
		tmp = append(rem, raw...)
		rem, err = facade.ProcessRaw(tmp, server.bufferChan)
		if err != nil {
			log.Error("error processing raw text: %s", err)
		}
		if err == io.EOF {
			if DEBUG_SERVER {
				log.Debug("recv end of file")
			}
			break
		}

	}
	return stream.SendAndClose(&facade.Status{Success: true})
}

func (server *Server) ReceiveText(textConn net.Conn, bufChan chan facade.TextSeq) {
	defer func() {
		if DEBUG_SERVER {
			log.Debug("close text %s", textConn.RemoteAddr().String())
		}
		textConn.Close()
	}()
	const BUFFER_SIZE = 1024
	var buf []byte = make([]byte, BUFFER_SIZE)
	var rem []byte = []byte{}
	var tmp []byte
	reader := bufio.NewReader(textConn)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("text read %s error: %s", textConn.RemoteAddr().String(), err)
			break
		}
		if DEBUG_SERVER_DUMP {
			log.Debug("recv %d byte:\n%s", n, log.Dump(buf, n, 0))
		} else if DEBUG_SERVER {
			log.Debug("recv %d byte", n)
		}

		tmp = append(rem, buf[:n]...)
		//		log.Debug("PROCESS %d byte:\n%s",len(tmp),log.Dump(tmp,len(tmp),0))
		rem, err = facade.ProcessRaw(tmp, bufChan)
		if err != nil {
			log.Error("text process error: %s", err)
			//            log.Debug("RETURN %d byte:\n%s",len(rem),log.Dump(rem,len(rem),0))
		}
	}
}

func (server *Server) Listen(
	confChan chan facade.Config,
	bufferChan chan facade.TextSeq,
	queryChan chan (chan string),
) {
	var err error

	server.confChan = confChan
	server.bufferChan = bufferChan
	server.queryChan = queryChan

	server.connStr = fmt.Sprintf("%s:%d", server.host, server.confPort)

	if DEBUG_SERVER {
		log.Debug("listen %s", server.connStr)
	}
	listener, err := net.Listen("tcp", server.connStr)
	if err != nil {
		log.PANIC("fail to listen %s: %s", server.connStr, err)
	}

	serv := grpc.NewServer()
	facade.RegisterFacadeServer(serv, &Server{confChan: confChan, bufferChan: bufferChan, queryChan: queryChan})
	if DEBUG_SERVER {
		log.Debug("serve %s", server.connStr)
	}
	err = serv.Serve(listener)
	if err != nil {
		log.PANIC("fail to serve: %s", err)
	}
	if DEBUG_SERVER {
		log.Debug("listen done.")
	}
}
