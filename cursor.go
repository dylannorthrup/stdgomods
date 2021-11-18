package stdgomods

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

var (
	err         error                // Global error
	esc                     = "\x1b" // Escape character what for making ANSI escape sequences
	TERM_INITED bool        = false  // Only want to do init once
	oldState    *term.State          // So we can keep track of how things were before we got here
	stdin       *os.File             // File pointer for STDIN
	stdinFd     int                  // Integer for file descriptor of STDIN
	sw, sh      int                  // Screen width and height
	XPos, YPos  int                  // X and Y position of the cursor as reported by 'UpdateCursorLocation'
)

func tput(strFormat string, args ...interface{}) string {
	if !TERM_INITED {
		InitTerm()
	}
	return fmt.Sprintf("%s%s", esc, fmt.Sprintf(strFormat, args...))
}

// Works around race condition when using 'Check' from inside
// InitTerm (but after we start mucking with the terminal window)
func TermCheck(err error, msg string) {
	prevTI := TERM_INITED
	TERM_INITED = true
	Check(err, msg)
	TERM_INITED = prevTI
}

func InitTerm() {
	// We should only need to do this once.
	if TERM_INITED {
		return
	}
	stdin = os.Stdin
	stdinFd = int(stdin.Fd())
	oldState, err = term.MakeRaw(stdinFd)
	TermCheck(err, "Problems running term.MakeRaw()")

	// And, before we leave, let folks know we've already been here
	TERM_INITED = true
}

func RestoreTerm() {
	// If we didn't do any init, then no need to be here
	if !TERM_INITED {
		return
	}
	term.Restore(stdinFd, oldState)
}

// Methods with semi-friendly names
func setTermSize() {
	if !TERM_INITED {
		InitTerm()
	}
	sh, sw, err = term.GetSize(stdinFd)
	Check(err, "Error getting TermSize()")
}

func MoveCursorTo(x, y int) {
	fmt.Print(TPCursorPosition(x, y))
}

func ScreenWidth() int {
	return TPCols()
}
func ScreenHeight() int {
	return TPRows()
}
func ScreenSize() (x, y int) {
	y = TPCols()
	x = TPRows()
	return x, y
}

//// Currently busted because I don't have a reliable way to get the results back
//// from TPDeviceStatusReport().  *IF* I get those results back, it works great.
//// TODO: Get reliable way to read buffer
func UpdateCursorLocation() error {
	time.Sleep(time.Duration(20) * time.Millisecond)
	buf := make([]byte, 12)
	code := make([]byte, 12)
	fmt.Print(TPDeviceStatusReport())
	time.Sleep(time.Duration(20) * time.Millisecond)
	// Expected output is of the form `ESC [ 06 ; 12 R`
	// Should be at least 8 bytes worth of data....
	numRead := 0
	// Loop until we get the ending 'R', or we exceed
	// the loop counter.
	loopCounter := 10
	for {
		numRead, _ = stdin.Read(buf)
		// Uncomment if we want to see what's inside 'buf'
		// fmt.Printf(" === %v === \r\n", buf)
		// Copy the bytes from buf to a temporary slice
		tmp := buf[:numRead]
		// Append `buf` to `code` (we want to do this regardless of whether or not we have
		// hit our terminating condition)
		code = append(code, buf[:numRead]...)
		// Check the end of the `tmp` slice to see if it contains an 'R' (indicating the end
		// of the ANSI response). If so, escape the loop
		if bytes.Equal(tmp[:1], []byte("R")) {
			break
		}
		time.Sleep(time.Duration(20) * time.Millisecond)
		loopCounter -= 1
		if loopCounter <= 0 {
			return fmt.Errorf("Could not get output from stdin.Read")
		}
	}

	coords := fmt.Sprintf("%s", code[2:numRead-1])
	coordsAry := strings.Split(coords, ";")
	XPos, _ = strconv.Atoi(coordsAry[0])
	YPos, _ = strconv.Atoi(coordsAry[1])
	return nil
}
func GetCurrentCursorLocation() (int, int) {
	UpdateCursorLocation()
	return XPos, YPos
}
func CUp() {
	MoveCursorUp(1)
}
func CDown() {
	MoveCursorDown(1)
}
func CLeft() {
	MoveCursorLeft(1)
}
func CRight() {
	MoveCursorRight(1)
}
func MoveCursorUp(times int) {
	fmt.Print(TPCursorUp(times))
}
func MoveCursorDown(times int) {
	fmt.Print(TPCursorDown(times))
}
func MoveCursorRight(times int) {
	fmt.Print(TPCursorRight(times))
}
func MoveCursorLeft(times int) {
	fmt.Print(TPCursorLeft(times))
}
func SavePosition() {
	fmt.Print(TPSaveCursor())
}
func RestorePosition() {
	fmt.Print(TPRestoreCursor())
}
func ClearEntireLine() {
	fmt.Print(TPEraseLine())
}
func ClearToBeginningOfLine() {
	fmt.Print(TPEraseLineBackward())
}
func ClearToEndOfLine() {
	fmt.Print(TPEraseLineForward())
}
func ClearFromHereDown() {
	fmt.Print(TPClearToEndOfScreen())
}
func ClearFromHereUp() {
	fmt.Print(TPClearToTopOfScreen())
}
func ClearScreen() {
	fmt.Print(TPClearScreen())
}

