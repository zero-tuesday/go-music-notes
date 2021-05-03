/*
+build darwin
+build amd64
*/

package main

/*
#cgo CFLAGS: -I/usr/local/Cellar/aubio/0.4.9_1/include/
#cgo LDFLAGS: -L/usr/local/lib/ -laubio
#include <stdlib.h>
#include <stdio.h>
#include <aubio/aubio.h>
*/
import "C"

import (
    "fmt"
    "os"
    "syscall"
    "unsafe"
)

func noteData(array *C.smpl_t) []C.smpl_t {
    return (*[1 << 28]C.smpl_t)(unsafe.Pointer(array))[:3:3]
}

func main() {
    defer C.aubio_cleanup()

    hopSize := C.uint(256)
    bufferSize := C.uint(512)
    sampleRate := C.uint(0)

    args := os.Args[1:]
    if len(args) < 1 {
        fmt.Println("No file provided")
        syscall.Exit(1)
    }

    filename := args[0]
    var c_fn *C.char= C.CString(filename)
    defer C.free(unsafe.Pointer(c_fn))
    source := C.new_aubio_source(c_fn, sampleRate, hopSize)
    if source == nil {
        fmt.Println("File not found:", filename)
        syscall.Exit(1)
    }
    defer C.aubio_source_close(source)
    defer C.del_aubio_source(source)

    if sampleRate == 0 {
        sampleRate = C.aubio_source_get_samplerate(source)
    }

    fmt.Println("Sample Rate:", sampleRate)

    var c_mthd *C.char = C.CString("default")
    defer C.free(unsafe.Pointer(c_mthd))
    var notes = C.new_aubio_notes(c_mthd, bufferSize, hopSize, sampleRate)
    if notes == nil {
        syscall.Exit(1)
    }
    defer C.del_aubio_notes(notes)

    C.aubio_notes_set_minioi_ms(notes, 0.030)
    C.aubio_notes_set_silence(notes, -90.0)
    C.aubio_notes_set_release_drop(notes, 10)

    var in *C.fvec_t = C.new_fvec(hopSize)
    defer C.del_fvec(in)
    var out *C.fvec_t = C.new_fvec(hopSize)
    defer C.del_fvec(out)
    read := C.uint(0)
    for ok := true; ok; ok = read == hopSize {
        C.aubio_source_do(source, in, &read) // reads data
        C.aubio_notes_do(notes, in, out) // analyze data for notes/pitch
        slice := noteData(out.data)
        if slice[0] > 0 {
            fmt.Println("midi:", slice[0], "old:", slice[2])
        }
    }

}
