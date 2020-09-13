package paramsutil

import "github.com/SevereCloud/vksdk/v2/api"

func Append(name string, value interface{}, params ...api.Params) []api.Params {
	if len(params) < 1 {
		params = append(params, api.Params{})
	}

	params[0][name] = value
	return params
}

func Find(name string, params ...api.Params) (interface{}, bool) {
	for _, set := range params {
		if v, ok := set[name]; ok {
			return v, ok
		}
	}

	return nil, false
}

func Delete(name string, params ...api.Params) bool {
	for _, set := range params {
		if _, ok := set[name]; ok {
			delete(set, name)
			return ok
		}
	}

	return false
}
