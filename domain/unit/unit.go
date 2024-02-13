package unit

import (
	"encoding/json"
	"fmt"

	"github.com/corstijank/mekstrike/domain/storage"
)

func (s *Stats) Marshal() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Stats) Unmarshal(b []byte) (storage.Readable, error) {
	err := json.Unmarshal(b, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (u *Stats) GetKey() string {
	return fmt.Sprintf("%s-%s", u.Name, u.Model)
}

func (u *Stats) GetIndices() []string {
	return []string{"_units", fmt.Sprintf("_units_%s", u.Type), fmt.Sprintf("_units_%s_%d", u.Type, u.Size)}
}
