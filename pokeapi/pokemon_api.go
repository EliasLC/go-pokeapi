package pokeapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiUrl = "https://pokeapi.co/api/v2/"

func MakeApiRequest(resourseUrl string) (map[string]interface{}, error) {
	requestUrl := fmt.Sprintf("%s%s", apiUrl, resourseUrl)
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
