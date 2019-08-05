package mas_test

import (
	"bytes"
	"encoding/json"
	"log"
	"reflect"
	"testing"

	"github.com/nminhquan/mas/model"
	"github.com/nminhquan/mas/utils"
)

func TestMarshallingToJSON(t *testing.T) {
	var buff = new(bytes.Buffer)
	var encoder = json.NewEncoder(buff)

	accInfo := model.AccountInfo{
		Id:      utils.NewSHAHash("accountNumber"),
		Number:  "accountNumber",
		Balance: 0.0,
	}

	ins := model.Instruction{
		Type: "account",
		Data: accInfo,
	}
	err := encoder.Encode(ins)
	if err != nil {
		log.Fatalln("cannot encode: ", ins)
	}
	log.Println("encoded data: ", buff.String())
}

func TestUnMarshallingFromJSON(t *testing.T) {
	encoded := "{\"Type\":\"account\",\"Data\":{\"Id\":\"194d30a1ecd8d5dfd8a5947d64eaff81bc9fac38\",\"Number\":\"accountNumber\",\"Balance\":0}}"
	var decoder = json.NewDecoder(bytes.NewBufferString(encoded))
	ins := model.Instruction{Data: model.AccountInfo{}}
	err := decoder.Decode(&ins)
	if err != nil {
		log.Fatalln("cannot decode: ", err)
	}
	log.Println("decoded data: type:", reflect.TypeOf(ins), " data = ", ins)
}