// This returns a string for a progress bar that looks like so:
//    [###########=====]
// Colorization might be a thing... we'll see
// Params:
//  - size:  Number of discrete units in the bar
//  - total: Total items to be done
//  - done:  Number of items that ARE done
func ProgressBar(size, total, done int64) string {
	//// TODO: CURRENTLY not working as intended... Needs to be fixed
	// Do some simple bounds checking
	if size <= 2 || total < done || total < 1 || done < 0 {
		return "[#]"
	}
	// What's the pct done?
	pDone := IntPct(done, total)
	// How many "done units" does that percent represent?
	pctUnits := 100 / size
	numDoneUnits := int(pDone / pctUnits)
	numTBDUnits := int(size) - numDoneUnits
	pb := fmt.Sprintf("%s[%s%s]", TPEraseLineForward(), strings.Repeat("#", numDoneUnits), strings.Repeat("=", numTBDUnits))
	return pb
}

// These functions are roughly named after the arguments to the `tput` command.
// https://tldp.org/HOWTO/Bash-Prompt-HOWTO/x405.html for reference

func TPCursorPosition(x, y int) string {
	return tput("[%d;%dH", x, y)
}

func TPTermSize() (width, height int) {
	if sw == 0 && sh == 0 {
		setTermSize()
	}
	return sw, sh
}

func TPRows() int {
	if sw == 0 && sh == 0 {
		setTermSize()
	}
	return sh
}

func TPCols() int {
	if sw == 0 && sh == 0 {
		setTermSize()
	}
	return sh
}
func TPDeviceStatusReport() string {
	return tput("[6n")
}
func TPSaveCursor() string {
	return tput("7")
}
func TPRestoreCursor() string {
	return tput("8")
}
func TPEraseLine() string {
	return tput("[2K")
}
func TPEraseLineBackward() string {
	return tput("[1K")
}
func TPEraseLineForward() string {
	return tput("[0K")
}
func TPClearScreen() string {
	return tput("[2J")
}
func TPClearToTopOfScreen() string {
	return tput("[1J")
}
func TPClearToEndOfScreen() string {
	return tput("[0J")
}
func TPCursorUp(times int) string {
	return strings.Repeat(tput("[1A"), times)
}
func TPCursorDown(times int) string {
	return strings.Repeat(tput("[1B"), times)
}
func TPCursorRight(times int) string {
	return strings.Repeat(tput("[1C"), times)
}
func TPCursorLeft(times int) string {
	return strings.Repeat(tput("[1D"), times)
}
