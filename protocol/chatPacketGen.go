package chatprotocol

import (
	"io"
	"strings"

	"github.com/withliyh/chat/mempool"
)

const (
	FLAG_PACKET_ONCE = iota
	FLAG_PACKET_START
	FLAG_PACKET_END
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

var (
	identifythegenerator uint32
)

func getIdentify() uint32 {
	identifythegenerator++
	return identifythegenerator
}

func NewChatCommandPacketWithLargetText(from, to uint32, text string) ([]*ChatCommandPacket, error) {
	reader := strings.NewReader(text)
	return NewChatCommandPacketWithReader(
		COMMAND_TYPE_STREAM,
		CONTENT_TYPE_TEXT,
		from,
		to,
		reader,
	)
}

func NewChatCommandPacketWithReader(dtype, ctype, from, to uint32, reader io.Reader) ([]*ChatCommandPacket, error) {
	identify := getIdentify()

	buf := mempool.Pool.Get()
	n, err := io.ReadFull(reader, buf)
	if err == io.EOF {
		mempool.Pool.Give(buf)
		return nil, err
	}

	packets := make([]*ChatCommandPacket, 0, 1024)

	packet := &ChatCommandPacket{}
	packet.ident = identify
	packet.CommandType = dtype
	packet.ContentType = ctype
	packet.From = from
	packet.To = to
	packet.Padding = buf[0:n]

	if err == io.ErrUnexpectedEOF {
		packet.flag = FLAG_PACKET_ONCE
		packets = append(packets, packet)
		return packets, nil
	}

	var seq uint32 = 0
	packet.flag = FLAG_PACKET_START
	for {
		buf := mempool.Pool.Get()
		n, err := io.ReadFull(reader, buf)
		packet := &ChatCommandPacket{}
		packet.CommandType = dtype
		packet.ContentType = ctype
		packet.From = from
		packet.To = to
		packet.Padding = buf[0:n]
		packet.ident = identify
		packet.seq = seq
		seq++
		if err != nil {
			packet.flag = FLAG_PACKET_END
			packets = append(packets, packet)
			return packets, nil
		}
		packets = append(packets, packet)
	}

	return packets, nil
}
