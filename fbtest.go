package main

import (
	"bufio"
	"fmt"
	"github.com/kaey/framebuffer"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var xsize, ysize int

type pixel struct {
	x, y, a, r, g, b int
}

func main() {
	c := make(chan pixel)
	cs := (chan<- pixel)(c)
	cr := (<-chan pixel)(c)
	go fb_mgmt(cr)
	//c <- pixel{i, j, 0, 0, 255, 0}
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("TCP")
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
		px6match, _ := regexp.MatchString("^PX [[:digit:]]+ [[:digit:]]+ [[:xdigit:]]{6}\r\n$", status)
		px8match, _ := regexp.MatchString("^PX [[:digit:]]+ [[:digit:]]+ [[:xdigit:]]{8}\r\n$", status)
		sizematch, _ := regexp.MatchString("^SIZE\r\n$", status)
		if px6match {
			data := strings.Split(status, " ")
			X, _ := strconv.Atoi(data[1])
			Y, _ := strconv.Atoi(data[2])
			R := data[3][0:2]
			G := data[3][2:4]
			B := data[3][4:6]
			R_int, _ := strconv.ParseInt(R, 16, 32)
			G_int, _ := strconv.ParseInt(G, 16, 32)
			B_int, _ := strconv.ParseInt(B, 16, 32)
			fmt.Println("PX %d %d %d %d %d", X, Y, R_int, G_int, B_int)
			c <- pixel{X, Y, 0, int(R_int), int(G_int), int(B_int)}

		} else if px8match {
			fmt.Println("PX mit 8 Farbwerten empfangen")
		} else if sizematch {
			fmt.Print("SIZE x y")
		}
	}
}
