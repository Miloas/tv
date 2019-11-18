package tv

import (
	"os"
	"github.com/iancoleman/orderedmap"
	"encoding/json"
	"io/ioutil"
)


// ReadVersionFile : read version info from semver json
func ReadVersionFile(path string) (*orderedmap.OrderedMap, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	o := orderedmap.New()
	err = json.Unmarshal(buf, &o)

	return o, err
}

// WriteVersionFile : write version info to semver json
func WriteVersionFile(o *orderedmap.OrderedMap, path string) error {
	buf, err := json.MarshalIndent(&o, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, buf, 0644)

	return err
}

// IsFileExist : check path exist, but notice that we dont catch permission error here
func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
