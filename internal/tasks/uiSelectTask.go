package tasks

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UiSelect struct {
	ItemText          string      `json:"item_text"`
	SecondaryItemText string      `json:"secondary_item_text"`
	DetailsText       string      `json:"details_text"`
	Data              interface{} `json:"data"`
	Register          string      `json:"register"`
}

func RunUiSelect(task UiSelect, tmplInputData map[string]interface{}) (Register, error) {
	app := tview.NewApplication()

	list := tview.NewList()
	textView := tview.NewTextView().SetDynamicColors(true).SetWrap(true)

	// Циклический переход по элементам списка
	list.SetWrapAround(true)

	// Действие по нажатию Esc
	// list.SetDoneFunc(func())

	// SetChangedFunc задает функцию, которая вызывается при переходе пользователя к элементу списка
	// func (l *List) SetChangedFunc(handler func(index int, mainText string, secondaryText string, shortcut rune)) *List
	// list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
	// 	// При переходе к элементу списка, отображаем его детальную информацию
	// 	detailedInfo := formatHostDetails(hosts[index])
	// 	textView.SetText(detailedInfo)
	// })

	// Configure the main Flex layout
	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(list, 0, 1, true).
		AddItem(textView, 0, 2, false)

	backgroundColor := tview.Styles.PrimitiveBackgroundColor
	r, g, b := backgroundColor.RGB()
	brighterColor := tcell.NewRGBColor(
		int32(min(int(r)+90, 255)),
		int32(min(int(g)+90, 255)),
		int32(min(int(b)+90, 255)),
	)
	list.SetSecondaryTextColor(brighterColor)

	var selectedItem interface{}

	for index, itemRaw := range task.Data.([]interface{}) {
		shortcut := func(index int) rune {
			if index < 9 {
				return rune(index + 49)
			} else if index < 35 {
				return rune(index + 88)
			} else {
				return 0
			}
		}

		item := itemRaw.(map[string]interface{})

		list.AddItem(item["name"].(string), item["ip"].(string), shortcut(index), func() {
			selectedItem = item
			app.Stop()
		})
	}

	if err := app.SetRoot(flex, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

	register := Register{
		Name: task.Register,
		Data: selectedItem,
	}

	return register, nil
}
