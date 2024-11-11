package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type HttpTask struct {
	Url      string            `json:"url"`
	Method   string            `json:"method"`
	Handlers map[string]string `json:"handlers"`
	Data     interface{}       `json:"data"`
	Register string            `json:"register"`
}

func RunHttpTask(ctx map[string]interface{}, task HttpTask) (Register, error) {
	req := resty.New().R()

	for k, v := range task.Handlers {
		req.SetHeader(k, v)
	}

	var resp *resty.Response
	var err error

	switch task.Method {
	case "GET":
		resp, err = req.Get(task.Url)
	case "POST":
		if task.Data != nil {
			req.SetBody(task.Data)
		}
		resp, err = req.Post(task.Url)
	case "PUT":
		if task.Data != nil {
			req.SetBody(task.Data)
		}
		resp, err = req.Put(task.Url)
	case "DELETE":
		resp, err = req.Delete(task.Url)
	case "PATCH":
		if task.Data != nil {
			req.SetBody(task.Data)
		}
		resp, err = req.Patch(task.Url)
	case "HEAD":
		resp, err = req.Head(task.Url)
	case "OPTIONS":
		resp, err = req.Options(task.Url)
	default:
		return Register{}, fmt.Errorf("unknown method: %s", task.Method)
	}

	if err != nil {
		return Register{}, err
	}

	body := resp.Body()
	if body == nil {
		return Register{}, fmt.Errorf("response body is nil")
	}

	register := Register{
		Name: task.Register,
	}

	var jsonResult interface{}
	if err := json.Unmarshal(body, &jsonResult); err != nil {
		jsonResult = string(body)
		register.Error = fmt.Errorf("error unmarshalling response body: %v", err)
	}

	register.Data = jsonResult

	return register, nil
}
