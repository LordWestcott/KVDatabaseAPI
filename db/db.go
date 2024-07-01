package db

import (
	"errors"
	"sync"
)

type Database struct {
	Data map[string]interface{}
	lock sync.RWMutex
}

type IDatabase interface {
	GetAllKeys() ([]string, error)
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

func NewDatabase() *Database {
	return &Database{
		Data: make(map[string]interface{}),
	}
}

func initCheck(d *Database) error {
	if d.Data == nil {
		return errors.New("database is not initialized")
	}
	return nil
}

func (d *Database) GetAllKeys() ([]string, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	out := make([]string, 0)

	if err := initCheck(d); err != nil {
		return nil, err
	}

	for k, _ := range d.Data {
		out = append(out, k)
	}

	return out, nil
}

func (d *Database) Get(key string) (interface{}, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if err := initCheck(d); err != nil {
		return nil, err
	}

	return d.Data[key], nil
}

func (d *Database) Set(key string, value interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if err := initCheck(d); err != nil {
		return err
	}

	d.Data[key] = value

	return nil
}

func (d *Database) Delete(key string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	if err := initCheck(d); err != nil {
		return err
	}

	delete(d.Data, key)

	return nil
}
