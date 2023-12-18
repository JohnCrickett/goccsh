package main

import (
	"fmt"
	"io"
	"os"
)

func Pwd(w *io.PipeWriter) {
	defer w.Close()
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, dir)
}
