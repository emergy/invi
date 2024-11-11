package template

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

func ProcessTemplates(data interface{}, t interface{}) (interface{}, error) {
	return processTemplatesReflect(reflect.ValueOf(data), t)
}

func processTemplatesReflect(v reflect.Value, t interface{}) (interface{}, error) {
	funcMap := template.FuncMap{
		"to_json":         interfaceToJSONString,
		"from_json":       jsonStringToInterface,
		"jmespath":        jmesPathFunc,
		"dict":            dict,
		"capitalizeFirst": capitalizeFirst,
		"toString":        toString,
		"contains":        strings.Contains,
		"ternary":         ternary,
		"isTerminal":      isTerminal,
		"to_yaml":         toYAML,
		"from_yaml":       fromYAML,
		"nindent":         nindent,
		"color":           Colorize,
		"default":         defaultValue,
	}

	if !v.IsValid() {
		return nil, nil
	}

	// Обработка указателей и интерфейсов путем разыменования
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		// Создаем новую карту для хранения обработанных значений
		newMap := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			processedVal, err := processTemplatesReflect(val, t)
			if err != nil {
				return nil, err
			}
			newMap.SetMapIndex(key, reflect.ValueOf(processedVal))
		}
		return newMap.Interface(), nil
	case reflect.Slice, reflect.Array:
		// Создаем новый срез для хранения обработанных элементов
		newSlice := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			processedElem, err := processTemplatesReflect(elem, t)
			if err != nil {
				return nil, err
			}
			newSlice.Index(i).Set(reflect.ValueOf(processedElem))
		}
		return newSlice.Interface(), nil
	case reflect.String:
		s := v.String()
		if strings.HasPrefix(s, "#yamltemplate") {
			// Извлекаем текст шаблона и применяем его
			templateText := strings.TrimPrefix(s, "#yamltemplate")
			templateText = strings.TrimLeft(templateText, "\n")

			tmpl, err := template.New("").Funcs(funcMap).Parse(templateText)
			if err != nil {
				return nil, fmt.Errorf("failed to parse template: %v", err)
			}
			var sb strings.Builder
			err = tmpl.Execute(&sb, t)
			if err != nil {
				return nil, fmt.Errorf("failed to execute template: %v", err)
			}
			resultStr := sb.String()

			// Попробуем распарсить результат как YAML
			var yamlResult interface{}
			err = yaml.Unmarshal([]byte(resultStr), &yamlResult)
			if err == nil {
				// Если удалось распарсить YAML, возвращаем результат
				return yamlResult, nil
			} else {
				// Если парсинг не удался, возвращаем строку
				return resultStr, nil
			}
		}
		if strings.HasPrefix(s, "#template") {
			// Извлекаем текст шаблона и применяем его
			templateText := strings.TrimPrefix(s, "#template")
			templateText = strings.TrimLeft(templateText, "\n")

			tmpl, err := template.New("").Funcs(funcMap).Parse(templateText)
			if err != nil {
				return nil, fmt.Errorf("failed to parse template: %v", err)
			}
			var sb strings.Builder
			err = tmpl.Execute(&sb, t)
			if err != nil {
				return nil, fmt.Errorf("failed to execute template: %v", err)
			}
			return sb.String(), nil
		}
		return s, nil
	default:
		// Возвращаем другие типы как есть
		return v.Interface(), nil
	}
}
