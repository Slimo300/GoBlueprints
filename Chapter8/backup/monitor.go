package backup

import (
	"fmt"
	"path/filepath"
	"time"
)

// Monitor is a type for monitoring file system
type Monitor struct {
	Paths       map[string]string
	Archiver    Archiver
	Destination string
}

// Now is a method that checks current state of monitored paths and in
// case of finding some changes creates new backup and updates monitor info
func (m *Monitor) Now() (int, error) {

	var counter int

	for path, lastHash := range m.Paths {

		newHash, err := DirHash(path)
		if err != nil {
			return counter, err
		}

		if newHash != lastHash {
			err := m.act(path)
			if err != nil {
				return counter, err
			}

			m.Paths[path] = newHash
			counter++
		}
	}
	return counter, nil
}

func (m *Monitor) act(path string) error {
	dirname := filepath.Base(path)
	filename := fmt.Sprintf(m.Archiver.DestFmt(), time.Now().UnixNano())
	return m.Archiver.Archive(path, filepath.Join(m.Destination, dirname, filename))
}
