package ini

import (
	"fmt"
	"slices"
	"strings"
	"sync/atomic"
)

type Manager struct {
	options     Options
	sources     []*dataSource
	futures     []*dataSource
	sections    map[string]*Section
	sectionList []string
	batch       atomic.Bool
	mutex       Mutex
	ValueMapper func(string) string
}

func (m *Manager) Batch(fn func(m *Manager) error) error {
	if m.batch.Swap(true) == false {
		defer m.batch.Swap(false)
	}
	return fn(m)
}

// Append appends one or more data sources and reloads automatically.
func (m *Manager) Append(source any, others ...any) error {
	if err := m.append(source); err != nil {
		return err
	}

	for _, other := range others {
		if err := m.append(other); err != nil {
			return err
		}
	}

	if !m.batch.Load() {
		return m.flush()
	}

	return nil
}

func (m *Manager) append(source any) error {
	ds, err := parseDataSource(source)
	if err != nil {
		return err
	}
	m.futures = append(m.futures, ds)
	return nil
}

func (m *Manager) flush() error {
	for _, s := range m.sources {
		s.Lock()
	}
	defer func() {
		for _, s := range m.sources {
			s.Unlock()
		}
	}()
	for len(m.futures) > 0 {
		s := m.futures[0]
		if err := s.reload(m); err != nil {
			return err
		}
		s.Lock()
		m.futures = m.futures[1:]
		m.sources = append(m.sources, s)
	}
	return nil
}

// Reload reloads and parses all data sources.
func (m *Manager) Reload() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	clear(m.sections)
	clear(m.sectionList)
	m.sectionList = m.sectionList[:0]

	for _, s := range m.sources {
		if err := s.reload(m); err != nil {
			return err
		}
	}

	return nil
}

// NewSection creates a new section.
func (m *Manager) NewSection(name string) *Section {
	if (m.options.Insensitive || m.options.InsensitiveSections) && len(name) > 0 {
		name = strings.ToLower(name)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if slices.Contains(m.sectionList, name) {
		return m.sections[name]
	}

	m.sectionList = append(m.sectionList, name)
	m.sections[name] = newSection(m, name)

	return m.sections[name]
}

// GetSection returns section by given name.
func (m *Manager) GetSection(name string) (*Section, error) {
	if len(name) > 0 && m.options.Insensitive || m.options.InsensitiveSections {
		name = strings.ToLower(name)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if slices.Contains(m.sectionList, name) {
		return m.sections[name], nil
	}

	return nil, fmt.Errorf("section %q does not exist", name)
}

// HasSection returns true if the file contains a section with given name.
func (m *Manager) HasSection(name string) bool {
	section, _ := m.GetSection(name)
	return section != nil
}

// Section assumes named section exists and returns a zero-value when not.
func (m *Manager) Section(name string) *Section {
	sec, err := m.GetSection(name)
	if err != nil {
		sec = newSection(m, name)
	}
	return sec
}
