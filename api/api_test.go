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

package api

import (
	"fmt"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/base64"
	"github.com/Bit-Nation/BITNATION-Panthalassa/repo"
	"github.com/Bit-Nation/BITNATION-Panthalassa/tracker"	
	"github.com/gin-gonic/gin"
	"github.com/DeveloppSoft/go-ipfs-api"
)

type LedgerMock struct {
	Repo string

	sh *shell.Shell // IPFS api
}

func NewLedgerMock(repo_path string, ipfs_api string) *LedgerMock {
	return &LedgerMock{Repo: repo_path, sh: shell.NewShell(ipfs_api)}
}

/***************************/
/* Implement interface for mocking Ledger
/***************************/
func (l *LedgerMock) Sync() error {
	// First, add the repo to ipfs
	return nil
}

func (l *LedgerMock) GetMessage(peer_name string, sequence string) (string, error) {
	return fmt.Sprintf("My Message %s", sequence), nil
}

func (l *LedgerMock) GetLastSeq(peer_name string) (string, error) {
	return "Message", nil
}

func (l *LedgerMock) GetFeed(peer_name string) ([]string, error) {
	return []string{"Message 1", "Message 2"}, nil
}

func (l *LedgerMock) Whoami() string {
	return "me"
}

func (l *LedgerMock) About(peer_name string) (string, error) {
	about := "About me"
	return about, nil
}

func (l *LedgerMock) SetAbout(about repo.About) error {
	return nil
}

func (l *LedgerMock) Publish(data string) error {
	return nil
}

func (l *LedgerMock) AddRessource(b64 string) (string, error) {
	return "hash", nil
}

func (l *LedgerMock) GetRessource(id string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte{}), nil
}

func (l *LedgerMock) Resolve(name string) (string, error) {
	return "name", nil
}

/***************************/

func TestSync(t *testing.T) {
	// Make the repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	ipfsApi := "<host>:<port>"
	rep := NewLedgerMock("./", ipfsApi)

	// Load tracker and api
	trk, _ := tracker.NewTracker(ctx, "./", ipfsApi)
	api := API{Repo: rep, Tracker: trk}
	//Create a new request
	req, errRequest := http.NewRequest("GET", "/sync", nil)

	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	r := gin.Default()

	r.GET("/sync", api.sync)
	
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %d, was: %d", http.StatusOK, w.Code)
	}
}

func TestGetMessage(t *testing.T) {
	// Make the repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	ipfsApi := "<host>:<port>"
	rep := NewLedgerMock("./", ipfsApi)

	// Load tracker and api
	trk, _ := tracker.NewTracker(ctx, "./", ipfsApi)
	api := API{Repo: rep, Tracker: trk}
	//Create a new request
	seq := "1"
	req, errRequest := http.NewRequest("GET", fmt.Sprintf("/get_message/user1/%s", seq), nil)

	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	r := gin.Default()

	expectedMessage := fmt.Sprintf("\"My Message %s\"", seq)
	
	r.GET("/get_message/:user/:seq", api.getMessage)
	
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %d, was: %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != expectedMessage {
		t.Errorf("Message should be %s, was: %s", expectedMessage, w.Body)
	}
}
