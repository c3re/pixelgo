package main

import (
	"github.com/kaey/framebuffer"
	"log"
	"fmt"
)

func main() {
	fb, err := framebuffer.Init("/dev/fb0")
	if err != nil {
		log.Fatal(err)
	}
	x, y := fb.Size()
	fmt.Printf("Screensize: %dx%d\n", x, y)
	// red, green, blue, alpha
	for i := 100; i < 200; i++ {
		for j := 100; j < 200; j++ {
			fb.WritePixel(i, j, 0x00, 0x00, 0xFF, 0x00)
		}
	}
	fb.Close()
}
