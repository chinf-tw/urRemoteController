package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	var MESSAGE_TYPE_ROBOT_STATE = byte(16)
	conn, err := net.Dial("tcp", "192.168.1.107:30002")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	// status, err := bufio.NewReader(conn).ReadString('\n')
	// fmt.Println(status)
	var data []byte
	count := 2
	for index := 0; index < count; index++ {
		data = make([]byte, 2048)

		realLen, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		bufferInt := 4
		packageSize := binary.BigEndian.Uint32(data[:bufferInt])
		// 匹配messageSize與實際封包數是否相同
		if realLen != int(packageSize) {
			// 處理理想與實際不同
			fmt.Println("messageSize與實際封包數不同")
			continue // 丟掉該封包
		}
		// Robot messageType should be MESSAGE_TYPE_ROBOT_STATE.
		if data[bufferInt] != MESSAGE_TYPE_ROBOT_STATE {
			fmt.Println("The robot messageType isn't equal MESSAGE_TYPE_ROBOT_STATE")
			continue // 丟掉該封包
		}

		for {
			fmt.Println("\n*****")
			bufferInt++
			if bufferInt >= realLen {
				fmt.Println("*** end!!! ***")
				break
			}
			wantBuffer := bufferInt + 4
			subPackageSize := binary.BigEndian.Uint32(data[bufferInt:wantBuffer])
			// bufferInt = wantBuffer
			// fmt.Println(data[:realLen])
			if subPackageSize > 1000 {
				fmt.Println("*** Error! subPackageSize > 1000 ***")
				fmt.Println(data[bufferInt:wantBuffer])
				break
			}

			wantBuffer = bufferInt + int(subPackageSize)
			fmt.Println("subPackageSize: ", subPackageSize)
			fmt.Println("bufferInt: ", bufferInt)
			fmt.Println(data[bufferInt:wantBuffer])
			fmt.Println("subPackageType: ", data[bufferInt:wantBuffer][4])
			fmt.Println("subPackageSize equal real len: ", int(subPackageSize) == len(data[bufferInt:wantBuffer]))
			bufferInt = wantBuffer - 1

		}

		// fmt.Println(data[:realLen])
		// 釋放data資料
		data = nil
	}

}
