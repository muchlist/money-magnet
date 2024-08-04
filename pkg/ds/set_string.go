package ds

import (
	"sort"
	"strings"
)

type setStringMapBased struct {
	setMap map[string]struct{}
}

func NewStringSet() *setStringMapBased {
	return &setStringMapBased{
		setMap: make(map[string]struct{}),
	}
}

func (s *setStringMapBased) Add(text string) (success bool) {
	if _, exist := s.setMap[text]; exist {
		success = false
		return
	}
	s.setMap[text] = struct{}{}
	success = true
	return
}

func (s *setStringMapBased) AddAll(texts []string) {
	for _, id := range texts {
		s.Add(id)
	}
}

func (s *setStringMapBased) Remove(text string) {
	delete(s.setMap, text)
}

func (s *setStringMapBased) Reveal() []string {
	stringSlice := make([]string, 0)
	for key := range s.setMap {
		stringSlice = append(stringSlice, key)
	}
	return stringSlice
}

func (s *setStringMapBased) RevealNotEmpty() []string {
	stringSlice := make([]string, 0)
	for key := range s.setMap {
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey != "" {
			stringSlice = append(stringSlice, trimmedKey)
		}
	}
	return stringSlice
}

func (s *setStringMapBased) RevealSorted() []string {
	stringSlice := s.Reveal()
	sort.Slice(stringSlice, func(i, j int) bool {
		return stringSlice[i] < stringSlice[j]
	})
	return stringSlice
}

func (s *setStringMapBased) RevealNotEmptySorted() []string {
	stringSlice := s.RevealNotEmpty()
	sort.Slice(stringSlice, func(i, j int) bool {
		return stringSlice[i] < stringSlice[j]
	})
	return stringSlice
}

func (s *setStringMapBased) Len() int {
	return len(s.setMap)
}
