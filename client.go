package main

import (
	"log"
	"net"

	"github.com/withliyh/chat/protocol"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	//echoProtocol := &echo.EchoProtocol{}

	packet := chatprotocol.NewChatCommandPacketWithText(0, 1, "hello")
	conn.Write(packet.Serialize())
	/*
			// read
			p, err := echoProtocol.ReadPacket(conn)
			if err == nil {
				echoPacket := p.(*echo.EchoPacket)
				fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
			}

			time.Sleep(2 * time.Second)
		}
	*/

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
