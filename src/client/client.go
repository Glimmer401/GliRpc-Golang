package client

import (
	"erpc/codec"
	"erpc/server"
	"encoding/json"
	"io"
	"sync"
	"net"
	"log"
	"errors"
	"fmt"
)

var ErrShutdown = errors.New("connection is shut down")

type Client struct {
	cc			codec.Codec
	opt			*server.Option
	sending		sync.Mutex	 // lock for cc codec
	mtx  		sync.Mutex	 // lock for Client
	header		codec.Header // send when request, as atomic operation, only need one copy
	seq			uint64		 // next seq
	pending		map[uint64]*Call	// call in flight
	isClose		bool	// close by user
	isShutdown	bool	// shutdown by server
}

var _ io.Closer = (*Client)(nil)

// Close the connection
func (client *Client) Close() error {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	// exception: shutdown by server
	if client.isClose {
		return ErrShutdown
	}
	// user close
	client.isClose = true
	return client.cc.Close()
}

// IsAvailable return true if the client does work
func (client *Client) IsAvailable() bool {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	return !client.isShutdown && !client.isClose
}

// request for a call
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	// the Client is not working
	if client.isClose || client.isShutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call  // track outstanding call
	client.seq++
	return call.Seq, nil
}

func (client *Client) removeCall(seq uint64) *Call {
	client.mtx.Lock()
	defer client.mtx.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mtx.Lock()
	defer client.mtx.Unlock()
	client.isShutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

// receive during the client work flow
func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}
		// the call has finished and handle its out-come
		call := client.removeCall(h.Seq)
		switch {
		case call == nil:	// invalid call but handled by server
			err = client.cc.ReadBody(nil)
		case h.Error != "":	// error happend for the call
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			call.done()
		default:			// everything is fine
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	// error occurs, so terminateCalls pending calls
	client.terminateCalls(err)
}

func (client *Client) send(call *Call) {
	// send a complete request without interupt
	client.sending.Lock()
	defer client.sending.Unlock()

	// register this call in flight
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	// prepare request header
	client.header.MethodName = call.MethodName
	client.header.Seq = seq
	client.header.Error = ""

	// encode and send the request
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		call := client.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

// Go invokes the function asynchronously.
func (client *Client) Go(methodName string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}
	call := &Call{
		MethodName: methodName,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	client.send(call)
	return call
}

// Call invokes the named function, waits for it to complete
func (client *Client) Call(methodName string, args, reply interface{}) error {
	call := <-client.Go(methodName, args, reply, make(chan *Call, 1)).Done
	return call.Error
}


func newClientCodec(cc codec.Codec, opt *server.Option) *Client {
	client := &Client{
		seq:     1, // seq starts with 1, 0 means invalid call
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
	}
	// start a thread receiving
	go client.receive()
	return client
}

// start a new client
func NewClient(conn net.Conn, opt *server.Option) (*Client, error) {
	// with option, create a codec
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("erpc client: codec error:", err)
		return nil, err
	}
	// send the option to the server
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("erpc client: option error: ", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}


func parseOptions(opts ...*server.Option) (*server.Option, error) {
	// if opts is nil or pass nil as parameter
	if len(opts) == 0 || opts[0] == nil {
		return server.DefaultOption, nil
	}
	if len(opts) != 1 {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.Magic = server.DefaultOption.Magic
	if opt.CodecType == "" {
		opt.CodecType = server.DefaultOption.CodecType
	}
	return opt, nil
}

// Dial connects to an RPC server at the specified network address
func Dial(network, address string, opts ...*server.Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	// close the connection if client is nil
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()
	return NewClient(conn, opt)
}
