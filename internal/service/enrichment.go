package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAge(name string) (int, error) {
	var res struct {
		Age int `json:"age"`
	}
	err := fetch(fmt.Sprintf("https://api.agify.io/?name=%s", name), &res)
	return res.Age, err
}

func GetGender(name string) (string, error) {
	var res struct {
		Gender string `json:"gender"`
	}
	err := fetch(fmt.Sprintf("https://api.genderize.io/?name=%s", name), &res)
	return res.Gender, err
}

func GetNationality(name string) (string, error) {
	var res struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	err := fetch(fmt.Sprintf("https://api.nationalize.io/?name=%s", name), &res)
	if err != nil || len(res.Country) == 0 {
		return "", err
	}
	return res.Country[0].CountryID, nil
}

func fetch(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
