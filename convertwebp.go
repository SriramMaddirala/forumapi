package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func convertimg(inputFile string) {
	cmd := exec.Command("/Users/ram/downloads/webpcon/bin/cwebp", inputFile, "-o", strings.Split(inputFile, ".")[0]+".webp")
	out, err := cmd.Output()
	fmt.Print(err)
	fmt.Print(out)
}
func convertgif(inputFile string) {
	cmd := exec.Command("/Users/ram/downloads/webpcon/bin/gif2webp", inputFile, "-o", strings.Split(inputFile, ".")[0]+".webp")
	out, err := cmd.Output()
	fmt.Print(err)
	fmt.Print(out)
}
