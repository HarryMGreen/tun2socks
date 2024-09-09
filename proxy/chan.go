package proxy

import (
	"net"
	"net/netip"

	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	M "github.com/xjasonlyu/tun2socks/v2/metadata"
	"github.com/xjasonlyu/tun2socks/v2/proxy/proto"
)

var _ Proxy = (*Chan)(nil)

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

type Chan struct {
	*Base
}

type PassData struct {
	Metadata   M.Metadata
	Conn       net.Conn
	PacketConn net.PacketConn
}

func NewChan() *Chan {
	return &Chan{
		Base: &Base{
			proto: proto.Chan,
		},
	}
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
