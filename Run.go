package DeCodeURInterface

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

var (
	DataNotMatchErr error
)

func GetRealtimeCommunicationsFormat(jsonFile string) (RealtimeCommunicationsFormat, error) {
	var (
		output RealtimeCommunicationsFormat
		format []byte
		err    error
	)
	if format, err = ioutil.ReadFile(jsonFile); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(format, &output); err != nil {
		return nil, err
	}
	return output, nil
}

type CommunicationsFloat64 struct {
	X  float64
	Y  float64
	Z  float64
	RX float64
	RY float64
	RZ float64
}

func RunURWithMoveJ(rCFormat RealtimeCommunicationsFormat, conn net.Conn, communications CommunicationsFloat64) error {
	var (
		data        []byte
		err         error
		actualpose  []float64
		actualposeI interface{}
		str         = "p["
	)
	if data, err = read(conn, rCFormat); err != nil {
		return err
	}
	toolVectorActual := rCFormat["Tool vector actual"]
	toolVectorActual.SetData(data[toolVectorActual.BeginIndex : toolVectorActual.BeginIndex+toolVectorActual.DataSize])

	if actualposeI, err = toolVectorActual.Output(); err != nil {
		return err
	}

	// 轉換型別
	switch actualposeI.(type) {
	case []float64:
		actualpose = actualposeI.([]float64)
		if len(actualpose) != toolVectorActual.NumberOfValues {
			return fmt.Errorf("Error: actualpose is not match toolVectorActual.NumberOfValues")
		}
		actualpose = addCommunications(actualpose, communications)
	default:
		return fmt.Errorf("Error: target interface type is not a []float64")
	}

	for _, f := range actualpose {
		s := strconv.FormatFloat(f, 'f', -1, 64)
		if err != nil {
			return err
		}
		str += s + ","
	}
	str = str[:len(str)-1]
	str += "]"

	moveStr := fmt.Sprintf("movej(%s)\n", str)
	if _, err := conn.Write([]byte(moveStr)); err != nil {
		return err
	}

	return nil
}

func read(conn net.Conn, rCFormat RealtimeCommunicationsFormat) ([]byte, error) {
	var (
		dataLen         int
		err             error
		targetInterface interface{}
		targetLen       uint32
	)
	var data = make([]byte, 2048)

	if dataLen, err = conn.Read(data); err != nil {
		return nil, err
	}

	data = data[:dataLen]
	messageSize := rCFormat["Message Size"]
	messageSize.SetData(data[messageSize.BeginIndex:messageSize.DataSize])

	if targetInterface, err = messageSize.Output(); err != nil {
		return nil, err
	}

	// 轉換型別
	switch targetInterface.(type) {
	case uint32:
		targetLen = targetInterface.(uint32)
	default:
		return nil, fmt.Errorf("Error: target interface type is not a int")
	}

	if int(targetLen) != dataLen {
		return nil, DataNotMatchErr
	}

	return data, nil
}

func addCommunications(currently []float64, communications CommunicationsFloat64) []float64 {
	if len(currently) != 6 {
		return nil
	}
	currently[0] += communications.X
	currently[1] += communications.Y
	currently[2] += communications.Z
	currently[3] += communications.RX
	currently[4] += communications.RY
	currently[5] += communications.RZ
	return currently
}
