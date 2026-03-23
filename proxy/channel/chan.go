package channel

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"net/url"

	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy"
)

var _ proxy.Proxy = (*Chan)(nil)

var (
	_tcpQueue = make(chan PassData)
	_udpQueue = make(chan PassData)
)

// TCPOut return fan-out TCP queue.
func TCPOut() <-chan PassData {
	return _tcpQueue
}

// UDPOut return fan-out UDP queue.
func UDPOut() <-chan PassData {
	return _udpQueue
}

type Chan struct{}

type PassData struct {
	Metadata   M.Metadata
	Conn       net.Conn
	PacketConn net.PacketConn
}

func New() (*Chan, error) {
	return &Chan{}, nil
}

func (r *Chan) DialContext(context.Context, *M.Metadata) (net.Conn, error) {
	return nil, errors.ErrUnsupported
}

func (r *Chan) DialUDP(*M.Metadata) (net.PacketConn, error) {
	return nil, errors.ErrUnsupported
}

func (c *Chan) PassTcp(conn adapter.TCPConn) error {
	id := conn.ID()
	srcIP, _ := netip.AddrFromSlice(id.RemoteAddress.AsSlice())
	dstIP, _ := netip.AddrFromSlice(id.LocalAddress.AsSlice())
	metadata := M.Metadata{
		Network: M.TCP,
		SrcIP:   srcIP,
		SrcPort: id.RemotePort,
		DstIP:   dstIP,
		DstPort: id.LocalPort,
	}
	pd := PassData{
		Metadata: metadata,
		Conn:     conn,
	}

	_tcpQueue <- pd
	return nil
}

func (b *Chan) PassUdp(conn adapter.UDPConn) error {
	id := conn.ID()
	srcIP, _ := netip.AddrFromSlice(id.RemoteAddress.AsSlice())
	dstIP, _ := netip.AddrFromSlice(id.LocalAddress.AsSlice())
	metadata := M.Metadata{
		Network: M.UDP,
		SrcIP:   srcIP,
		SrcPort: id.RemotePort,
		DstIP:   dstIP,
		DstPort: id.LocalPort,
	}
	pd := PassData{
		Metadata:   metadata,
		PacketConn: conn,
	}

	_udpQueue <- pd
	return nil
}

func Parse(*url.URL) (proxy.Proxy, error) { return New() }

func init() {
	proxy.RegisterProtocol("chan", Parse)
}
