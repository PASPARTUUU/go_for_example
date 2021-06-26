package mysaga

import (
	"encoding/json"
	"fmt"
)

type Pay struct {
	Num int
	Str string
}

type wantData struct {
	Payload struct {
		Noise struct {
			Num int
			Str string
		}
		Num int
		Str string
	}
}

func UserDo() {

	sagaClient := Client{
		OwnerHost: "",
	}

	saga, _ := sagaClient.AddStep(
		"trans-name",
		"step-name",
		[]string{"qwe"},
		Begining,
		fqwe,
		fqwe,
		2,
		nil,
	)

	var data = struct {
		Payload struct {
			Noise struct {
				Num int
				Str string
			}
			Num int
			Str string
		}
	}{
		Payload: struct {
			Noise struct {
				Num int
				Str string
			}
			Num int
			Str string
		}{Noise: struct {
			Num int
			Str string
		}{Num: 999, Str: "noice"},
			Num: 456, Str: "asd"},
	}
	_ = data

	saga.ParseMe(func(b [][]byte) ([]interface{}, error) {
		var want wantData

		err := json.Unmarshal(b[0], &want)
		if err != nil {
			return nil, err
		}

		fmt.Printf("%+v\n", want)

		var savingParams = make([]interface{}, 0)
		savingParams = append(savingParams, want.Payload.Num)
		savingParams = append(savingParams, want.Payload.Str)

		return savingParams, nil
	}, func(b [][]byte) ([]interface{}, error) {
		var want wantData

		err := json.Unmarshal(b[0], &want)
		if err != nil {
			return nil, err
		}

		var savingParams = make([]interface{}, 0)
		savingParams = append(savingParams, want.Payload.Num)
		savingParams = append(savingParams, want.Payload.Str)

		return savingParams, nil
	})

	saga.Play(data)

	saga.Listen()

}

func fqwe(num int, str string) (string, string, Pay, error) {
	fmt.Println(num)
	fmt.Println(str)
	return "private_key", "public_key", Pay{Num: num, Str: str}, nil
}
