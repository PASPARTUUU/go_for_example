package errpath

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	standart   modeType = "standart"
	functional modeType = "functional"
	extended   modeType = "extended"
	custom     modeType = "custom"
)

const currentFuncCaller int = 2

type modeType string
type modeStruct struct {
	Standart   modeType
	Functional modeType
	Extended   modeType
}

// Mods - список возможных отображений ошибки
var Mods modeStruct
var setMod modeType = custom

// SetMod - установить модификатор отображения ошибки
func SetMod(mod modeType) {
	setMod = mod
}

func init() {
	Mods.Standart = standart
	Mods.Functional = functional
	Mods.Extended = extended
}

// SuperError -
type SuperError struct {
	Err     error
	Message string
	Path    path
}
type path struct {
	FuncName string
	File     string
	Line     int
}

func (e SuperError) Error() string {
	if e.Err == nil {
		return ""
	}
	if setMod == standart {
		return fmt.Sprintf("%s:%s", e.Message, e.Err)
	}
	if setMod == extended {
		return fmt.Sprintf("(%s|%s|%v)%s -> %s", e.Path.FuncName, e.Path.File, e.Path.Line, e.Message, e.Err.Error())
	}
	if setMod == functional {
		return fmt.Sprintf("(%s)%s -> %s", e.Path.FuncName, e.Message, e.Err.Error())
	}

	if setMod == custom {
		return fmt.Sprintf("(%s|%v)%s -> %s", e.Path.FuncName, e.Path.Line, e.Message, e.Err.Error())
	}

	return fmt.Sprintf("%s", e.Err.Error()) // дефолтная реализация
}

// Infof -
func Infof(format string, a ...interface{}) string {
	path := getPath(currentFuncCaller)
	return fmt.Sprintf("(%s|%v)%s", path.FuncName, path.Line, fmt.Sprintf(format, a...))
}

// InfofWithFuncCaller - с номером функции в стэке
func InfofWithFuncCaller(caller int, format string, a ...interface{}) string {
	path := getPath(caller)
	return fmt.Sprintf("(%s|%v)%s", path.FuncName, path.Line, fmt.Sprintf(format, a...))
}

// Errorf - новая ошибка
func Errorf(format string, a ...interface{}) *SuperError {
	var sr SuperError
	sr.Err = fmt.Errorf("end")
	sr.Message = fmt.Sprintf(format, a...)
	sr.Path = getPath(currentFuncCaller)
	return &sr
}

// Err - описание ошибки
func Err(err error, text ...interface{}) *SuperError {
	var sr SuperError
	if err == nil {
		return nil
	}
	sr.Err = err
	sr.Message = fmt.Sprint(text...)
	sr.Path = getPath(currentFuncCaller)
	return &sr
}

// Func - возвращает имя функции
func Func() string {
	pp, _, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("not ok")
	}
	return funcNameParse(runtime.FuncForPC(pp).Name())
}

// FuncTrace - возвращает порядок вызовов функций
func FuncTrace() []string {
	var res []string
	pc := make([]uintptr, 1022) // at least 1 entry needed
	m := runtime.Callers(currentFuncCaller, pc)
	pc = pc[:m]

	for i, item := range pc {
		f := runtime.FuncForPC(item)
		res = append(res, fmt.Sprint(funcNameParse(f.Name()), "\t| ", i, "\n"))
	}

	return res
}

// Funcc - возвращает имя функции в скобках
func Funcc() string {
	pp, _, _, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("not ok")
	}
	return "(" + funcNameParse(runtime.FuncForPC(pp).Name()) + ")"
}

// DontPanic - парсинг введеного стека при панике для удобного чтения
func DontPanic(panic string) {
	var color int = 0
	firstSplit := strings.Split(panic, `\n\t`)
	for _, fs := range firstSplit {
		secondSplit := strings.Split(fs, `\n`)
		fmt.Println()
		color = 0
		for _, item := range secondSplit {
			if color == 0 {
				fmt.Printf("\x1b[33m%s\x1b[0m", string(item))
			}
			if color == 1 {
				fmt.Printf("\x1b[36m%s\x1b[0m", string(item))
			}
			color++
		}
	}
}

func getPath(caller int) path {
	pp, file, line, ok := runtime.Caller(caller)
	if !ok {
		fmt.Println("not ok")
	}
	fn := runtime.FuncForPC(pp).Name()
	return path{FuncName: funcNameParse(fn), File: fileNameParse(file), Line: line}
}

func funcNameParse(name string) (res string) {
	sp := strings.Split(name, ".")
	res = sp[len(sp)-1]
	return
}

func fileNameParse(name string) (res string) {
	sp := strings.Split(name, "/")
	res = sp[len(sp)-1]
	return
}
