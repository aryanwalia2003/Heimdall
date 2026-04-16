package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/aryanwalia/heimdall/pkg/wrap"
)

//export ProcessError
func ProcessError(input *C.char) *C.char {
	goInput := C.GoString(input)
	goOutput := wrap.ProcessError(goInput)
	return C.CString(goOutput)
}

//export FreeString
func FreeString(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

func main() {}
