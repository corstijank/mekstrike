package unit

import (
	"encoding/json"
	"fmt"

	"github.com/corstijank/mekstrike/src/common/go/storage"
)

func (s *UnitStats) Marshal() ([]byte, error) {
	result, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UnitStats) Unmarshal(b []byte) (storage.Readable, error) {
	err := json.Unmarshal(b, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (u *UnitStats) GetKey() string {
	return fmt.Sprintf("%s-%s", u.Name, u.Model)
}

func (u *UnitStats) GetIndices() []string {
	return []string{"_units", fmt.Sprintf("_units_%s", u.Type), fmt.Sprintf("_units_%s_%d", u.Type, u.Size)}
}
