package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/withliyh/chat/protocol"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.1.200:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	//echoProtocol := &echo.EchoProtocol{}

	/*
		packets, err := chatprotocol.NewPacketWithLargetText(2, 1, "hello")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, packet := range packets {
			conn.Write(packet.Serialize())
		}
	*/

	packets, e := chatprotocol.NewPacketWithFile(2, 1, "./server")
	if e != nil {
		fmt.Println(e.Error())
	}
	for _, p := range packets {
		conn.Write(p.Serialize())
	}
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
	time.Sleep(5 * time.Second)

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
