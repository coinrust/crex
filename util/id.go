package util

import (
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.MachineID = machineID
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

// AmazonEC2MachineID retrieves the private IP address of the Amazon EC2 instance
// and returns its lower 16 bits.
// It works correctly on Docker as well.
func machineID() (uint16, error) {
	return 0, nil
}

func NextID() (uint64, error) {
	id, err := sf.NextID()
	return id, err
}
