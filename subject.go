package redtape

type SubjectOption func(*Subject)

func SubjectName(n string) SubjectOption {
	return func(s *Subject) {
		s.Name = n
	}
}

func SubjectRole(role ...*Role) SubjectOption {
	return func(s *Subject) {
		s.Roles = append(s.Roles, role...)
	}
}

func SubjectMeta(meta ...map[string]interface{}) SubjectOption {
	return func(s *Subject) {
		if s.Meta == nil {
			s.Meta = make(map[string]interface{})
		}

		for _, md := range meta {
			for k, v := range md {
				s.Meta[k] = v
			}
		}
	}
}

type Subject struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Roles []*Role                `json:"roles"`
	Meta  map[string]interface{} `json:"meta"`
}

func NewSubject(id string, opts ...SubjectOption) Subject {
	sub := Subject{
		ID:   id,
		Meta: make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(&sub)
	}

	return sub
}

func (s Subject) EffectiveRoles() []*Role {
	var er []*Role
	for _, r := range s.Roles {
		er = append(er, r.EffectiveRoles()...)
	}

	return er
}

func (s Subject) String() string {
	return s.ID
}
