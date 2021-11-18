package stdgomods

// Any task-specific things we want to insure are run regardless.
// NOTE: Does not include conditional things like "Do we want to
//       make the terminal Raw?" because we only want to do that
//       *IFF* we need to use the cursor ansi escape codes.
func init() {
	InitOutputTools()
}

// Things we want to make sure we run on exit. Ideally, whoever calls
// us will put a `defer stdgomods.Cleanup()` very early in their `main()`
// function.
func Cleanup() {
	// Tests to see if we modified the terminal and, if so, restores it to
	// how it was configured originally.
	RestoreTerm()
}
