package urremoteController

import (
	"fmt"
	"net"
	"time"
)

func WaitMoveJ(rCFormat RealtimeCommunicationsFormat, conn net.Conn, targetPose []float64, timeout time.Duration) error {
	var (
		data        []byte
		err         error
		actualposeI interface{}
		done        chan bool
	)

	toolVectorActual := rCFormat["Tool vector actual"]
	go func() {
		time.Sleep(timeout)
		done <- true
	}()
	for {
		if data, err = read(conn, rCFormat, timeout); err != nil {
			return err
		}
		begin := toolVectorActual.BeginIndex
		end := toolVectorActual.BeginIndex + toolVectorActual.DataSize
		if len(data) < end {
			return fmt.Errorf("Error: data length: %d less than toolVectorActual end index: %d", len(data), end)
		}
		toolVectorActual.SetData(data[begin:end])

		if actualposeI, err = toolVectorActual.Output(); err != nil {
			return err
		}
		// 轉換型別
		switch actualposeI.(type) {
		case []float64:
			actualpose := actualposeI.([]float64)
			if len(actualpose) != toolVectorActual.NumberOfValues {
				return fmt.Errorf("Error: actualpose is not match toolVectorActual.NumberOfValues")
			}
			if equalFloat64s(actualpose, targetPose) {
				return nil
			}
		default:
			return fmt.Errorf("Error: target interface type is not a []float64")
		}
		select {
		case <-done:
			return fmt.Errorf("Error: Timeout %d sec", timeout/time.Second)
		default:
		}
	}

	// return nil
}

func equalFloat64s(source []float64, target []float64) bool {

	if len(source) != len(target) {
		return false
	}
	for i := 0; i < len(source); i++ {
		if !floatEquals(source[i], target[i]) {
			return false
		}
	}
	return true
}

var EPSILON float64 = 0.0001

func floatEquals(source, target float64) bool {
	return (source-target) < EPSILON && (target-source) < EPSILON
}
