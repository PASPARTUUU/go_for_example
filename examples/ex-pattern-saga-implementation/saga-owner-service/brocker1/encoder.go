package brocker1

import "encoding/json"

type Encoder interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	ContentType() string
}

// 

type JsonEncoder struct{}

func (j *JsonEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JsonEncoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j *JsonEncoder) ContentType() string {
	return "application/json"
}
