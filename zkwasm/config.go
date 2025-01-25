package zkwasm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var Config = map[string]interface{}{}

func init() {
	// Get the directory of the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error getting current file directory")
		return
	}
	dir := filepath.Dir(filename)

	// Construct the path to the config.json file
	configPath := filepath.Join(dir, "config.json")

	// Open the JSON file
	jsonFile, err := os.Open(configPath)
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
