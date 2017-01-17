package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// von Wikipedia geklaut: https://de.wikipedia.org/wiki/Framebuffer#Linux-Framebuffer
var width, height, bitspp, bytespp int
var data byte[]
func main() {
	//	var row, col, width, height, bitspp, bytespp int
	// Framebuffer oeffnen
	fd, err := unix.Open("/dev/fb0", unix.O_RDWR, unix.S_IRWXU)
	if err != nil {
		log.Fatal(err)
	}
	var screeninfo fb_var_screeninfo
	ioctl(uintptr(fd), FBIOGET_VSCREENINFO, &screeninfo)
	fmt.Printf("Bildschirmauflösung: %dx%d\n", screeninfo.xres_virtual, screeninfo.yres_virtual)
	bitspp = screeninfo.bits_per_pixel
	width = screeninfo.xres
	height = screeninfo.yres
	bytespp = bitspp / 8

	if bitspp != 32 {
		fmt.Printf("Farbaufloesung = %d Bits pro Pixel\n", bitspp)
		fmt.Printf("Bitte aendern Sie die Farbtiefe auf 32 Bits pro Pixel\n")
		unix.Close(fd)
		os.Exit(1)
	}
	data, err := unix.Mmap(fd, 0, int(width*height*bytespp), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		log.Fatal(err)
	}

}

// copied from here: https://raw.githubusercontent.com/tianon/debian-golang-pty/master/ioctl.go
func ioctl(fd, cmd uintptr, ptr *fb_var_screeninfo) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, uintptr(unsafe.Pointer(ptr)))
	if e != 0 {
		return e
	}
	return nil
}

func draw(x, y int, a, r, g, b byte) {
	data[width*y+x] = a
	data[width*y+x+1] = r
	data[width*y+x+2] = g
	data[width*y+x+3] = b
}

const FBIOGET_VSCREENINFO uintptr = 0x4600

type pixel struct {
	a, r, g, b byte
}

// copied from here: https://www.kernel.org/doc/Documentation/fb/api.txt
type fb_bitfield struct {
	offset    uint32 /* beginning of bitfield	*/
	length    uint32 /* length of bitfield		*/
	msb_right uint32 /* != 0 : Most significant bit is */
}

// copied from here: https://www.kernel.org/doc/Documentation/fb/api.txt
type fb_var_screeninfo struct {
	xres           uint32 /* visible resolution		*/
	yres           uint32
	xres_virtual   uint32 /* virtual resolution		*/
	yres_virtual   uint32
	xoffset        uint32      /* offset from virtual to visible */
	yoffset        uint32      /* resolution			*/
	bits_per_pixel uint32      /* guess what			*/
	grayscale      uint32      /* 0 = color, 1 = grayscale,	*/
	red            fb_bitfield /* bitfield in fb mem if true color, */
	green          fb_bitfield /* else only length is significant */
	blue           fb_bitfield
	transp         fb_bitfield /* transparency			*/
	nonstd         uint32      /* != 0 Non standard pixel format */
	activate       uint32      /* see FB_ACTIVATE_*		*/
	height         uint32      /* height of picture in mm    */
	width          uint32      /* width of picture in mm     */
	accel_flags    uint32      /* (OBSOLETE) see fb_info.flags */
	/* Timing: All values in pixclocks, except pixclock (of course) */
	pixclock     uint32 /* pixel clock in ps (pico seconds) */
	left_margin  uint32 /* time from sync to picture	*/
	right_margin uint32 /* time from picture to sync	*/
	upper_margin uint32 /* time from sync to picture	*/
	lower_margin uint32
	hsync_len    uint32    /* length of horizontal sync	*/
	vsync_len    uint32    /* length of vertical sync	*/
	sync         uint32    /* see FB_SYNC_*		*/
	vmode        uint32    /* see FB_VMODE_*		*/
	rotate       uint32    /* angle we rotate counter clockwise */
	colorspace   uint32    /* colorspace for FOURCC-based modes */
	reserved     [4]uint32 /* Reserved for future compatibility */
}
