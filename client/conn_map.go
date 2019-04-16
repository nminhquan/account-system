package client

import (
	"google.golang.org/grpc"
)

type ConnPool interface {
	Get(address string) (*grpc.ClientConn, error)
	Insert(address string, conn *grpc.ClientConn) error
}

type ConnMap struct {
	pool map[string]*grpc.ClientConn
}

func NewConnMap() *ConnMap {
	pool := make(map[string]*grpc.ClientConn)
	return &ConnMap{pool}
}

func (p *ConnMap) Get(address string) (*grpc.ClientConn, error) {
	conn := p.pool[address]
	if conn == nil {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		p.Insert(address, conn)
		return conn, err
	}
	return conn, nil
}

func (p *ConnMap) Insert(address string, conn *grpc.ClientConn) error {
	p.pool[address] = conn
	return nil
}

func (p *ConnMap) GetMap() map[string]*grpc.ClientConn {
	return p.pool
}
