package ds

import "sort"

type Num interface {
	int | int64 | uint64 | float64
}

type setNumMapBased[T Num] struct {
	setMap map[T]struct{}
}

func NewNumberSet[T Num]() *setNumMapBased[T] {
	return &setNumMapBased[T]{
		setMap: make(map[T]struct{}),
	}
}

func (s *setNumMapBased[T]) Add(num T) (success bool) {
	if _, exist := s.setMap[num]; exist {
		success = false
		return
	}
	s.setMap[num] = struct{}{}
	success = true
	return
}

func (s *setNumMapBased[T]) Remove(num T) {
	delete(s.setMap, num)
}

func (s *setNumMapBased[T]) Reveal() []T {
	int64Slice := make([]T, 0)
	for key := range s.setMap {
		int64Slice = append(int64Slice, key)
	}
	return int64Slice
}

func (s *setNumMapBased[T]) RevealSorted() []T {
	int64Slice := s.Reveal()
	sort.Slice(int64Slice, func(i, j int) bool {
		return int64Slice[i] < int64Slice[j]
	})
	return int64Slice
}

func (s *setNumMapBased[T]) Len() int {
	return len(s.setMap)
}
