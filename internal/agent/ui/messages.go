package ui

type Messages struct {
	M map[string]string
}

func (m *Messages) init() {
	m.M = map[string]string{}
}

func (m *Messages) Set(key, value string) {
	if m.M == nil {
		m.init()
	}
	m.M[key] = value
}

func (m *Messages) Get(key string) string {
	if m.M == nil {
		return ""
	}
	return m.M[key]
}

func (m *Messages) Clear(key string) {
	if m.M != nil {
		delete(m.M, key)
	}
}

func (m *Messages) ClearAll() {
	if m.M != nil {
		m.M = map[string]string{}
	}
}