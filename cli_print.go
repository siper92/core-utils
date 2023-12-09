package core_utils

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

func PrintError(err string, params ...interface{}) {
	errMessage := "### " + fmt.Sprintf(err, params...) + " ###"

	val := int(math.Min(float64(len(errMessage)), 120))
	ErrorMessage(strings.Repeat("#", val))
	ErrorMessage(errMessage)
	ErrorMessage(strings.Repeat("#", val))

	if !debug {
		return
	}

	ok := true
	var file string
	var lineN int
	callerIndex := 1
	for ok {
		_, file, lineN, ok = runtime.Caller(callerIndex)
		if ok {
			ErrorMessage(">> %s[%d]\n", file, lineN)
		}
		callerIndex++
	}
}

func DebugError(err error) {
	if err != nil {
		ImportantDebug(err.Error())
	}
}

func CliError(msg ...interface{}) {
	for _, m := range msg {
		switch m.(type) {
		case string:
			ErrorMessage(m.(string))
		case error:
			ErrorMessage(m.(error).Error())
		default:
			PrintError("...Unknown error type")
		}
	}
}

func ErrorMessage(msg string, p ...interface{}) {
	color.Red(fmt.Sprintf(msg, p...))
}

func PrintAndPanic(err interface{}) {
	if err == nil {
		return
	}

	if v, ok := err.(error); ok {
		PrintError(v.Error())
	}

	panic(err)
}

func StopOnError(err interface{}) {
	if err == nil {
		return
	}

	if v, ok := err.(error); ok {
		PrintError(v.Error())
	}

	os.Exit(1)
}

func StopAndNotify(msg string, p ...interface{}) {
	PrintError(msg, p...)
	os.Exit(1)
}

func Warning(f string, a ...interface{}) {
	color.Yellow(f, a...)
}

var (
	notice = false
	debug  = false
)

func AllowNotice() {
	notice = true
}

func AllowDebug() {
	debug = true
}

func Notice(f string, a ...interface{}) {
	if notice {
		color.Magenta(f, a...)
	}
}

func Debug(f string, a ...interface{}) {
	if IsDebugMode() {
		color.Yellow(f, a...)
	}
}

func ImportantDebug(f string, a ...interface{}) {
	if IsDebugMode() {
		color.HiBlue("###################################################")
		color.HiBlue(f, a...)
		color.HiBlue("###################################################")
	}
}

func HiGreen(f string, a ...interface{}) {
	color.HiGreen("###################################################")
	color.HiGreen(f, a...)
	color.HiGreen("###################################################")
}

func TestFailMessage(f string, a ...interface{}) string {
	return fmt.Sprintf("%s\n%s\n%s\n",
		color.RedString(">>>>>>>>>>>>>>>>>>>>>"),
		color.RedString(fmt.Sprintf(f, a...)),
		color.RedString(">>>>>>>>>>>>>>>>>>>>>"),
	)
}

func ErrorWarning(e error) {
	if e != nil {
		color.Yellow(e.Error())
	}
}

func PrintMessage(f string, a ...interface{}) {
	color.Green(" >> %s\n", fmt.Sprintf(f, a...))
}

func PrintServerResponse(f string, a ...interface{}) {
	color.Blue("  #[%s]# >> %s\n", time.Now().Format("2006/01/02-15:04:05"), fmt.Sprintf(f, a...))
}

func PrintWarningMessage(f string, a ...interface{}) {
	s := fmt.Sprintf(f, a...)
	line := fmt.Sprintf("######%s######\n", strings.Repeat("#", len(s)))
	Warning(line)
	Warning("##### %s #####\n", s)
	Warning(line)
}

func FormattedJson(v interface{}) ([]byte, error) {
	jsonVal, err := json.Marshal(v)
	if err != nil {
		return []byte{}, err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, jsonVal, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return prettyJSON.Bytes(), nil
}

func LogTimeElapsed(start time.Time, name string) {
	elapsed := time.Since(start)
	Notice("%s took %s", name, elapsed)
}

func IsDebugMode() bool {
	return debug || flag.Lookup("test.v") != nil
}
