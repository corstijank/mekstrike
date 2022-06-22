package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/mitchellh/copystructure"
)

// Store is simple access layer to a dapr store
type Store struct {
	client dapr.Client
	store  string
}

type Saveable interface {
	GetKey() string
	GetIndices() []string
	Marshal() ([]byte, error)
}

type Readable interface {
	Unmarshal([]byte) (Readable, error)
}

func New(name string) (Store, error) {
	var err error
	client, err := dapr.NewClient()
	if err != nil {
		return Store{}, err
	}

	return Store{
		client: client,
		store:  name,
	}, nil
}

func (s *Store) Close() {
	s.client.Close()
}

func (s *Store) ReadMany(ctx context.Context, index string, helper Readable) ([]Readable, error) {
	value, err := s.client.GetState(ctx, s.store, index, make(map[string]string))
	if err != nil {
		return nil, err
	}
	if len(value.Value) == 0 {
		return make([]Readable, 0), nil
	}
	keys := make([]string, 0)
	err = json.Unmarshal(value.Value, &keys)
	if err != nil {
		return nil, err
	}

	result := make([]Readable, 0)
	for _, k := range keys {
		helper, err := s.Read(ctx, k, helper)
		if err != nil {
			log.Printf("Error reading %s: %+v", k, err)
		}

		val, err := copystructure.Copy(helper)
		if err != nil {
			log.Printf("Error copying %s to new instance: %+v", k, err)
		}
		result = append(result, val.(Readable))
	}
	return result, nil
}

func (s *Store) Read(ctx context.Context, key string, v Readable) (Readable, error) {
	value, err := s.client.GetState(ctx, s.store, key, make(map[string]string))
	if err != nil {
		return v, err
	}

	if len(value.Value) == 0 {
		return nil, fmt.Errorf("no item with key %s", key)
	}

	v, err = v.Unmarshal(value.Value)
	if err != nil {
		return v, err
	}
	return v, nil
}

func (s *Store) Persist(ctx context.Context, e Saveable) error {
	log.Printf("Persisting %s on indexes %+v", e, e.GetIndices())

	b, err := e.Marshal()
	if err != nil {
		return err
	}
	if err := s.client.SaveState(ctx, s.store, e.GetKey(), b, make(map[string]string)); err != nil {
		return err
	}

	for _, index := range e.GetIndices() {
		kf := false
		keys := make([]string, 0)

		value, err := s.client.GetState(ctx, s.store, index, make(map[string]string))
		if err != nil {
			return err
		}

		if len(value.Value) != 0 {
			err = json.Unmarshal(value.Value, &keys)
			if err != nil {
				return err
			}

			for _, k := range keys {
				if k == e.GetKey() {
					kf = true
				}
			}
		}
		if !kf {
			keys = append(keys, e.GetKey())
			b, err := json.Marshal(keys)
			if err != nil {
				return err
			}
			if err := s.client.SaveState(ctx, s.store, index, b, make(map[string]string)); err != nil {
				return err
			}
		}

	}
	return nil
}
