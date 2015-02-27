// +build plan9

package drawfcall

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"image"
	"os"
	"strings"
)

type Conn struct {
	winname              []byte
	newfd, datafd, ctlfd *os.File
	mousedev             *bufio.Reader
	kbddev, consctl      *os.File
	kbdreader            *bufio.Reader
}

func New() (*Conn, error) {
	//fmt.Println("New")
	winname, err := ioutil.ReadFile("/dev/winname")
	if err != nil {
		return nil, err
	}
	newfd, err := os.OpenFile("/dev/draw/new", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 12*12)
	_, err = newfd.Read(buf)
	if err != nil {
		return nil, err
	}
	params := strings.Fields(string(buf))
	dirno := params[0]
	datafd, err := os.OpenFile("/dev/draw/"+dirno+"/data", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	ctlfd, err := os.OpenFile("/dev/draw/"+dirno+"/ctl", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	mousedev, err := os.OpenFile("/dev/mouse", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	mousebuf := bufio.NewReader(mousedev)
	kbddev, err := os.OpenFile("/dev/cons", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	consctl, err := os.OpenFile("/dev/consctl", os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}
	consctl.WriteString("rawon")
	kbdreader := bufio.NewReader(kbddev)

	return &Conn{
		winname:   winname,
		newfd:     newfd,
		datafd:    datafd,
		ctlfd:     ctlfd,
		mousedev:  mousebuf,
		kbddev:    kbddev,
		consctl:   consctl,
		kbdreader: kbdreader,
	}, nil
}

// TODO
func (c *Conn) Close() error {
	fmt.Println("Close")
	return nil
}

func (c *Conn) Init(label, winsize string) error {
	fmt.Println("Init")
	return nil
}

func (c *Conn) ReadMouse() (m Mouse, resized bool, err error) {
	//fmt.Println("ReadMouse")
	var symbol uint8
	var x, y, buttons, ms int
	_, err = fmt.Fscanf(c.mousedev, "%c %d %d %d %d ", &symbol, &x, &y, &buttons, &ms)
	m = Mouse{image.Pt(x, y), buttons, ms}
	if symbol == 'r' {
		resized = true
	}
	return
}

func (c *Conn) ReadKbd() (r rune, err error) {
	//fmt.Println("ReadKbd")
	r, _, err = c.kbdreader.ReadRune()
	return
}

func (c *Conn) MoveTo(p image.Point) error {
	fmt.Println("MoveTo")
	return nil
}

func (c *Conn) Cursor(cursor *Cursor) error {
	fmt.Println("Cursor")
	return nil
}

func (c *Conn) BounceMouse(m *Mouse) error {
	fmt.Println("BounceMouse")
	return nil
}

func (c *Conn) Label(label string) error {
	fmt.Println("Label")
	return nil
}

// Return values are bytes copied, actual size, error.
func (c *Conn) ReadSnarf() ([]byte, error) {
	//fmt.Println("ReadSnarf")
	return ioutil.ReadFile("/dev/snarf")
}

func (c *Conn) WriteSnarf(snarf []byte) error {
	fmt.Println("WriteSnarf")
	return nil
}

func (c *Conn) Top() error {
	fmt.Println("Top")
	return nil
}

func (c *Conn) Resize(r image.Rectangle) error {
	fmt.Println("Resize")
	return nil
}

func (c *Conn) ReadDraw(b []byte) (int, error) {
	//fmt.Println("ReadDraw")
	return c.datafd.Read(b)
}

func (c *Conn) WriteDraw(b []byte) (int, error) {
	//fmt.Println("WriteDraw", b)
	return c.datafd.Write(b)
}

func (c *Conn) AttachScreen() (error) {
	//fmt.Println("AttachScreen")
	buf := make([]byte, 1+4+1+len(c.winname))
	buf[0] = 'n'
	bplong(buf[1:], 1)
	buf[5] = byte(len(c.winname))
	copy(buf[6:], c.winname)
	//fmt.Println(buf)
	_, err := c.datafd.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conn) GetInfo() ([]byte, error) {
	//fmt.Println("GetInfo")
	buf := make([]byte, 12*12)
	_, err := c.ctlfd.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func bplong(b []byte, n uint32) {
	binary.LittleEndian.PutUint32(b, n)
}

func bpshort(b []byte, n uint16) {
	binary.LittleEndian.PutUint16(b, n)
}
