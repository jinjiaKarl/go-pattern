package grpc

import (
	"math"
	"net"
	"sync"
	"time"
)

const (
	defaultServerMaxReceiveMessageSize = 1024 * 1024 * 4
	defaultServerMaxSendMessageSize    = math.MaxInt32
	// http2IOBufSize specifies the buffer size for sending frames.
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

//grpc source code, just show basic use
type Server struct {
	opts serverOptions
	mu   sync.Mutex // guards following
	lis  map[net.Listener]bool
	//.....
}

type serverOptions struct {
	//...
	maxConcurrentStreams  uint32
	maxReceiveMessageSize int
	maxSendMessageSize    int
	//...
	writeBufferSize   int
	readBufferSize    int
	connectionTimeout time.Duration
	//...
}

var defaultServerOptions = serverOptions{
	maxReceiveMessageSize: defaultServerMaxReceiveMessageSize,
	maxSendMessageSize:    defaultServerMaxSendMessageSize,
	connectionTimeout:     120 * time.Second,
	writeBufferSize:       defaultWriteBufSize,
	readBufferSize:        defaultReadBufSize,
}

type ServerOption interface {
	apply(*serverOptions)
}
type EmptyServerOption struct{}

func (EmptyServerOption) apply(*serverOptions) {}

// funcServerOption wraps a function that modifies serverOptions into an
// implementation of the ServerOption interface.
type funcServerOption struct {
	f func(*serverOptions)
}

func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

//return interface
func WriteBUfferSize(s int) ServerOption {
	return newFuncServerOption(func(options *serverOptions) {
		options.writeBufferSize = s
	})
}
func ReadBufferSize(s int) ServerOption {
	return newFuncServerOption(func(options *serverOptions) {
		options.readBufferSize = s
	})
}

func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	s := &Server{
		lis:  make(map[net.Listener]bool), //监听地址列表
		opts: opts,
		//....
	}
	//...
	return s
}
