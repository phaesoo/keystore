package store

import (
	"fmt"
)

var validStatuses = []string{"running", "stopped"}

type Master struct {
	Status string `json:"status"`
}

// NewMaster creates a store of master
func NewMaster(status string) (Master, error) {
	var master Master
	if err := validateStatus(status); err != nil {
		return master, err
	}
	return master, nil
}

func validateStatus(status string) error {
	for _, s := range validStatuses {
		if status == s {
			return nil
		}
	}
	return fmt.Errorf("status is not one of %s", validStatuses)
}

func (m *Master) IsRunning() bool {
	return m.Status == "running"
}

func (m *Master) SetStatus(status string) error {
	if err := validateStatus(status); err != nil {
		return err
	}

	m.Status = status
	return nil
}
