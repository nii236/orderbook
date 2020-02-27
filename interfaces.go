package orderbook

type KV interface {
	Set(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Del(key []byte) error
}
