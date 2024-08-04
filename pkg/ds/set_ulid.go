package ds

import "github.com/muchlist/moneymagnet/pkg/xulid"

type setULIDMapBased struct {
	setMap map[xulid.ULID]struct{}
}

func NewULIDSet() *setULIDMapBased {
	return &setULIDMapBased{
		setMap: make(map[xulid.ULID]struct{}),
	}
}

func (s *setULIDMapBased) Add(ulidVal xulid.ULID) (success bool) {
	if _, exist := s.setMap[ulidVal]; exist {
		success = false
		return
	}
	s.setMap[ulidVal] = struct{}{}
	success = true
	return
}

func (s *setULIDMapBased) AddAll(ulids []xulid.ULID) {
	for _, id := range ulids {
		s.Add(id)
	}
}

func (s *setULIDMapBased) Remove(ulidVal xulid.ULID) {
	delete(s.setMap, ulidVal)
}

func (s *setULIDMapBased) Reveal() []xulid.ULID {
	ulidSlice := make([]xulid.ULID, 0)
	for key := range s.setMap {
		ulidSlice = append(ulidSlice, key)
	}
	return ulidSlice
}

func (s *setULIDMapBased) Len() int {
	return len(s.setMap)
}
