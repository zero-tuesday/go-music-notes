/*
+build darwin
+build amd64
*/

package main

/*
#cgo CFLAGS: -I/usr/local/Cellar/aubio/0.4.9_1/include/
#cgo LDFLAGS: -L/usr/local/lib/ -laubio
#include <stdio.h>
#include <aubio/aubio.h>
*/
import "C"

import (
    "fmt"
    "os"
    "syscall"
)

func main() {
    hopSize := C.uint(1024)

    args := os.Args[1:]
    if len(args) < 1 {
        fmt.Println("No file provided")
        syscall.Exit(1)
    }

    filename := args[0]
    source := C.new_aubio_source(C.CString(filename), 0, hopSize)
    if source == nil {
        fmt.Println("File not found: ", filename)
        syscall.Exit(1)
    }

    fvec := C.new_fvec(hopSize)
    read := C.uint(0)

    for ok := true; ok; ok = (read == hopSize) {
        C.aubio_source_do(source, fvec, &read)
        C.fvec_print(fvec)
        fmt.Println("Read", read)
    }

    C.del_fvec(fvec)

    C.aubio_source_close(source)
    C.del_aubio_source(source)
    C.aubio_cleanup()
}
