package stdgomods

import (
	"fmt"
	"log"
	filepath "path/filepath"
	"regexp"
	"runtime"
	"strings"

	aec "github.com/morikuni/aec"
)

var (
	ES_DEBUG         bool // Enable PDebug output
	ES_DEBUG_TOGGLED bool // Track if we togged debug output or not in ToggleDebugIfNeeded()
)

// Print debug info. Uses logger with one-time config options set in Init()
// It will only print if the global ES_DEBUG variable is 'true'
func PDebug(msg ...string) {
	if !ES_DEBUG {
		return
	}
	// Combine the strings in `msg` into a single string
	fullMsg := strings.Join(msg, " ")

	// Use this to print info about where we are debugging from
	pc, absFileName, line, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("Could not get runtime info from inside PDebug(). Exiting")
	}
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	funcName := parts[pl-1]
	fileName := filepath.Base(absFileName)
	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
	}
	log.Printf("%s:%d _func %s()_ :\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n\t%s\n", fileName, line, funcName, fullMsg)
}

// Simple utility function to check the error status and, if it is not nil,
// log.Panic a message and exit
func PanicIfError(err error, msg string) {
	if err != nil {
		log.Panic("%s : %s\n", msg, err.Error())
	}
}

// This is used to allow for fine-grained DEBUG output with multiple flags.
// Caveats:
// - Do not overlap areas that flags are used. The intent is for flags to be specific to a function
// - This is intended used in conjunction with PDebug().
func ToggleDebugIfNeeded(flag bool) {
	// If we toggled ES_DEBUG, toggle it back and set ES_DEBUG_TOGGLED to false
	// Then return because our work here is done
	if ES_DEBUG_TOGGLED {
		ES_DEBUG = !ES_DEBUG
		ES_DEBUG_TOGGLED = false
		return
	}
	// If the flag is true and ES_DEBUG is false, we should enable ES_DEBUG but also
	// set the ES_DEBUG_TOGGLED flag so we know to disable it later
	if flag && !ES_DEBUG {
		ES_DEBUG = true
	}
	// If the flag is false and ES_DEBUG is true, we turn off ES_DEBUG for now and note
	// we need to turn it back on again later
	if !flag && ES_DEBUG {
		ES_DEBUG = false
		ES_DEBUG_TOGGLED = true
	}
}

func bold(msg string) string {
	return fmt.Sprintf(aec.Bold.Apply(msg))
}
func red(msg string) string {
	return fmt.Sprintf("%s", bold(aec.RedF.Apply(msg)))
}
func yellow(msg string) string {
	return fmt.Sprintf("%s", bold(aec.YellowF.Apply(msg)))
}
func blue(msg string) string {
	return fmt.Sprintf("%s", bold(aec.BlueF.Apply(msg)))
}
func cyan(msg string) string {
	return fmt.Sprintf("%s", bold(aec.CyanF.Apply(msg)))
}
func purple(msg string) string {
	return fmt.Sprintf("%s", bold(aec.MagentaF.Apply(msg)))
}
func white(msg string) string {
	return fmt.Sprintf("%s", bold(aec.WhiteF.Apply(msg)))
}
func orange(msg string) string {
	Orange := aec.Color8BitF(aec.NewRGB8Bit(0xac, 0x81, 0x00))
	return fmt.Sprintf("%s", Orange.Apply(msg))
}
