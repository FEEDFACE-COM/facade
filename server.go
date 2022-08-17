package main

import (
	"FEEDFACE.COM/facade/facade"
	"FEEDFACE.COM/facade/log"
	"bufio"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"io"
	"net"
	"time"
)

const DEBUG_SERVER = false
const DEBUG_SERVER_DUMP = false

type Server struct {
	host     string
	confPort uint
	textPort uint

	timeout    float64
	bufferSize uint

	transport string

	connStr    string
	connection *grpc.ClientConn

	bufferChan chan facade.TextSeq
	confChan   chan facade.Config
}

func NewServer(host string, confPort uint, textPort uint, timeout float64, noIPv4 bool, noIPv6 bool) *Server {
	ret := Server{
		host:       host,
		confPort:   confPort,
		textPort:   textPort,
		timeout:    timeout,
		bufferSize: TEXT_BUFFER_SIZE,
		transport:  "tcp",
	}

	if noIPv4 && noIPv6 {
		ret.transport = ""
	} else if noIPv4 {
		ret.transport = "tcp6"
	} else if noIPv6 {
		ret.transport = "tcp4"
	}
	return &ret
}

func (server *Server) ListenText(bufChan chan facade.TextSeq) {

	if server.transport == "" {
		return
	}

	textListenStr := fmt.Sprintf("%s:%d", server.host, server.textPort)
	textListener, err := net.Listen(server.transport, textListenStr)
	if err != nil {
		log.PANIC("%s fail listen on %s: %s", server.Desc(), textListenStr, err)
	}
	defer func() { textListener.Close() }()
	log.Notice("%s listen for text on %s", server.Desc(), textListener.Addr())

	for {
		textConn, err := textListener.Accept()
		if err != nil {
			log.Error("%s fail accept on %s: %s", server.Desc(), textListenStr, err)
			continue
		}
		if DEBUG_SERVER {
			log.Debug("%s accept text from %s", server.Desc(), textConn.RemoteAddr().String())
		}
		if server.timeout == 0.0 {
			textConn.SetReadDeadline(time.Time{})
		} else {
			textConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		}
		go server.ReceiveText(textConn, bufChan)

	}
}

func (server *Server) Conf(ctx context.Context, config *facade.Config) (*facade.Status, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		if DEBUG_SERVER {
			log.Debug("%s receive conf from %s", server.Desc(), peer.Addr)
		}
	} else {
		log.Warning("%s receive conf from unknown peer", server.Desc())
	}

	server.confChan <- *config
	return &facade.Status{Success: true}, nil
}

func (server *Server) Pipe(stream facade.Facade_PipeServer) error {

	peer, ok := peer.FromContext(stream.Context())
	if ok {
		if DEBUG_SERVER {
			log.Debug("%s receive stream from %s", server.Desc(), peer.Addr)
		}
	} else {
		log.Warning("%s receive stream from unknown peer", server.Desc())
	}

	var rem []byte = []byte{}
	var tmp []byte
	for {
		msg, err := stream.Recv()
		if err != nil && err != io.EOF {
			if DEBUG_SERVER {
				log.Error("%s fail to receive: %s", server.Desc(), err)
			}
			return log.NewError("fail to receive: %s", err)
		}
		raw := msg.GetRaw()
		if DEBUG_SERVER_DUMP {
			log.Debug("%s recv %d byte raw:\n%s", server.Desc(), len(raw), log.Dump(raw, len(raw), 0))
		} else if DEBUG_SERVER {
			log.Debug("%s recv %d byte raw", server.Desc(), len(raw))
		}
		tmp = append(rem, raw...)
		rem, err = facade.ProcessRaw(tmp, server.bufferChan)
		if err != nil {
			log.Warning("%s fail process raw text: %s", server.Desc(), err)
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
	var buf []byte = make([]byte, server.bufferSize)
	var rem []byte = []byte{}
	var tmp []byte
	reader := bufio.NewReader(textConn)
	if DEBUG_SERVER {
		log.Debug("%s receive text from %s", server.Desc(), textConn.RemoteAddr())
	}
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			if DEBUG_SERVER {
				log.Debug("%s end of text from %s", server.Desc(), textConn.RemoteAddr())
			}
			break
		}
		if err != nil {
			log.Error("%s text read %s error: %s", server.Desc(), textConn.RemoteAddr().String(), err)
			break
		}
		if DEBUG_SERVER_DUMP {
			log.Debug("%s recv %d byte:\n%s", server.Desc(), n, log.Dump(buf, n, 0))
		} else if DEBUG_SERVER {
			log.Debug("%s recv %d byte", server.Desc(), n)
		}
		tmp = append(rem, buf[:n]...)
		//if DEBUG_SERVER {
		//      log.Debug("%s PROCESS %d byte:\n%s",server.Desc(),len(tmp),log.Dump(tmp,len(tmp),0))
		//}
		rem, err = facade.ProcessRaw(tmp, bufChan)
		if err != nil {
			log.Error("%s text process error: %s", server.Desc(), err)
			//if DEBUG_SERVER {
			//      log.Debug("%s RETURN %d byte:\n%s",server.Desc(),len(rem),log.Dump(rem,len(rem),0))
			//}
		}
	}
}

func (server *Server) Listen(
	confChan chan facade.Config,
	bufferChan chan facade.TextSeq,
) {
	var err error

	server.confChan = confChan
	server.bufferChan = bufferChan

	if server.transport == "" {
		return
	}

	server.connStr = fmt.Sprintf("%s:%d", server.host, server.confPort)

	if DEBUG_SERVER {
		log.Debug("%s listen for conf on %s", server.Desc(), server.connStr)
	}
	listener, err := net.Listen(server.transport, server.connStr)
	if err != nil {
		log.PANIC("%s fail listen on %s: %s", server.Desc(), server.connStr, err)
	}

	serv := grpc.NewServer()
	facade.RegisterFacadeServer(serv, &Server{confChan: confChan, bufferChan: bufferChan})

	log.Notice("%s listen for conf on %s", server.Desc(), listener.Addr())

	err = serv.Serve(listener)
	if err != nil {
		log.PANIC("%s fail serve: %s", server.Desc(), err)
	}
	if DEBUG_SERVER {
		log.Debug("%s listen done.", server.Desc())
	}
}

func (server *Server) Desc() string {
	ret := "server["
	ret += fmt.Sprintf("%s", server.transport)
	ret += "]"
	return ret
}
