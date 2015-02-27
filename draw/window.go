package draw

import (
	"fmt"
	"image"
)

var screenid uint32

func (i *Image) AllocScreen(fill *Image, public bool) (*Screen, error) {
	i.Display.mu.Lock()
	defer i.Display.mu.Unlock()
	return i.allocScreen(fill, public)
}

func (i *Image) allocScreen(fill *Image, public bool) (*Screen, error) {
	d := i.Display
	if d != fill.Display {
		return nil, fmt.Errorf("allocscreen: image and fill on different displays")
	}
	var id uint32
	for try := 0; ; try++ {
		if try >= 25 {
			return nil, fmt.Errorf("allocscreen: cannot find free id")
		}
		a := d.bufimage(1 + 4 + 4 + 4 + 1)
		screenid++
		id = screenid
		a[0] = 'A'
		bplong(a[1:], id)
		bplong(a[5:], i.id)
		bplong(a[9:], fill.id)
		if public {
			a[13] = 1
		}
		if err := d.flush(false); err == nil {
			break
		}
	}
	s := &Screen{
		Display: d,
		id:      id,
		Fill:    fill,
	}
	return s, nil
}

// Free frees the server resources associated with the screen.
func (s *Screen) Free() error {
	s.Display.mu.Lock()
	defer s.Display.mu.Unlock()
	return s.free()
}

func (s *Screen) free() error {
	if s == nil {
		return nil
	}
	d := s.Display
	a := d.bufimage(1 + 4)
	a[0] = 'F'
	bplong(a[1:], s.id)
	// flush(true) because screen is likely holding the last reference to window,
	// and we want it to disappear visually.
	return d.flush(true)
}

func allocwindow(i *Image, s *Screen, r image.Rectangle, ref int, val Color) (*Image, error) {
	d := s.Display
	i, err := allocImage(d, i, r, d.ScreenImage.Pix, false, val, s.id, ref)
	if err != nil {
		return nil, err
	}
	i.Screen = s
	i.next = s.Display.Windows
	s.Display.Windows = i
	return i, nil
}
