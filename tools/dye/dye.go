package dye

import (
	"fmt"
	"runtime"
	"strings"
)

// Dye - строка которая будет окрашена
type Dye string

// New - приготовить строку к окраске
func New(text ...interface{}) Dye {
	var prep string
	for i, item := range text {
		prep += fmt.Sprint(item) + " "
		if i != len(text)-1 {
			prep += "| "
		}
	}
	return Dye(prep)
}

// Text - возвращает окрашенную строку
func (text Dye) Text() string { return string(text) }

// Print -
func (text Dye) Print() { fmt.Println(text.Text()) }

// Next -
func Next(text ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("not ok")
	}
	var s string = ""
	for _, item := range text {
		s += " " + fmt.Sprint(item) + " |"
	}
	New(fmt.Sprint("(", fileNameParse(file), "|", line, ")") + s).TextPurple().FontBolt().Print()
}

// Шрифты ---------------------------
// 0	нормальный режим
// 1	жирный
// 3	курсив (возможно не работает)
// 4	подчеркнутый
// 5	мигающий
// 7	инвертированные цвета
// 8	невидимый
// 9	зачеркнутый (возможно не работает)

// FontBolt - жирный текст
func (text Dye) FontBolt() Dye { return Dye(fmt.Sprintf("\x1b[1m%s\x1b[0m", string(text))) }

// FontUnderlined - подчеркнутый текст
func (text Dye) FontUnderlined() Dye { return Dye(fmt.Sprintf("\x1b[4m%s\x1b[0m", string(text))) }

// FontBlink - мигаюший текст
func (text Dye) FontBlink() Dye { return Dye(fmt.Sprintf("\x1b[5m%s\x1b[0m", string(text))) }

// Цвета текста ------------
// 30	черный
// 31	красный
// 32	зеленый
// 33	желтый
// 34	синий
// 35	пурпурный
// 36	голубой
// 37	белый

// TextBlack - черный цвет текста
func (text Dye) TextBlack() Dye { return Dye(fmt.Sprintf("\x1b[30m%s\x1b[0m", string(text))) }

// TextRed - красный цвет текста
func (text Dye) TextRed() Dye { return Dye(fmt.Sprintf("\x1b[31m%s\x1b[0m", string(text))) }

// TextGreen - зеленый цвет текста
func (text Dye) TextGreen() Dye { return Dye(fmt.Sprintf("\x1b[32m%s\x1b[0m", string(text))) }

// TextYellow - желтый цвет текста
func (text Dye) TextYellow() Dye { return Dye(fmt.Sprintf("\x1b[33m%s\x1b[0m", string(text))) }

// TextBlue - синий цвет текста
func (text Dye) TextBlue() Dye { return Dye(fmt.Sprintf("\x1b[34m%s\x1b[0m", string(text))) }

// TextPurple - пурпурный цвет текста
func (text Dye) TextPurple() Dye { return Dye(fmt.Sprintf("\x1b[35m%s\x1b[0m", string(text))) }

// TextCyan - голубой цвет текста
func (text Dye) TextCyan() Dye { return Dye(fmt.Sprintf("\x1b[36m%s\x1b[0m", string(text))) }

// TextWhite - белый цвет текста
func (text Dye) TextWhite() Dye { return Dye(fmt.Sprintf("\x1b[37m%s\x1b[0m", string(text))) }

// Цвета фона ------------
// 40	черный
// 41	красный
// 42	зеленый
// 43	желтый
// 44	синий
// 45	пурпурный
// 46	голубой
// 47	белый

// BackBlack - черный цвет фона
func (text Dye) BackBlack() Dye { return Dye(fmt.Sprintf("\x1b[40m%s\x1b[0m", string(text))) }

// BackRed - красный цвет фона
func (text Dye) BackRed() Dye { return Dye(fmt.Sprintf("\x1b[41m%s\x1b[0m", string(text))) }

// BackGreen - зеленый цвет фона
func (text Dye) BackGreen() Dye { return Dye(fmt.Sprintf("\x1b[42m%s\x1b[0m", string(text))) }

// BackYellow - желтый цвет фона
func (text Dye) BackYellow() Dye { return Dye(fmt.Sprintf("\x1b[43m%s\x1b[0m", string(text))) }

// BackBlue - синий цвет фона
func (text Dye) BackBlue() Dye { return Dye(fmt.Sprintf("\x1b[44m%s\x1b[0m", string(text))) }

// BackPurple - пурпурный цвет фона
func (text Dye) BackPurple() Dye { return Dye(fmt.Sprintf("\x1b[45m%s\x1b[0m", string(text))) }

// BackCyan - голубой цвет фона
func (text Dye) BackCyan() Dye { return Dye(fmt.Sprintf("\x1b[46m%s\x1b[0m", string(text))) }

// BackWhite - белый цвет фона
func (text Dye) BackWhite() Dye { return Dye(fmt.Sprintf("\x1b[47m%s\x1b[0m", string(text))) }

// ===================================

func fileNameParse(name string) (res string) {
	sp := strings.Split(name, "/")
	res = sp[len(sp)-1]
	return
}
