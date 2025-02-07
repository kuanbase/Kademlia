package kencode

import (
	"fmt"
	"strings"
)

const (
	PING     = "PING"
	PONG     = "PONG"
	GETID    = "GETID"
	RETURNID = "RETURNID"
)

type KenCode struct {
	Commands []string
	Values   []any
}

func NewKenCode() *KenCode {
	return &KenCode{Commands: make([]string, 0), Values: make([]any, 0)}
}

type Encoder struct {
	kenCode *KenCode
}

func (e *Encoder) Ping(address string) *Encoder {
	e.kenCode.Commands = append(e.kenCode.Commands, PING)
	e.kenCode.Values = append(e.kenCode.Values, address)
	return e
}

func (e *Encoder) ResponsePing(address string) *Encoder {
	e.kenCode.Commands = append(e.kenCode.Commands, PONG)
	e.kenCode.Values = append(e.kenCode.Values, address)
	return e
}

func (e *Encoder) GetID(address string) *Encoder {
	e.kenCode.Commands = append(e.kenCode.Commands, GETID)
	e.kenCode.Values = append(e.kenCode.Values, address) // 請求方需要出事地址
	return e
}

func (e *Encoder) ResponseGETID(id string) *Encoder {
	e.kenCode.Commands = append(e.kenCode.Commands, RETURNID)
	e.kenCode.Values = append(e.kenCode.Values, id) // 返回Node ID
	return e
}

func (e *Encoder) Store(data []byte) *Encoder {
	e.kenCode.Commands = append(e.kenCode.Commands, "STORE")
	e.kenCode.Values = append(e.kenCode.Values, data)
	return e
}

func (e *Encoder) Encode() string {
	var kencodes []string

	for i := 0; i < len(e.kenCode.Values); i++ {
		str := fmt.Sprintf("[%s]=[%v];", e.kenCode.Commands[i], e.kenCode.Values[i])
		kencodes = append(kencodes, str)
	}

	return strings.Join(kencodes, "")
}

func NewEncoder() *Encoder {
	return &Encoder{kenCode: NewKenCode()}
}

type Decoder struct {
	kenCode string
}

func NewDecoder(kenCode string) *Decoder {
	return &Decoder{kenCode: kenCode}
}

func (d *Decoder) Decode() *KenCode {
	kenCodes := strings.Split(d.kenCode, ";")

	kenCode := NewKenCode()

	for i := 0; i < len(kenCodes); i++ {
		pair := strings.Split(kenCodes[i], "=")
		if len(pair) != 2 {
			continue
		}

		command := pair[0]
		value := pair[1]

		command = strings.TrimPrefix(command, "[")
		command = strings.TrimSuffix(command, "]")

		value = strings.TrimPrefix(value, "[")
		value = strings.TrimSuffix(value, "]")

		kenCode.Commands = append(kenCode.Commands, command)
		kenCode.Values = append(kenCode.Values, value)
	}

	return kenCode
}
