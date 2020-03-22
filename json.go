package abuse

import (
	"encoding/json"
	"io"
)

func jsonDecode(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
