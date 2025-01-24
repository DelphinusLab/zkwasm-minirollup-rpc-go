package zkwasm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Config = map[string]interface{}{}

func init() {
	// Open the JSON file
	jsonFile, err := os.Open("./zkwasm/config.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()

	// Read the JSON file
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal the JSON data into a map
	err = json.Unmarshal(byteValue, &Config)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		return
	}
}
