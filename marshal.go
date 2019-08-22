package lk

// MarshalText implements the MarshalText interface
func (k *Key) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// UnmarshalText implements the UnmarshalText interface
func (k *Key) UnmarshalText(text []byte) error {
	k, err := Parse(string(text))
	if err != nil {
		return err
	}
	return nil
}