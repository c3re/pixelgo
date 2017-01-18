package main

import (
	"github.com/kaey/framebuffer"
	"log"
	"fmt"
	"net"
	"bufio"
	"regexp"
)

var xsize, ysize int

type pixel struct {
	x, y, a, r, g, b int
}


func main() {
	c := make(chan pixel)
	cs := (chan<- pixel) (c)
	cr := (<-chan pixel) (c)
	go fb_mgmt(cr)
	//c <- pixel{i, j, 0, 0, 255, 0}
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, cs)
	}
}

func fb_mgmt(c <-chan pixel) {
	fb, err := framebuffer.Init("/dev/fb0")
	if err != nil {
		log.Fatal(err)
	}
	xsize, ysize := fb.Size()
	fmt.Println(xsize, ysize)
	for {
		px := <-c
		fb.WritePixel(px.x, px.y, px.r, px.g, px.b, px.a)
	}
}

func handleConnection(conn net.Conn, c chan<- pixel) {
	var err error
	var status string
	fmt.Println("Neue Connection: ", conn.RemoteAddr().String())
	for err == nil {
		status, err = bufio.NewReader(conn).ReadString('\n')
		fmt.Print("received: ", status[:len(status) - 1])
		matched, _ := regexp.Match("insert regex here", []byte(status))
		if matched {
			fmt.Println(" good")
		} else {
			fmt.Println(" bad")
		}
	}
	fmt.Print("Connection closed.")
}
