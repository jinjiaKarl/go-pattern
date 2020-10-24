package stuff

import "fmt"

type StuffClient interface {
	DoStuff() error
}

//two configuration options (timeout and retries)
type stuffClient struct {
	conn string
	opts StuffClientOptions
}

func (s *stuffClient) DoStuff() error {
	fmt.Println(s.conn, "+", s.opts.timeout, "+", s.opts.retries)
	return nil
}

type StuffClientOptions struct {
	retries int
	timeout int
}

var defaultStuffClientOptions = StuffClientOptions{
	retries: 3,
	timeout: 5,
}

/*
//first solution: we could pass a config struct
//That struct is private, so we should provide some sort of constructor
func NewStuffClient(conn string, options StuffClientOptions) StuffClient {
	return &stuffClient{
		conn:    conn,
		timeout: options.Timeout,
		retries: options.Retries,
	}
}
*/

//second solution: Functional Options Pattern
//define StuffClientOption which is just a function which accepts our options struct as a parameter
type StuffClientOption func(*StuffClientOptions)

func WithTimeout(t int) StuffClientOption {
	return func(options *StuffClientOptions) {
		options.timeout = t
	}
}
func WithRetries(r int) StuffClientOption {
	return func(options *StuffClientOptions) {
		options.retries = r
	}
}
func NewStuffClient(conn string, opts ...StuffClientOption) StuffClient {
	options := defaultStuffClientOptions
	for _, o := range opts {
		o(&options)
	}
	return &stuffClient{
		conn: conn,
		opts: options,
	}
}
