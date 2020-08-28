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

const DEBUG_SERVER = false
const DEBUG_SERVER_DUMP = false

type Server struct {
	host     string
	confPort uint
	textPort uint
	timeout  float64

	transport string

	connStr    string
	connection *grpc.ClientConn

	bufferChan chan facade.TextSeq
	confChan   chan facade.Config
	queryChan  chan (chan string)
}

func NewServer(host string, confPort uint, textPort uint, timeout float64, forceIPv4 bool, forceIPv6 bool) *Server {
	ret := Server{host: host, confPort: confPort, textPort: textPort, timeout: timeout, transport: "tcp"}
	if forceIPv4 {
		ret.transport = "tcp4"
	} else if forceIPv6 {
		ret.transport = "tcp6"
	}
	return &ret
}

func (server *Server) ListenText(bufChan chan facade.TextSeq) {
	textListenStr := fmt.Sprintf("%s:%d", server.host, server.textPort)
	textListener, err := net.Listen(server.transport, textListenStr)
	if err != nil {
		log.PANIC("%s fail listen on %s: %s", server.Desc(), textListenStr, err)
	}
	defer func() { /*log.Debug("stop listen text on %s",textListener.Addr().String());*/ textListener.Close() }()
	log.Info("%s text on %s", server.Desc(), textListener.Addr().String())

	for {
		textConn, err := textListener.Accept()
		if err != nil {
			log.Error("%s fail accept on %s: %s", server.Desc(), textListenStr, err)
			continue
		}
		if DEBUG_SERVER {
			log.Info("%s receive text from %s", server.Desc(), textConn.RemoteAddr().String())
		}
		if server.timeout == 0.0 {
			textConn.SetReadDeadline(time.Time{})
		} else {
			textConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		}
		go server.ReceiveText(textConn, bufChan)

	}
}

func (server *Server) Info(ctx context.Context, empty *facade.Empty) (*facade.Status, error) {
	if DEBUG_SERVER {
		log.Debug("%s received info request", server.Desc())
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
			log.Debug("%s query channel time out", server.Desc())
		}
		return &facade.Status{Success: false, Error: "timeout"}, log.NewError("timeout")

	}

	if DEBUG_SERVER {
		log.Debug("%s respond query info: %s", server.Desc(), ret.Info)
	}

	ret.Success = true
	return ret, nil

}

func (server *Server) Conf(ctx context.Context, config *facade.Config) (*facade.Status, error) {
	if DEBUG_SERVER {
		log.Debug("%s receive conf %s", server.Desc(), config.Desc())
	}

	server.confChan <- *config
	return &facade.Status{Success: true}, nil
}

func (server *Server) Pipe(stream facade.Facade_PipeServer) error {
	//	if DEBUG_SERVER {
	log.Info("%s receive text stream", server.Desc())
	//	}
	var rem []byte = []byte{}
	var tmp []byte
	for {
		msg, err := stream.Recv()
		if err != nil && err != io.EOF {
			if DEBUG_SERVER {
				log.Debug("%s fail to receive: %s", server.Desc(), err)
			}
			return log.NewError("fail to receive: %s", err)
		}
		raw := msg.GetRaw()
		if DEBUG_SERVER {
			if DEBUG_SERVER_DUMP {
				log.Debug("%s recv %d byte raw:\n%s", server.Desc(), len(raw), log.Dump(raw, len(raw), 0))
			} else {
				log.Debug("%s recv %d byte raw", server.Desc(), len(raw))
			}
		}
		tmp = append(rem, raw...)
		rem, err = facade.ProcessRaw(tmp, server.bufferChan)
		if err != nil {
			log.Error("%s error processing raw text: %s", server.Desc(), err)
		}
		if err == io.EOF {
			if DEBUG_SERVER {
				log.Debug("%s recv end of file", server.Desc())
			}
			break
		}

	}
	return stream.SendAndClose(&facade.Status{Success: true})
}

func (server *Server) ReceiveText(textConn net.Conn, bufChan chan facade.TextSeq) {
	defer func() {
		if DEBUG_SERVER {
			log.Debug("%s close text %s", server.Desc(), textConn.RemoteAddr().String())
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
			log.Error("%s text read %s error: %s", server.Desc(), textConn.RemoteAddr().String(), err)
			break
		}
		if DEBUG_SERVER {
			if DEBUG_SERVER_DUMP {
				log.Debug("%s recv %d byte:\n%s", server.Desc(), n, log.Dump(buf, n, 0))
			} else {
				log.Debug("%s recv %d byte", server.Desc(), n)
			}
		}
		tmp = append(rem, buf[:n]...)
		//		log.Debug("%s PROCESS %d byte:\n%s",server.Desc(),len(tmp),log.Dump(tmp,len(tmp),0))
		rem, err = facade.ProcessRaw(tmp, bufChan)
		if err != nil {
			log.Error("%s text process error: %s", server.Desc(), err)
			//            log.Debug("%s RETURN %d byte:\n%s",server.Desc(),len(rem),log.Dump(rem,len(rem),0))
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
		log.Info("%s listen %s", server.Desc(), server.connStr)
	}
	listener, err := net.Listen(server.transport, server.connStr)
	if err != nil {
		log.PANIC("%s fail to listen %s: %s", server.Desc(), server.connStr, err)
	}

	serv := grpc.NewServer()
	facade.RegisterFacadeServer(serv, &Server{confChan: confChan, bufferChan: bufferChan, queryChan: queryChan})
	//	if DEBUG_SERVER {
	log.Info("%s serve on %s", server.Desc(), server.connStr)
	//	}
	err = serv.Serve(listener)
	if err != nil {
		log.PANIC("%s fail to serve: %s", server.Desc(), err)
	}
	if DEBUG_SERVER {
		log.Debug("%s listen done.", server.Desc())
	}
}

func (server *Server) Desc() string {
	ret := "server["
	ret += server.transport
	ret += "]"
	return ret
}
