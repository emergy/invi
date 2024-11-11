package template

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jmespath/go-jmespath"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func jsonStringToInterface(jsonStr string) (interface{}, error) {
	var data interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func interfaceToJSONString(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// jmesPathFunc выполняет выражение JMESPath над входными данными
func jmesPathFunc(expression string, data interface{}) (interface{}, error) {
	// Преобразуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Преобразуем JSON обратно в map[string]interface{}
	var mapData interface{}
	err = json.Unmarshal(jsonData, &mapData)
	if err != nil {
		return nil, err
	}

	// Выполняем поиск JMESPath
	result, err := jmespath.Search(expression, mapData)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// dict создает карту из переданных аргументов
func dict(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		panic("dict expects an even number of arguments")
	}
	m := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			panic("dict keys must be strings")
		}
		m[key] = values[i+1]
	}
	return m
}

// capitalizeFirst делает первую букву строки заглавной
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func ternary(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func toYAML(data interface{}) (string, error) {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

func fromYAML(yamlStr string) (interface{}, error) {
	var data interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func nindent(n int, s string) string {
	indent := strings.Repeat(" ", n)
	return indent + strings.ReplaceAll(s, "\n", "\n"+indent)
}

func defaultValue(def, val interface{}) interface{} {
	// Проверяем, является ли значение пустым, используя reflect
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		if v.Len() == 0 {
			return def
		}
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		if v.Len() == 0 {
			return def
		}
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return def
		}
	case reflect.Bool:
		if !v.Bool() {
			return def
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 0 {
			return def
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() == 0 {
			return def
		}
	case reflect.Float32, reflect.Float64:
		if v.Float() == 0.0 {
			return def
		}
	default:
		if !v.IsValid() {
			return def
		}
	}

	// Если ни одно из условий не сработало, возвращаем val
	return val
}
