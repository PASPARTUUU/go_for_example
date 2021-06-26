package mysaga

import (
	"encoding/json"
	"errors"
)

var ErrNotFound = errors.New("ErrNotFound")

type SagaStore interface {
	SetIncomingData(key string, data interface{})
	GetIncomingData(key string) (interface{}, error)
	SetOutgoingData(key string, data interface{})
	GetOutgoingData(key string) (interface{}, error)

	GetRedoCount(key string) (int, error)
	IncRedoCount(key string) error

	// // SetParamsForCompensate -
	// SetParamsForExecute
	// // GetParamsForExecute - получить параметры для исполнения шага по ключу идемпотентности
	// GetParamsForExecute(key string) ([]interface{}, error)
	// // SetParamsForCompensate -
	// SetParamsForCompensate(key string, data interface{})
	// // GetParamsForExecute - получить параметры для компенсации шага по ключу идемпотентности
	// GetParamsForCompensate(key string) ([]interface{}, error)

	// TODO: подумать над названием метода
	// GetForConsistancy - проверяет наличие записи с успешно выполненым шагом по ключу идемпотентности
	GetForConsistancy(key string) bool
}

type Store struct {
	data         map[string]interface{}
	execRes      map[string]interface{}
	redoTryCount map[string]int
}

func (s *Store) SetIncomingData(key string, data interface{}) {
	s.data[key] = data
}

func (s *Store) GetIncomingData(key string) (interface{}, error) {
	res, ok := s.data[key]
	if !ok {
		return nil, ErrNotFound
	}

	return res, nil
}

func (s *Store) SetOutgoingData(key string, data interface{}) {
	s.execRes[key] = data
}

func (s *Store) GetOutgoingData(key string) (interface{}, error) {
	execRes, ok := s.execRes[key]
	if !ok {
		return nil, ErrNotFound
	}

	b, err := json.Marshal(execRes)
	if err != nil {
		return nil, err
	}
	var res []interface{}
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Store) GetRedoCount(key string) (int, error) {
	return s.redoTryCount[key], nil
}

func (s *Store) IncRedoCount(key string) error {
	s.redoTryCount[key]++
	return nil
}

func (s *Store) GetForConsistancy(key string) bool {
	var consistancy bool
	return consistancy
}
