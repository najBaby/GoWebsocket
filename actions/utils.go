package actions

import "encoding/json"

type response struct {
	Errors    []string    `json:"errors,omitempty"`
	Data      interface{} `json:"data"`
	Operation string      `json:""operation`
}

func convert(in interface{}, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, out)
	if err != nil {
		return err
	}
	return nil
}

func convertResponseTo(in interface{}, out interface{}) error {
	res := response{}
	err := convert(in, &res)
	if err != nil {
		return err
	}
	err = convert(res.Data, &out)
	if err != nil {
		return err
	}
	return nil
}
