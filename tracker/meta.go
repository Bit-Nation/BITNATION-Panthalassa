/*
Copyright 2017 Eliott Teissonniere

Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without restriction,
including without limitation the rights to use, copy, modify, merge,
publish, distribute, sublicense, and/or sell copies of the Software,
and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

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
