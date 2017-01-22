# pixelgo
a pixelflut implementation written in go. pixelflut is a fun multiplayer game where everybody can draw single pixels onto a screen thats in the network. pixelgo uses the framebuffer to draw its pixels on the screen, so no X-Server is required (more precisely it will not work witha running X-Server). Using the framebuffer is a quite simple and low level approach with the downside that this program will only run under linux or systems that have an compatible fb-API.
## Install & Run
1. install go
2. set your GOPATH e.g. ```cd && mkdir && export GOPATH=~/go```
3. ```go get github.com/c3re/pixelgo```
4. run the program ```$GOPATH/bin/pixelflut```

The avalaible commandline args are: the framebuffer to use (defaults to /dev/fb0) and the TCP-port to use (default: 1234). If you want something else use the desired fb as first and the port as second argument.

## The protocol
the initial protocol is described here: https://cccgoe.de/wiki/Pixelflut
There a basically two Commands

1. SIZE -> outputs SIZE X Y where X and Y are the size of the screen you can draw to
2. PX X Y rrggbb -> returns nothing but paints a rrggbb-colored pixel to the position X Y on the screen.

The original protocol allows usage of rrggbbaa values too. But i dropped that one for performance reasons. And honestly: who wants transparency in a pixelwar?
