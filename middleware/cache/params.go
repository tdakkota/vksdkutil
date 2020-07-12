package cache

import (
	"sort"

	"github.com/SevereCloud/vksdk/api"
)

// Param is VK API method argument.
type Param struct {
	Key   string
	Value interface{}
}

// OrderedParams is ordered VK API method arguments.
type OrderedParams []Param

func CreateParams(params api.Params) OrderedParams {
	result := make([]Param, 0, len(params))

	for k, v := range params {
		result = append(result, Param{
			Key:   k,
			Value: v,
		})
	}

	sort.Sort(OrderedParams(result))
	return result
}

func (o OrderedParams) Len() int {
	return len(o)
}

func (o OrderedParams) Less(i, j int) bool {
	return o[i].Key < o[j].Key
}

func (o OrderedParams) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
