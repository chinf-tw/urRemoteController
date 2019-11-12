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

	var data []byte
	count := 5
	var dataMap = make(map[byte][]byte)
	isTest := true
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

			bufferInt++
			if bufferInt >= realLen {
				if isTest {
					fmt.Println("*** end!!! ***")
				}

				break
			}
			wantBuffer := bufferInt + 4
			subPackageSize := binary.BigEndian.Uint32(data[bufferInt:wantBuffer])

			if int(subPackageSize) > realLen {
				fmt.Println("*** Error! subPackageSize > Robot packageSize ***")
				fmt.Println(data[bufferInt:wantBuffer])
				break
			}

			wantBuffer = bufferInt + int(subPackageSize)

			if isTest {
				fmt.Println("\n*****")
				fmt.Println("subPackageSize: ", subPackageSize)
				fmt.Println("bufferInt: ", bufferInt)
				fmt.Println("subPackageType: ", data[bufferInt:wantBuffer][4])
				fmt.Println("subPackageSize equal real len: ", int(subPackageSize) == len(data[bufferInt:wantBuffer]))
				fmt.Println(data[bufferInt:wantBuffer])
			}

			dataMap[data[bufferInt:wantBuffer][4]] = data[bufferInt:wantBuffer]
			bufferInt = wantBuffer - 1

		}
		fmt.Println(dataMap[7])
		// 釋放data資料
		data = nil
	}

}
