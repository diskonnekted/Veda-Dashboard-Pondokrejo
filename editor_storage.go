package main

import (
	"encoding/json"
	"os"
	"sync"
)

type MarkerCoord struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type EditorStateFile struct {
	HouseholdCoords  map[string]MarkerCoord `json:"household_coords"`
	CustomHouseholds []Household            `json:"custom_households"`
	GeoLayers        map[string]json.RawMessage `json:"geo_layers"`
}

type EditorStore struct {
	mu   sync.RWMutex
	path string
	data EditorStateFile
}

func NewEditorStore(path string) *EditorStore {
	return &EditorStore{
		path: path,
		data: EditorStateFile{
			HouseholdCoords:  map[string]MarkerCoord{},
			CustomHouseholds: []Household{},
			GeoLayers:        map[string]json.RawMessage{},
		},
	}
}

func (s *EditorStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var loaded EditorStateFile
	if err := json.Unmarshal(b, &loaded); err != nil {
		return err
	}
	if loaded.HouseholdCoords == nil {
		loaded.HouseholdCoords = map[string]MarkerCoord{}
	}
	if loaded.CustomHouseholds == nil {
		loaded.CustomHouseholds = []Household{}
	}
	if loaded.GeoLayers == nil {
		loaded.GeoLayers = map[string]json.RawMessage{}
	}
	s.data = loaded
	return nil
}

func (s *EditorStore) Save() error {
	s.mu.RLock()
	b, err := json.MarshalIndent(s.data, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return err
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	_ = os.Remove(s.path)
	return os.Rename(tmp, s.path)
}

func (s *EditorStore) Apply(base []Household) []Household {
	s.mu.RLock()
	coords := s.data.HouseholdCoords
	custom := s.data.CustomHouseholds
	s.mu.RUnlock()

	out := make([]Household, 0, len(base)+len(custom))
	for _, hh := range base {
		if c, ok := coords[hh.NoKK]; ok {
			hh.Latitude = c.Latitude
			hh.Longitude = c.Longitude
		}
		out = append(out, hh)
	}
	out = append(out, custom...)
	return out
}

type EditorSaveRequest struct {
	Updates          []struct {
		NoKK      string  `json:"no_kk"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"updates"`
	CustomHouseholds []Household `json:"custom_households"`
	GeoLayers        map[string]json.RawMessage `json:"geo_layers"`
}

func (s *EditorStore) UpdateAndSave(req EditorSaveRequest) error {
	s.mu.Lock()
	if s.data.HouseholdCoords == nil {
		s.data.HouseholdCoords = map[string]MarkerCoord{}
	}
	for _, u := range req.Updates {
		if u.NoKK == "" {
			continue
		}
		s.data.HouseholdCoords[u.NoKK] = MarkerCoord{Latitude: u.Latitude, Longitude: u.Longitude}
	}
	if req.CustomHouseholds != nil {
		s.data.CustomHouseholds = req.CustomHouseholds
	}
	if req.GeoLayers != nil {
		if s.data.GeoLayers == nil {
			s.data.GeoLayers = map[string]json.RawMessage{}
		}
		for k, v := range req.GeoLayers {
			if k == "" {
				continue
			}
			s.data.GeoLayers[k] = v
		}
	}
	s.mu.Unlock()

	return s.Save()
}

func (s *EditorStore) GeoLayersSnapshot() map[string]json.RawMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string]json.RawMessage, len(s.data.GeoLayers))
	for k, v := range s.data.GeoLayers {
		if v == nil {
			continue
		}
		b := make([]byte, len(v))
		copy(b, v)
		out[k] = b
	}
	return out
}

func (s *EditorStore) GeoLayer(name string) (json.RawMessage, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data.GeoLayers[name]
	if !ok || v == nil {
		return nil, false
	}
	b := make([]byte, len(v))
	copy(b, v)
	return b, true
}
