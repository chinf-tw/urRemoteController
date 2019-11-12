package main

import (
	"bytes"
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

		var num int32
		bufferInt := 4
		r := bytes.NewReader(data[:bufferInt])
		binary.Read(r, binary.BigEndian, &num)

		// 匹配messageSize與實際封包數是否相同
		if int32(realLen) != num {
			// 處理理想與實際不同
			fmt.Println("messageSize與實際封包數不同")
			continue // 丟掉該封包
		}
		// Robot messageType should be MESSAGE_TYPE_ROBOT_STATE.
		if data[bufferInt] != MESSAGE_TYPE_ROBOT_STATE {
			fmt.Println("The robot messageType isn't equal MESSAGE_TYPE_ROBOT_STATE")
			continue // 丟掉該封包
		}
		checkSize := 0
		dataMap := make(map[int][]byte)
		for wantBuffer := 0; bufferInt < realLen; {
			bufferInt++ // show next number
			subPackageType := int(data[bufferInt])
			wantBuffer = bufferInt + 4
			subPackageSize := binary.BigEndian.Uint32(data[bufferInt:wantBuffer])
			bufferInt = wantBuffer

			wantBuffer = bufferInt + int(subPackageSize)
			dataMap[subPackageType] = data[bufferInt:wantBuffer]
			bufferInt = wantBuffer

			fmt.Println(subPackageType)
			checkSize += int(subPackageSize)
		}
		if checkSize != realLen {
			fmt.Print("\n*** oh no... ***\n")
			fmt.Println(data[:realLen])
		}
		// fmt.Println(dataMap[1])

		fmt.Println()
		// 釋放data資料
		data = nil
	}

}
