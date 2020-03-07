package urRemoteController

import (
	"net"
	"testing"
	"time"
)

func TestRunMoveJ(t *testing.T) {
	var (
		jsonFile   = "./format3_9.json"
		rCFormat   RealtimeCommunicationsFormat
		err        error
		targetPose []float64
		timeout    = time.Second * 5
	)
	conn, err := net.Dial("tcp", "192.168.1.107:30003")
	if rCFormat, err = GetRealtimeCommunicationsFormat(jsonFile); err != nil {
		t.Error(err)
	}
	defer conn.Close()
	cf := CommunicationsFloat64{
		Z: -0.06,
	}
	if targetPose, err = RunURWithMoveJ(rCFormat, conn, cf, timeout); err != nil {
		t.Error(err)
	}

	beforeTime := time.Now()
	if err = WaitMoveJ(rCFormat, conn, targetPose, timeout); err != nil {
		t.Error(err)
	}
	afterTime := time.Now()
	if afterTime.Sub(beforeTime) >= timeout {
		t.Error("Timeout")
	}

}
