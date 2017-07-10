package tracker

import (
	"os"
	"io"
	"bufio"

	"encoding/csv"
)

type MetaTracker struct {
	Following map[string]string

	f *os.File
}

func MetaOpen(path string) (MetaTracker, error) {
	metatrack := MetaTracker{}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	metatrack.f = f

	return metatrack, err
}

func (m *MetaTracker) Close() error {
	return m.f.Close()
}

func (m *MetaTracker) Parse() error {
	reader := csv.NewReader(bufio.NewReader(m.f))
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
	// First, add it to the file
	writer := csv.NewWriter(m.f)
	defer writer.Flush()

	err := writer.Write([]string{peer, hash})

	// Now, save it in memory
	if err == nil {
		m.Following[peer] = hash
	}

	return err
}
