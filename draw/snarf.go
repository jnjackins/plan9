package draw

// TODO
func (d *Display) ReadSnarf() ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.conn.ReadSnarf()
}

// WriteSnarf writes the data to the snarf buffer.
func (d *Display) WriteSnarf(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	err := d.conn.WriteSnarf(data)
	if err != nil {
		return err
	}
	return nil
}
