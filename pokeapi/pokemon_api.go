package pokeapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const ApiUrl = "https://pokeapi.co/api/v2/"

func MakeApiRequest(requestUrl string) (map[string]interface{}, error) {
	res, err := http.Get(requestUrl)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	resultMap := map[string]interface{}{}
	json.Unmarshal(body, &resultMap)

	return resultMap, err
}
