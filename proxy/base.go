package proxy

import (
	"context"
	"errors"
	"net"

	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy/proto"
)

var _ Proxy = (*Base)(nil)

type Base struct {
	addr  string
	proto proto.Proto
}

func (b *Base) Addr() string {
	return b.addr
}

func (b *Base) Proto() proto.Proto {
	return b.proto
}

func (b *Base) IsChan() bool {
	return proto.Chan == b.proto
}

func (b *Base) PassTcp(conn adapter.TCPConn) error {
	return errors.ErrUnsupported
}

func (b *Base) PassUdp(conn adapter.UDPConn) error {
	return errors.ErrUnsupported
}

func (b *Base) DialContext(context.Context, *M.Metadata) (net.Conn, error) {
	return nil, errors.ErrUnsupported
}

func (b *Base) DialUDP(*M.Metadata) (net.PacketConn, error) {
	return nil, errors.ErrUnsupported
}
