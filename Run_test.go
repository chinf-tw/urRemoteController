package DeCodeURInterface

import (
	"net"
	"testing"
)

func TestRunMoveJ(t *testing.T) {
	var (
		jsonFile = "./format3_9.json"
		rCFormat RealtimeCommunicationsFormat
		err      error
	)
	conn, err := net.Dial("tcp", "192.168.1.107:30003")
	if rCFormat, err = GetRealtimeCommunicationsFormat(jsonFile); err != nil {
		t.Error(err)
	}

	cf := CommunicationsFloat64{
		Z: -0.06,
	}
	if err = RunURWithMoveJ(rCFormat, conn, cf); err != nil {
		t.Error(err)
	}
}
