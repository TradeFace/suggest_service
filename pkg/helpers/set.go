package helpers

var exists = struct{}{}

type Set struct {
	m map[string]struct{}
}

func NewSet(values interface{}) *Set {
	s := &Set{}
	s.m = make(map[string]struct{})
	switch values := values.(type) {
	case []string:
		s.Append(values)
	}
	return s
}

func (s *Set) Add(value string) {
	s.m[value] = exists
}

func (s *Set) Append(values []string) {
	for _, value := range values {
		s.m[value] = exists
	}
}

func (s *Set) Remove(value string) {
	delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *Set) Get() []string {
	values := make([]string, 0)
	for k := range s.m {
		values = append(values, k)
	}
	return values
}
