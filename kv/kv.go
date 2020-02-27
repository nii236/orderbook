package kv

import (
	badger "github.com/dgraph-io/badger/v2"
)

type Badger struct {
	*badger.DB
}

func NewBadger(conn *badger.DB) *Badger {
	return &Badger{conn}
}

func (b *Badger) Set(key, value []byte) error {
	return b.Update(func(tx *badger.Txn) error {
		tx.Set(key, value)
		return nil
	})
}
func (b *Badger) Get(key []byte) ([]byte, error) {
	result := []byte{}
	err := b.Update(func(tx *badger.Txn) error {
		item, err := tx.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			copy(result, val)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (b *Badger) Del(key []byte) error {
	return b.Update(func(tx *badger.Txn) error {
		tx.Delete(key)
		return nil
	})
}
