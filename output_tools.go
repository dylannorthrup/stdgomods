package stdgomods

import (
	"fmt"
	"log"
	filepath "path/filepath"
	"runtime"
	"strings"
	"time"

	aec "github.com/morikuni/aec"
)

type statusFunction func() string

var (
	DEBUG           bool             // Enable PDebug output
	debug_toggled   bool             // Track if we togged debug output or not in ToggleDebugIfNeeded()
	StatusFunctions []statusFunction // List of functions that should be called by `UpdateStatus`
	updateTicker    *time.Ticker     // A ticker to periodically ping the `UpdateStatus` function
	stopUpdates     chan struct{}    // A channel to tell the goroutine calling the `UpdateStatus` function to stop
	updateInterval  int              // Number of seconds to wait between doing updates. Default it to 3.
	updatesSetUp    bool
	customColors    map[string]aec.ANSI // Map for holding our colors
)

func InitOutputTools() {
	// Instantiating a map for my custom colors
	customColors = make(map[string]aec.ANSI)

	///// For PDebug
	// Set logger params... disable flags and add easily deliniated prefix
	log.SetFlags(0)
	log.SetPrefix("=============================================\nES_DEBUG: ")
}

// Print debug info. Uses logger with one-time config options set in InitOutputTools()
// It will only print if the global DEBUG variable is 'true'
func PDebug(msg ...string) {
	if !DEBUG {
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
	log.Printf("\r%s:%d _func %s()_ :\r\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\r\n\t%s\r\n", fileName, line, funcName, fullMsg)
}

func FmtPDebug(frmt string, vars ...interface{}) {
	PDebug(fmt.Sprintf("%s", fmt.Sprintf(frmt, vars...)))
}

// Simple utility function to check the error status and, if it is not nil,
// log.Panic a message and exit
func PanicIfError(err error, msg string) {
	if err != nil {
		log.Panic(fmt.Sprintf("%s : %s\n", msg, err.Error()))
	}
}

func Check(err error, msg string) {
	if err != nil {
		if TERM_INITED {
			RestoreTerm()
		}
		log.Panic(fmt.Sprintf("%s : %s\n", msg, err.Error()))
	}
}

// This is used to allow for fine-grained DEBUG output with multiple flags.
// Caveats:
// - Do not overlap areas that flags are used. The intent is for flags to be specific to a function
// - This is intended used in conjunction with PDebug().
func ToggleDebugIfNeeded(flag bool) {
	// If we toggled DEBUG, toggle it back and set debug_toggled to false
	// Then return because our work here is done
	if debug_toggled {
		DEBUG = !DEBUG
		debug_toggled = false
		return
	}
	// If the flag is true and DEBUG is false, we should enable DEBUG but also
	// set the debug_toggled flag so we know to disable it later
	if flag && !DEBUG {
		DEBUG = true
	}
	// If the flag is false and DEBUG is true, we turn off DEBUG for now and note
	// we need to turn it back on again later
	if !flag && DEBUG {
		DEBUG = false
		debug_toggled = true
	}
}

//// Convert Bytes to String and String to Bytes

func B2S(s string) []byte {
	b := []byte(s)
	return b
}

func S2B(b []byte) string {
	s := string(b)
	return s
}

//// Colorization functions
func Colorize(frmt string, color aec.ANSI, vars ...interface{}) string {
	return color.Apply(fmt.Sprintf(frmt, vars...))
}
func PColorize(frmt string, color aec.ANSI, vars ...interface{}) {
	fmt.Printf("%s", color.Apply(fmt.Sprintf(frmt, vars...)))
}
func SPColorize(frmt string, color aec.ANSI, vars ...interface{}) string {
	return fmt.Sprintf("%s", color.Apply(fmt.Sprintf(frmt, vars...)))
}
func Bold(msg string) string {
	return fmt.Sprintf(aec.Bold.Apply(msg))
}
func Red(msg string) string {
	return Colorize("%s", aec.Bold, Colorize("%s", aec.RedF, msg))
}
func PRed(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Red(msg))
}
func SPRed(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Red(msg))
}
func Green(msg string) string {
	return Colorize("%s", aec.Bold, Colorize("%s", aec.GreenF, msg))
}
func PGreen(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Green(msg))
}
func SPGreen(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Green(msg))
}
func Yellow(msg string) string {
	return Colorize("%s", aec.Bold, Colorize("%s", aec.YellowF, msg))
}
func PYellow(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Yellow(msg))
}
func SPYellow(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Yellow(msg))
}

