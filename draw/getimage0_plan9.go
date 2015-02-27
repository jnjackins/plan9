// +build plan9

package draw

import "strings"

func (d *Display) getimage0(i *Image) (*Image, error) {
	if i != nil {
		i.free()
		*i = Image{}
	}

	err := d.conn.AttachScreen()
	if err != nil {
		return nil, err
	}
	info, err := d.conn.GetInfo()
	if err != nil {
		return nil, err
	}

	pix, _ := ParsePix(strings.TrimSpace(string(info[2*12 : 3*12])))
	if i == nil {
		i = new(Image)
	}
	*i = Image{
		Display: d,
		id:      1,
		Pix:     pix,
		Depth:   pix.Depth(),
		Repl:    atoi(info[3*12:]) > 0,
		R:       ator(info[4*12:]).Inset(4),
		Clipr:   ator(info[8*12:]).Inset(4),
	}
	d.imageid++
	return i, nil
}