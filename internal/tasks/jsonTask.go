package tasks

import (
	"encoding/json"
	"fmt"
)

type JsonTask struct {
	Data     interface{} `json:"data"`
	Register string      `json:"register"`
}

func RunJsonTask(registers []Register, task JsonTask) (Register, error) {
	var data interface{}

	if task.Data != nil {
		data = task.Data
	} else {
		data = registers
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return Register{}, fmt.Errorf("error marshalling json: %v", err)
	}

	register := Register{
		Name: task.Register,
	}

	register.Data = string(jsonBytes)

	return register, nil
}
