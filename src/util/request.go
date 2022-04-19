package util

import (
	"encoding/json"
	"fmt"
  )

type Request struct {
	Name string			`json:"name"`
	Args interface{}	`json:"args"`
	Rets interface{}	`json:"rets`
}

func (r *Request) Decode(data string) {
	err := json.Unmarshal([]byte(data), r)
	if err != nil {
		fmt.Println("error occured during decode json to struct")
	}
}

func (r *Request) Encode() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("error occured during encode struct to json")
		return nil
	}
	return data
}