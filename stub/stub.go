package stub

type Stub struct {
}

func (k *Stub) Set(key, value []byte) error {
	return nil
}
func (k *Stub) Get(key []byte) ([]byte, error) {
	return nil, nil
}
func (k *Stub) Del(key []byte) error {
	return nil
}
