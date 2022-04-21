package server

import (
	"erpc/codec"
	"erpc/service"

	"encoding/json"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
	"errors"
	"strings"
)


/*
 * a server goes through following phase
 * 1. accept a connection
 * 2. start a gorouting to serve the conn
 * 3. read the conn Option, modify the conn into codec
 * 4. serve the codec, read the request for header and body,
 *	  header for method name,
 *    body for argv
 * 5. process
 * 6. response the request
 */
type Server struct {
	services	sync.Map
}


func NewServer() *Server {
	return &Server{}
}


// register a struct as a service
func (server *Server) Register(rcvr interface{}) error {
	s := service.NewService(rcvr)
	if _, dup := server.services.LoadOrStore(s.Name, s); dup {
		return errors.New("erpc: the services has been registered")
	}
	return nil
}

// find a certain method with its service from a serviceMethod name
func (server *Server) findService(serviceMethod string) (svc *service.Service, mtype *service.MethodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("erpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := server.services.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service.Service)
	mtype = svc.Methods[methodName]
	if mtype == nil {
		err = errors.New("erpc server: can't find method " + methodName)
	}
	return
}


func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	if opt.Magic != Magic {
		log.Printf("rpc server: invalid magic number %x", opt.Magic)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	server.serveCodec(f(conn))
}

func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // it's not possible to recover, so close the connection
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

/*
 * read a remote request
 * get method name from header
 * get argv from interface
 */ 
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	
	req.svc, req.mtype, err = server.findService(h.MethodName)
	if err != nil {
		return req, err
	}

	req.argv = req.mtype.NewArgv()
	req.replyv = req.mtype.NewReplyv()

	// make argv pointer
	argvptr := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvptr = req.argv.Addr().Interface()
	}

	// read from the codec
	if err = cc.ReadBody(argvptr); err != nil {
		log.Println("erpc server: read body err:", err)
		return req, err
	}

	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	// concurrency management needed
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO, should call registered rpc methods to get the right replyv
	defer wg.Done()

	err := req.svc.Call(req.mtype, req.argv, req.replyv)
	// contain the error info in the header with a nil body
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(cc, req.h, invalidRequest, sending)
		return
	}

	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

// start a server by handing in a net.listener
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}
