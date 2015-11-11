package chatprotocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/gansidui/gotcp"
	"github.com/withliyh/chat/mempool"
)

type HeaderPacket struct {
	PackageLen  uint32
	CommandType uint32
	ContentType uint32
	From        uint32
	To          uint32
}

// Packet
type ChatCommandPacket struct {
	HeaderPacket
	Padding     []byte
	internalBuf []byte
}

func (this *ChatCommandPacket) Serialize() []byte {
	this.internalBuf = mempool.Pool.Get()

	this.PackageLen = 20 + uint32(len(this.Padding))

	binary.BigEndian.PutUint32(this.internalBuf[0:4], uint32(this.PackageLen))
	binary.BigEndian.PutUint32(this.internalBuf[4:8], uint32(this.CommandType))
	binary.BigEndian.PutUint32(this.internalBuf[8:12], uint32(this.ContentType))
	binary.BigEndian.PutUint32(this.internalBuf[12:16], uint32(this.From))
	binary.BigEndian.PutUint32(this.internalBuf[16:20], uint32(this.To))
	copy(this.internalBuf[20:], this.Padding)
	return this.internalBuf[:this.PackageLen]
}

func NewChatCommandPacketWithText(from, to uint32, text string) *ChatCommandPacket {
	packet := &ChatCommandPacket{}
	packet.CommandType = 0
	packet.ContentType = 0
	packet.From = from
	packet.To = to
	packet.Padding = []byte(text)
	return packet
}

type ChatProtocol struct {
}

func (this *ChatProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {

	packet := &ChatCommandPacket{}

	buf := mempool.Pool.Get()
	n, err := io.ReadFull(conn, buf[0:20])
	if err != nil {
		return nil, err
	}
	packet.PackageLen = binary.BigEndian.Uint32(buf[0:4])
	packet.CommandType = binary.BigEndian.Uint32(buf[4:8])
	packet.ContentType = binary.BigEndian.Uint32(buf[8:12])
	packet.From = binary.BigEndian.Uint32(buf[12:16])
	packet.To = binary.BigEndian.Uint32(buf[16:20])

	n, err = io.ReadFull(conn, buf[20:packet.PackageLen-20])
	if err != nil {
		return nil, err
	}
	packet.Padding = buf[20:n]
	return packet, nil
}

type ChatCallback struct {
}

func (this *ChatCallback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	c.AsyncWritePacket(NewChatCommandPacketWithText(0, 0, "hello"), 0)
	return true
}

func (this *ChatCallback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*ChatCommandPacket)
	fmt.Println(string(packet.Padding))
	return true
}

func (this *ChatCallback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}
