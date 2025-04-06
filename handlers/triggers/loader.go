package triggers

import (
	"io"

	"github.com/julianshen/nebula/data"
	"gopkg.in/yaml.v3"
)

// LoadTrigger loads a trigger definition from a YAML reader
func LoadTrigger(r io.Reader) (*data.Trigger, error) {
	// Read all content from the reader
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Parse YAML into trigger struct
	var trigger data.Trigger
	err = yaml.Unmarshal(content, &trigger)
	if err != nil {
		return nil, err
	}

	return &trigger, nil
}
