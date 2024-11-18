package functions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func FetchDataFromFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(target)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		return err
	}
	return nil
}

func FetchData(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching data from %s: %v", url, err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body from %s: %v", url, err)
		return err
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		log.Printf("Error unmarshaling JSON from %s: %v", url, err)
		log.Printf("Raw response: %s", string(body))
		return err
	}

	return nil
}
