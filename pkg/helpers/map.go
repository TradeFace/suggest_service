package helpers

// type Map struct {
// 	m map[string]interface{}
// }

// func NewMap(values map[string]interface{}) *Map {
// 	s := &Map{}
// 	s.m = make(map[string]interface{})
// 	switch values.(type) {
// 	case map[string]interface{}:
// 		s.Append(values)
// 	}
// 	return s
// }

// func (s *Map) Add(key string, value interface{}) {
// 	s.m[key] = value
// }

// func (s *Map) Append(values map[string]interface{}) {
// 	for key, value := range values {
// 		s.m[key] = value
// 	}
// }

// func (s *Map) Remove(key string) {
// 	delete(s.m, key)
// }

// func (s *Map) ContainsKey(key string) bool {
// 	_, c := s.m[key]
// 	return c
// }

// func (s *Map) ContainsValue(value interface{}) bool {
// 	// _, c := s.m[key]
// 	return c
// }

// func (s *Map) Get() []string {
// 	values := make([]string, 0)
// 	for k := range s.m {
// 		values = append(values, k)
// 	}
// 	return values
// }
