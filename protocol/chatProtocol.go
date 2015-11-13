package chatprotocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/gansidui/gotcp"
	"github.com/withliyh/chat/mempool"
)

const (
	FLAG_PACKET_ONCE = iota
	FLAG_PACKET_START
	FLAG_PACKET_END
	FLAG_PACKET_STREAMING
)
const (
	COMMAND_TYPE_STREAM = iota
	COMMAND_TYPE_PACKET
)

const (
	CONTENT_TYPE_TEXT = iota
	CONTENT_TYPE_IMAGE
	CONTENT_TYPE_AUDIO
	CONTENT_TYPE_VIDEO
	CONTENT_TYPE_FILE
)

type HeaderPacket struct {
	PackageLen  uint32
	CommandType uint32
	ContentType uint32
	From        uint32
	To          uint32
	Seq         uint32
	Flag        uint32
	Ident       uint32
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
	binary.BigEndian.PutUint32(tmpbuf[20:24], uint32(this.Seq))
	binary.BigEndian.PutUint32(tmpbuf[24:28], uint32(this.Flag))
	binary.BigEndian.PutUint32(tmpbuf[28:32], uint32(this.Ident))
	copy(tmpbuf[32:], this.Padding[:])
	mempool.LimitPool.Give(this.Padding)

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
	packet.Seq = binary.BigEndian.Uint32(buf[20:24])
	packet.Flag = binary.BigEndian.Uint32(buf[24:28])
	packet.Ident = binary.BigEndian.Uint32(buf[28:32])

	_, err = io.ReadFull(conn, buf[32:packet.PackageLen])
	if err != nil {
		return nil, err
	}
	packet.Padding = buf[32:packet.PackageLen]
	//	fmt.Println(packet.internalBuf[0:packet.PackageLen])
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

var recvfile *os.File

func (this *ChatCallback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*ChatCommandPacket)
	if packet.Flag == FLAG_PACKET_START {
		recvfile, _ = os.Create("./recv.raw")
	} else if packet.Flag == FLAG_PACKET_END {
		recvfile.Close()
		return true
	}
	recvfile.Write(packet.Padding)
	//	s := string(packet.Padding)
	//	fmt.Println(s)
	mempool.Pool.Give(packet.GetInternalBuf())
	return true
}

func (this *ChatCallback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}
