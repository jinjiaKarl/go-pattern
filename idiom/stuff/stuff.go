package stuff

import "fmt"

type StuffClient interface {
	DoStuff() error
}

//two configuration options (timeout and retries)
type stuffClient struct {
	conn    string
	timeout int
	retries int
}

func (s *stuffClient) DoStuff() error {
	fmt.Println(s.conn, "+", s.timeout, "+", s.retries)
	return nil
}

type StuffClientOptions struct {
	Retries int
	Timeout int
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
		options.Timeout = t
	}
}
func WithRetries(r int) StuffClientOption {
	return func(options *StuffClientOptions) {
		options.Retries = r
	}
}
func NewStuffClient(conn string, opts ...StuffClientOption) StuffClient {
	options := StuffClientOptions{
		Timeout: 10,
		Retries: 2,
	}
	for _, o := range opts {
		o(&options)
	}
	return &stuffClient{
		conn:    conn,
		timeout: options.Timeout,
		retries: options.Retries,
	}
}
