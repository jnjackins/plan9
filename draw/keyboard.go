package draw

// Keyboardctl is the source of keyboard events.
type Keyboardctl struct {
	C <-chan rune // Channel on which keyboard characters are delivered.
}

// InitKeyboard connects to the keyboard and returns a Keyboardctl to listen to it.
func (d *Display) InitKeyboard() *Keyboardctl {
	ch := make(chan rune, 20)
	go kbdproc(d, ch)
	return &Keyboardctl{ch}
}

func kbdproc(d *Display, ch chan rune) {
	for {
		r, err := d.conn.ReadKbd()
		if err != nil && err.Error() != "EOF" {
			panic(err)
		} else if err != nil {
			d.ExitC <- struct{}{} // signal client to shut down the display
			return
		}
		ch <- r
	}
}
