package urRemoteController

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

var (
	DataNotMatchErr error
)

// GetRealtimeCommunicationsFormat Get RealtimeCommunicationsFormat, RealtimeCommunicationsFormat is important with this package
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

// CommunicationsFloat64 UR pose struct
type CommunicationsFloat64 struct {
	X  float64 `json:"X"`
	Y  float64 `json:"Y"`
	Z  float64 `json:"Z"`
	RX float64 `json:"RX"`
	RY float64 `json:"RY"`
	RZ float64 `json:"RZ"`
}

// Add add all values
func (source *CommunicationsFloat64) Add(target *CommunicationsFloat64) {
	source.X += target.X
	source.Y += target.Y
	source.Z += target.Z
	source.RX += target.RX
	source.RY += target.RY
	source.RZ += target.RZ
}

// RunURWithMoveJ Using URScript with movej to remote control
func RunURWithMoveJ(rCFormat RealtimeCommunicationsFormat, conn net.Conn, communications CommunicationsFloat64, timeout time.Duration) ([]float64, error) {
	var (
		data       []byte
		err        error
		actualpose []float64
		str        string
	)
	if data, err = Read(conn, rCFormat, timeout); err != nil {
		return nil, err
	}

	if actualpose, err = DecodeActualPose(rCFormat, data); err != nil {
		return nil, err
	}
	actualpose = AddCommunications(actualpose, communications)
	if str, err = Float64ToPose(actualpose); err != nil {
		return nil, err
	}

	moveStr := fmt.Sprintf("movej(%s)\n", str)
	if _, err := conn.Write([]byte(moveStr)); err != nil {
		return nil, err
	}

	return actualpose, nil
}

// Read Read UR data with remote control protocol
func Read(conn net.Conn, rCFormat RealtimeCommunicationsFormat, timeout time.Duration) ([]byte, error) {
	var (
		dataLen         int
		err             error
		targetInterface interface{}
		targetLen       uint32
		done            chan bool
	)
	var data = make([]byte, 2048)

	if dataLen, err = conn.Read(data); err != nil {
		return nil, err
	}

	// data = data[:dataLen]
	if dataLen == 0 {
		return nil, fmt.Errorf("Error: dataLen is 0")
	}
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
	go func() {
		time.Sleep(timeout)
		done <- true
	}()
	for int(targetLen) != dataLen {
		if dataLen, err = conn.Read(data); err != nil {
			return nil, err
		}
		select {
		case <-done:
			return nil, fmt.Errorf("Error: read UR data Timeout %d sec", timeout/time.Second)
		default:
		}
	}
	if int(targetLen) != dataLen {
		return nil, fmt.Errorf("Error: int(targetLen) != dataLen")
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("Error: data is empty")
	}

	return data, nil
}

// AddCommunications Add communications
func AddCommunications(currently []float64, communications CommunicationsFloat64) []float64 {
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

// Float64ToPose Convert float list to UR pose
func Float64ToPose(target []float64) (string, error) {
	var (
		err error
		str string
	)
	str = "p["
	for _, f := range target {
		s := strconv.FormatFloat(f, 'f', -1, 64)
		if err != nil {
			return "", err
		}
		str += s + ","
	}
	str = str[:len(str)-1]
	str += "]"
	return str, nil
}

// DecodeActualPose Deconde UR remote Control Protocol, data need to be all bytes
func DecodeActualPose(rCFormat RealtimeCommunicationsFormat, data []byte) ([]float64, error) {
	var (
		actualposeI interface{}
		err         error
		actualpose  []float64
	)
	toolVectorActual := rCFormat["Tool vector actual"]
	begin := toolVectorActual.BeginIndex
	end := toolVectorActual.BeginIndex + toolVectorActual.DataSize
	if len(data) < end {
		return nil, fmt.Errorf("Error: data length: %d less than toolVectorActual end index: %d", len(data), end)
	}
	toolVectorActual.SetData(data[begin:end])

	if actualposeI, err = toolVectorActual.Output(); err != nil {
		return nil, err
	}

	// 轉換型別
	switch actualposeI.(type) {
	case []float64:
		actualpose = actualposeI.([]float64)
		if len(actualpose) != toolVectorActual.NumberOfValues {
			return nil, fmt.Errorf("Error: actualpose is not match toolVectorActual.NumberOfValues")
		}
		return actualpose, nil
	default:
		return nil, fmt.Errorf("Error: target interface type is not a []float64")
	}

}
