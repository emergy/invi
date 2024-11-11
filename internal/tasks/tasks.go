package tasks

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/emergy/invi/internal/template"
	"github.com/mitchellh/mapstructure"
)

type Register struct {
	Name  string      `json:"name"`
	Data  interface{} `json:"data"`
	Error error       `json:"error"`
}

func RunTasks(ctx map[string]interface{}) error {
	var registers []Register

	for _, taskRaw := range ctx["tasks"].([]interface{}) {
		tmplInputData := map[string]interface{}{
			"query": ctx["args"].([]string)[0],
			"flags": ctx["flags"],
		}

		for _, register := range registers {
			tmplInputData[register.Name] = register.Data
		}

		taskRaw, err := template.ProcessTemplates(
			taskRaw.(map[string]interface{}),
			tmplInputData,
		)
		if err != nil {
			return fmt.Errorf("error processing templates: %v", err)
		}

		task := taskRaw.(map[string]interface{})

		if task["type"] == nil {
			return fmt.Errorf("task type is required for task: %v", task)
		}

		if task["register"] == nil {
			task["register"] = "result"
		}

		switch task["type"].(string) {
		case "http":
			var httpTask HttpTask

			if err := mapstructure.Decode(task, &httpTask); err != nil {
				return fmt.Errorf("error decoding http task: %v", err)
			}

			register, err := RunHttpTask(ctx, httpTask)
			if err != nil {
				return err
			}

			registers = append(registers, register)

		case "json":
			var jsonTask JsonTask

			if err := mapstructure.Decode(task, &jsonTask); err != nil {
				return fmt.Errorf("error decoding json task: %v", err)
			}

			register, err := RunJsonTask(registers, jsonTask)
			if err != nil {
				return err
			}

			registers = append(registers, register)

		case "dump":
			spew.Dump(registers)

		case "show":
			if task["data"] != nil {
				fmt.Println(task["data"])
			} else {
				data := registers[len(registers)-1].Data
				fmt.Println(data)
			}

		case "ui_select":
			var uiSelect UiSelect

			if err := mapstructure.Decode(task, &uiSelect); err != nil {
				return fmt.Errorf("error decoding ui task: %v", err)
			}

			register, err := RunUiSelect(uiSelect, tmplInputData)
			if err != nil {
				return err
			}

			registers = append(registers, register)

		case "command":
			var commandTask CommandTask

			if err := mapstructure.Decode(task, &commandTask); err != nil {
				return fmt.Errorf("error decoding command task: %v", err)
			}

			register, err := RunCommandTask(commandTask)
			if err != nil {
				return err
			}

			registers = append(registers, register)

		default:
			return fmt.Errorf("unknown task type: %s", task["type"])
		}
	}
	return nil
}
