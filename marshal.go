package lk

// MarshalText implements the TextMarshaler interface
func (k *Key) MarshalText() (text []byte, err error) {
	text, err = []byte(k.String()), nil
	return text, err
}

// UnmarshalText implements the TextUnmarshaler interface
func (k *Key) UnmarshalText(text []byte) (err error) {
	k, err = Parse(string(text))
	return err
}
