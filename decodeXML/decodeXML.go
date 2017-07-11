package decodeXML

import (
	"encoding/xml"
	"os"
)

func DecodeXML(v interface{}, filename *string) error {
	file, err := os.Open(*filename)
	if err != nil {
		return err
	}

	return xml.NewDecoder(file).Decode(v)
}
