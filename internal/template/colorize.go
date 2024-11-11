package template

import (
	"fmt"
	"strings"
)

// Список наиболее используемых цветов и стилей
var colors = map[string]string{
	// Обычные цвета
	"black":   "30",
	"red":     "31",
	"green":   "32",
	"yellow":  "33",
	"blue":    "34",
	"magenta": "35",
	"cyan":    "36",
	"white":   "37",

	// Яркие (светлые) цвета
	"bright_black":   "90",
	"bright_red":     "91",
	"bright_green":   "92",
	"bright_yellow":  "93",
	"bright_blue":    "94",
	"bright_magenta": "95",
	"bright_cyan":    "96",
	"bright_white":   "97",

	// Стили
	"bold":      "1",
	"underline": "4",
	"reversed":  "7",

	// Сброс стиля
	"reset": "0",
}

// Функция для добавления ANSI-кодов цвета и стиля к тексту
func Colorize(text string, codes ...string) string {
	var ansiCodes []string
	for _, code := range codes {
		if val, exists := colors[code]; exists {
			ansiCodes = append(ansiCodes, val)
		} else {
			// Если код не найден в colors, считаем, что это ANSI код напрямую
			ansiCodes = append(ansiCodes, code)
		}
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", strings.Join(ansiCodes, ";"), text)
}