// This is a lighter shade of green than the normal green. It's defined
// because my normal screen colors are green on black.
func Grean(msg string) string {
	if customColors["GreanF"] == nil {
		customColors["GreanF"] = aec.Color8BitF(aec.NewRGB8Bit(0x4c, 0xfc, 0x4c))
	}
	return Colorize("%s", aec.Bold, Colorize("%s", customColors["GreanF"], msg))
}
func PGrean(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Grean(msg))
}
func SPGrean(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Grean(msg))
}
func Blue(msg string) string {
	return Colorize("%s", aec.Bold, Colorize("%s", aec.BlueF, msg))
}
func PBlue(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Blue(msg))
}
func SPBlue(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Blue(msg))
}
func Cyan(msg string) string {
	return Colorize("%s", aec.Bold, Colorize("%s", aec.CyanF, msg))
}
func PCyan(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Cyan(msg))
}
func SPCyan(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Cyan(msg))
}
func Purple(msg string) string {
	return Colorize("%s", aec.Bold, Colorize("%s", aec.MagentaF, msg))
}
func PPurple(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Purple(msg))
}
func SPPurple(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Purple(msg))
}
func White(msg string) string {
	return Colorize("%s", aec.Bold, Colorize("%s", aec.WhiteF, msg))
}
func PWhite(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", White(msg))
}
func SPWhite(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", White(msg))
}
func Orange(msg string) string {
	if customColors["OrangeF"] == nil {
		customColors["OrangeF"] = aec.Color8BitF(aec.NewRGB8Bit(0xac, 0x81, 0x00))
	}
	return Colorize("%s", aec.Bold, Colorize("%s", customColors["OrangeF"], msg))
}
func POrange(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Orange(msg))
}
func SPOrange(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Orange(msg))
}
func Pink(msg string) string {
	if customColors["PinkF"] == nil {
		customColors["PinkF"] = aec.Color8BitF(aec.NewRGB8Bit(0xff, 0x14, 0x00))
	}
	return Colorize("%s", aec.Bold, Colorize("%s", customColors["PinkF"], msg))
}
func PPink(frmt string, vars ...interface{}) {
	msg := fmt.Sprintf(frmt, vars...)
	fmt.Printf("%s", Pink(msg))
}
func SPPink(frmt string, vars ...interface{}) string {
	msg := fmt.Sprintf(frmt, vars...)
	return fmt.Sprintf("%s", Pink(msg))
}

func setUpUpdateStatus() {
	updateInterval = 1
	updateTicker = time.NewTicker(time.Duration(updateInterval) * time.Second)
	stopUpdates = make(chan struct{})
	// Set up background function to do status updates
	go func() {
		for {
			select {
			case <-updateTicker.C:
				UpdateStatus()
			case <-stopUpdates:
				updateTicker.Stop()
				return
			}
		}
	}()
}

// Use our cursor movement abilities to make a small status section in the upper right corner
// to print out whatever we thing is important up there.
func UpdateStatus() {
	// If we don't have any status functions, move along.
	if StatusFunctions == nil {
		return
	}
	// Make sure we have a proper screen width set and, if not, set one
	if sw == 0 {
		setTermSize()
	}
	// Gather up the output from all our functions
	output := ""
	for _, f := range StatusFunctions {
		output = fmt.Sprintf("%s\r%s\n", output, f())
	}
	now := fmt.Sprintf("%s", time.Now().Format(time.RubyDate))
	offset := ScreenWidth() - len(now) - 5
	output = fmt.Sprintf("%s\r%s%s\n", output, TPCursorRight(offset), now)
	// Save where we were
	UpdateCursorLocation()
	prevX := XPos
	prevY := YPos
	// Move to upper right corner
	MoveCursorTo(1, sw)
	// Print our output
	fmt.Print(output)
	// Return to our original location
	MoveCursorTo(prevX, prevY)
}

func AddStatusFunction(f statusFunction) {
	if !updatesSetUp {
		setUpUpdateStatus()
	}
	StatusFunctions = append(StatusFunctions, f)
}
