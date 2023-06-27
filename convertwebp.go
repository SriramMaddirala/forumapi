package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func convertimg(inputFile string) {
	cmd := exec.Command("/Users/ram/downloads/webpcon/bin/cwebp", inputFile, "-o", "/Users/ram/Pictures"+strings.Split(inputFile, ".")[0]+".webp")
	_, err := cmd.Output()
	if err != nil {
		fmt.Print(err)
	}
}
func convertgif(inputFile string) {
	cmd := exec.Command("/Users/ram/downloads/webpcon/bin/gif2webp", inputFile, "-o", "/Users/ram/Pictures"+strings.Split(inputFile, ".")[0]+".webp")
	_, err := cmd.Output()
	if err != nil {
		fmt.Print(err)
	}
}
