package dry

// FirstArg returns the first passed argument,
// can be used to extract first result value
// from a function call to pass it on to functions like fmt.Printf
func FirstArg(args ...any) any {
	if len(args) == 0 {
		return nil
	}
	return args[0]
}
