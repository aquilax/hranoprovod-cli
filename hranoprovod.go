package main

import (

)

type Hranoprovod struct {}

func NewHranoprovod() *Hranoprovod {
	return &Hranoprovod{}
}

func (hr *Hranoprovod) register() error {
	return nil
}

func (hr *Hranoprovod) search(q string) error {
	api := NewAPIClient(GetDefaultAPIClientOptions())
	api.Search(q)
	return nil
}