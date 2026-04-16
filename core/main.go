package main

import "C"
import "github.com/aryanwalia/heimdall/pkg/wrap"

//export ProcessError
func ProcessError(input *C.char) *C.char {
	goInput := C.GoString(input)
	goOutput := wrap.ProcessError(goInput)
	return C.CString(goOutput)
}

// Placeholder for shared library logic
func main() {}
