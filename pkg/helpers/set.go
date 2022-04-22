package helpers

var exists = struct{}{}

type set struct {
	m map[string]struct{}
}

func NewSet() *set {
	s := &set{}
	s.m = make(map[string]struct{})
	return s
}

func (s *set) Add(value string) {
	s.m[value] = exists
}

func (s *set) Append(values []string) {
	for _, value := range values {
		s.m[value] = exists
	}
}

func (s *set) Remove(value string) {
	delete(s.m, value)
}

func (s *set) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

func (s *set) Get() []string {
	values := make([]string, 0)
	for k := range s.m {
		values = append(values, k)
	}
	return values
}
