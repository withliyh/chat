package chatprotocol

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/withliyh/chat/mempool"
)

var (
	identifythegenerator uint32
)

func getIdentify() uint32 {
	identifythegenerator++
	return identifythegenerator
}

func NewPacketWithLargetText(from, to uint32, text string) ([]*ChatCommandPacket, error) {
	reader := strings.NewReader(text)
	return NewPacketWithReader(
		COMMAND_TYPE_STREAM,
		CONTENT_TYPE_TEXT,
		from,
		to,
		reader,
	)
}

func NewPacketWithFile(from, to uint32, name string) ([]*ChatCommandPacket, error) {
	f, err := os.Open(name)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	return NewPacketWithReader(
		COMMAND_TYPE_STREAM,
		CONTENT_TYPE_FILE,
		from,
		to,
		f,
	)
}

func NewPacketWithReader(dtype, ctype, from, to uint32, reader io.Reader) ([]*ChatCommandPacket, error) {
	identify := getIdentify()

	packets := make([]*ChatCommandPacket, 0, 1024)
	var seq uint32 = 0
	for {
		buf := mempool.LimitPool.Get()
		n, err := io.ReadFull(reader, buf)
		if err != nil {
			fmt.Println(err.Error())
		}

		packet := &ChatCommandPacket{}
		packet.CommandType = dtype
		packet.ContentType = ctype
		packet.From = from
		packet.To = to
		packet.Padding = buf[0:n]
		packet.Ident = identify
		packet.Seq = seq

		if err == io.EOF {
			if seq == 0 {
				mempool.LimitPool.Give(buf)
				return nil, err
			} else {
				packet.Flag = FLAG_PACKET_END
				packets = append(packets, packet)
				return packets, nil
			}
		}

		if err == io.ErrUnexpectedEOF {
			if seq == 0 {
				packet.Flag = FLAG_PACKET_ONCE
			} else {
				packet.Flag = FLAG_PACKET_END
			}
			packets = append(packets, packet)
			return packets, nil
		}

		if seq == 0 {
			packet.Flag = FLAG_PACKET_START
		} else {
			packet.Flag = FLAG_PACKET_STREAMING
		}
		packets = append(packets, packet)

		seq++
	}
}
