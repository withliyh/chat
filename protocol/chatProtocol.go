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
	seq         uint32
	flag        uint32
	ident       uint32
}

// Packet
type ChatCommandPacket struct {
	HeaderPacket
	Padding     []byte
	internalBuf []byte
}

func (this *ChatCommandPacket) Serialize() []byte {
	this.PackageLen = 32 + uint32(len(this.Padding))
	this.internalBuf = mempool.Pool.Get()
	tmpbuf := this.internalBuf
	binary.BigEndian.PutUint32(tmpbuf[0:4], uint32(this.PackageLen))
	binary.BigEndian.PutUint32(tmpbuf[4:8], uint32(this.CommandType))
	binary.BigEndian.PutUint32(tmpbuf[8:12], uint32(this.ContentType))
	binary.BigEndian.PutUint32(tmpbuf[12:16], uint32(this.From))
	binary.BigEndian.PutUint32(tmpbuf[16:20], uint32(this.To))
	binary.BigEndian.PutUint32(tmpbuf[20:24], uint32(this.seq))
	binary.BigEndian.PutUint32(tmpbuf[24:28], uint32(this.flag))
	binary.BigEndian.PutUint32(tmpbuf[28:32], uint32(this.ident))
	copy(tmpbuf[32:], this.Padding[:])
	mempool.Pool.Give(this.Padding)

	fmt.Println(tmpbuf[0:this.PackageLen])
	fmt.Println(len(tmpbuf[0:this.PackageLen]))
	return tmpbuf[0:this.PackageLen]
}

func (this *ChatCommandPacket) GetInternalBuf() []byte {
	return this.internalBuf
}

type ChatProtocol struct {
}

func (this *ChatProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {

	packet := &ChatCommandPacket{}
	packet.internalBuf = mempool.Pool.Get()
	buf := packet.internalBuf
	_, err := io.ReadFull(conn, buf[0:32])
	if err != nil {
		return nil, err
	}
	packet.PackageLen = binary.BigEndian.Uint32(buf[0:4])
	packet.CommandType = binary.BigEndian.Uint32(buf[4:8])
	packet.ContentType = binary.BigEndian.Uint32(buf[8:12])
	packet.From = binary.BigEndian.Uint32(buf[12:16])
	packet.To = binary.BigEndian.Uint32(buf[16:20])
	packet.seq = binary.BigEndian.Uint32(buf[20:24])
	packet.flag = binary.BigEndian.Uint32(buf[24:28])
	packet.ident = binary.BigEndian.Uint32(buf[28:32])

	_, err = io.ReadFull(conn, buf[32:packet.PackageLen])
	if err != nil {
		return nil, err
	}
	packet.Padding = buf[32:packet.PackageLen]
	fmt.Println(packet.internalBuf[0:packet.PackageLen])
	return packet, nil
}

type ChatCallback struct {
}

func (this *ChatCallback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	//	c.AsyncWritePacket(NewChatCommandPacketWithText(0, 0, "hello"), 0)
	return true
}

func (this *ChatCallback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*ChatCommandPacket)
	s := string(packet.Padding)
	fmt.Println(s)
	mempool.Pool.Give(packet.GetInternalBuf())
	return true
}

func (this *ChatCallback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}
