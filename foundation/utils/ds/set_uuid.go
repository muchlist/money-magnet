package ds

import "github.com/google/uuid"

type setUUIDMapBased struct {
	setMap map[uuid.UUID]struct{}
}

func NewUUIDSet() *setUUIDMapBased {
	return &setUUIDMapBased{
		setMap: make(map[uuid.UUID]struct{}),
	}
}

func (s *setUUIDMapBased) Add(uuidVal uuid.UUID) (success bool) {
	if _, exist := s.setMap[uuidVal]; exist {
		success = false
		return
	}
	s.setMap[uuidVal] = struct{}{}
	success = true
	return
}

func (s *setUUIDMapBased) AddAll(uuids []uuid.UUID) {
	for _, id := range uuids {
		s.Add(id)
	}
}

func (s *setUUIDMapBased) Remove(uuidVal uuid.UUID) {
	delete(s.setMap, uuidVal)
}

func (s *setUUIDMapBased) Reveal() []uuid.UUID {
	uuidSlice := make([]uuid.UUID, 0)
	for key := range s.setMap {
		uuidSlice = append(uuidSlice, key)
	}
	return uuidSlice
}

func (s *setUUIDMapBased) Len() int {
	return len(s.setMap)
}
