package utils

import (
	"bytes"
	"encoding/gob"
	"mas/model"
)

func SerializeMessage(ins model.Instruction) []byte {
	b := new(bytes.Buffer)
	gob.Register(model.AccountInfo{})
	gob.Register(model.PaymentInfo{})
	if err := gob.NewEncoder(b).Encode(ins); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func DeserializeMessage(in []byte) model.Instruction {
	gob.Register(model.AccountInfo{})
	gob.Register(model.PaymentInfo{})
	dec := gob.NewDecoder(bytes.NewBuffer(in))
	var ins model.Instruction
	if err := dec.Decode(&ins); err != nil {
		panic(err)
	}
	return ins
}
