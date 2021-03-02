// +build !confonly

package encoding

import (
	"bytes"
	"context"
	"io"
	"net"
	"time"
)

type ClientConn struct {
	client GunService_TunClient
	reader io.Reader
	over   context.CancelFunc
}

func (s *ClientConn) Read(b []byte) (n int, err error) {
	if s.reader == nil {
		h, err := s.client.Recv()
		if err != nil {
			return 0, newError("unable to read from gun tunnel").Base(err)
		}
		s.reader = bytes.NewReader(h.Data)
	}
	n, err = s.reader.Read(b)
	if err == io.EOF {
		return 0, nil
	}
	return 0, err
}

func (s *ClientConn) Write(b []byte) (n int, err error) {
	err = s.client.Send(&Hunk{Data: b[:]})
	if err != nil {
		return 0, newError("Unable to send data over gun").Base(err)
	}
	return len(b), nil
}

func (s *ClientConn) Close() error {
	return s.client.CloseSend()
}

func (s ClientConn) LocalAddr() net.Addr {
	panic("implement me")
}

func (s ClientConn) RemoteAddr() net.Addr {
	panic("implement me")
}

func (s ClientConn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (s ClientConn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (s ClientConn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}

func NewClientConn(client GunService_TunClient) *ClientConn {
	return &ClientConn{
		client: client,
		reader: nil,
	}
}
