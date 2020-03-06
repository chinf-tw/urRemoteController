package urremoteController

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
)

// RealtimeCommunicationsInterface is "Universal Robots" Remote Control Via TCP/IP - 16496
type RealtimeCommunicationsInterface interface {
	Meaning() string     // get the meaning
	SetData([]byte)      // set the data
	GetData() []byte     // get the data
	Output() interface{} // finally Output
}

// RealtimeCommunications Handle all "The Realtime communications Interface" infomation
type RealtimeCommunications struct {
	Meaning        string `json:"meaning"`
	BeginIndex     int    `json:"beginIndex"`
	NumberOfValues int    `json:"Number of values"`
	URType         string `json:"Type"`
	DataSize       int    `json:"Size in bytes"`
	No             int    `json:"No"`
	data           []byte
}

/*RealtimeCommunicationsFormat that
JSON Format with data
	{
	"Controller Timer":{
		"meaning":"Controller Timer",
		"beginIndex":740,
		"Number of values":1,
		"Type":"double",
		"Size in bytes":8,
		"No":19
	},
	"Digital input bits":{
		"meaning":"Digital input bits",
		"beginIndex":684,
		"Number of values":1,
		"Type":"double",
		"Size in bytes":8,
		"No":17
	}
		...
	}
*/
type RealtimeCommunicationsFormat map[string]RealtimeCommunications

func (m RealtimeCommunications) GetData() []byte {
	return m.data
}

func (m *RealtimeCommunications) SetData(d []byte) {
	m.data = d
}

func (m RealtimeCommunications) Output() (interface{}, error) {

	if m.data == nil {
		return nil, errors.New("Error: data is nil")
	}
	if len(m.data) != m.DataSize {
		return nil, errors.New("Error: data size is not match")
	}
	var numberofBits int

	switch m.URType {
	case "double":
		numberofBits = 8
		if m.NumberOfValues > 1 {
			output := []float64{}
			for index := 0; index < m.NumberOfValues; index++ {
				output = append(output, Float64frombytes(m.data[index*numberofBits:(index+1)*numberofBits]))
			}
			return output, nil
		}

		return Float64frombytes(m.data), nil
	case "integer":
		if m.NumberOfValues > 1 {
			output := []uint32{}
			numberofBits = 4
			for index := 0; index < m.NumberOfValues; index++ {
				output = append(output, binary.BigEndian.Uint32(m.data[index*numberofBits:(index+1)*numberofBits]))
			}
			return output, nil
		}

		return binary.BigEndian.Uint32(m.data), nil
	default:
		return nil, fmt.Errorf("Not handle this type %s now", m.URType)
	}
}

func RealTime(conn net.Conn) ([]byte, error) {
	var data = make([]byte, 2048)
	realLen, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	data = data[:realLen]
	bufferInt := 4
	packageSize := binary.BigEndian.Uint32(data[:bufferInt])
	if int(packageSize) != realLen {
		return nil, errors.New("packageSize is not equal realLen")
	}

	tcpForceStart := 540 // 4+8+48*11 = 540 bytes
	tcpForceData := data[tcpForceStart : tcpForceStart+48]

	data = nil
	return tcpForceData, nil
}

// Float64frombytes can convert bytes to float64, that need to give 64 bits (8 bytes) to convert to float64
// and that used BigEndian to convert.
func Float64frombytes(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
