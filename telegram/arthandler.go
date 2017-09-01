package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StringRating struct {
	Rating string `json:"rating"`
}

type FurryNetwork struct {
	Rating int `json:"rating"`
}

func FALookup(id string) string {
	url := fmt.Sprintf("http://faexport.fortinj.com/submission/%s.json", id)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var stringRating StringRating
	json.Unmarshal(body, &stringRating)
	return stringRating.Rating
}

func FNLookup(id string) int {
	url := fmt.Sprintf("https://beta.furrynetwork.com/api/artwork/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var intRating FurryNetwork
	json.Unmarshal(body, &intRating)
	return intRating.Rating
}

func E621IDLookup(id string) string {
	url := fmt.Sprintf("https://e621.net/post/show.json?id=%s", id)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var stringRating StringRating
	json.Unmarshal(body, &stringRating)
	return stringRating.Rating
}

func E621MD5Lookup(id string) string {
	url := fmt.Sprintf("https://e621.net/post/show.json?md5=%s", id)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var stringRating StringRating
	json.Unmarshal(body, &stringRating)
	return stringRating.Rating
}
