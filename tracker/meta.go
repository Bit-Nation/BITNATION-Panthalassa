package tracker

import (
	"os"
	"io"
	"bufio"

	"encoding/csv"
)

type MetaTracker struct {
	Following map[string]string
	path string
}

func MetaOpen(meta_path string) (MetaTracker, error) {
	metatrack := MetaTracker{path: meta_path}

	f, err := metatrack.open()
	if err != nil {
		return metatrack, err
	}
	defer f.Close()

	err = metatrack.parse(f)
	return metatrack, err
}

func (m *MetaTracker) open() (*os.File, error) {
	return os.OpenFile(m.path, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
}

func (m *MetaTracker) parse(f *os.File) error {
	reader := csv.NewReader(bufio.NewReader(f))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// File is organized as <peer_id>,<current_hash>
		m.Following[line[0]] = line[1]
	}

	return nil
}

func (m *MetaTracker) Append(peer string, hash string) error {
	f, err := m.open()
	if err != nil {
		return err
	}
	defer f.Close()

	// First, add it to the file
	writer := csv.NewWriter(f)
	defer writer.Flush()

	err = writer.Write([]string{peer, hash})

	// Now, save it in memory
	if err == nil {
		m.Following[peer] = hash
	}

	return err
}
