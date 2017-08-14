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
	"github.com/DeveloppSoft/go-ipfs-api"
	"encoding/json"
)

type LedgerMock struct {
	Repo string

	sh *shell.Shell // IPFS api
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageFeed struct {
	Messsages []string `json:"messages"`
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
	if sequence == "1" {
		return "My Message 1", nil
	} else if sequence == "undefined" {
		return "", fmt.Errorf("Invalid sequence: %s", sequence)
	}

	return "", nil
}

func (l *LedgerMock) GetLastSeq(peer_name string) (string, error) {
	return "Message", nil
}

func (l *LedgerMock) GetFeed(peer_name string) ([]string, error) {
	if peer_name == "user1" {
		return []string{"Message 1", "Message 2"}, nil
	} else {
		return []string{}, fmt.Errorf("Can't find feed for user: %s", peer_name)
	}
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
	api := NewAPI("1234", rep, trk)
	
	//Create a new request
	req, errRequest := http.NewRequest("GET", "/v0/sync", nil)
	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	api.r.ServeHTTP(w, req)

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
	api := NewAPI("1234", rep, trk)
	
	//Create a new request
	seq := "1"
	req, errRequest := http.NewRequest("GET", fmt.Sprintf("/v0/messages/user1/%s", seq), nil)

	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	api.r.ServeHTTP(w, req)

	expectedMessage := fmt.Sprintf("\"My Message %s\"", seq)
	
	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %d, was: %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != expectedMessage {
		t.Errorf("Message should be %s, was: %s", expectedMessage, w.Body)
	}
}

func TestGetMessageNotFound(t *testing.T) {
	// Make the repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	ipfsApi := "<host>:<port>"
	rep := NewLedgerMock("./", ipfsApi)

	// Load tracker and api
	trk, _ := tracker.NewTracker(ctx, "./", ipfsApi)
	api := NewAPI("1234", rep, trk)
	//Create a new request
	seq := "2"
	req, errRequest := http.NewRequest("GET", fmt.Sprintf("/v0/messages/user1/%s", seq), nil)

	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	api.r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Response code should be %d, was: %d", http.StatusNotFound, w.Code)
	}
}

func TestGetMessageWrongSequence(t *testing.T) {
	// Make the repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ipfsApi := "<host>:<port>"
	rep := NewLedgerMock("./", ipfsApi)

	// Load tracker and api
	trk, _ := tracker.NewTracker(ctx, "./", ipfsApi)
	api := NewAPI("1234", rep, trk)
	//Create a new request
	req, errRequest := http.NewRequest("GET", "/v0/messages/user1/undefined", nil)

	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	api.r.ServeHTTP(w, req)

	expectedError := "Invalid sequence: undefined"

	if w.Code != http.StatusBadRequest {
		t.Errorf("Response code should be %d, was: %d", http.StatusBadRequest, w.Code)
	}

	var err ErrorResponse
	errParse := json.Unmarshal(w.Body.Bytes(), &err)
	
	if errParse != nil {
		t.Errorf("Error parsing response: %s", errParse)
	}

	if err.Error != expectedError {
		t.Errorf("Message should be %s, was: %s", expectedError, err.Error)
	}
}

func TestGetMessagesFeed(t *testing.T) {
	// Make the repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ipfsApi := "<host>:<port>"
	rep := NewLedgerMock("./", ipfsApi)

	// Load tracker and api
	trk, _ := tracker.NewTracker(ctx, "./", ipfsApi)
	api := NewAPI("1234", rep, trk)
	//Create a new request
	req, errRequest := http.NewRequest("GET", "/v0/messages/user1", nil)

	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	api.r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response code should be %d, was: %d", http.StatusOK, w.Code)
	}

	var feed MessageFeed
	errParse := json.Unmarshal(w.Body.Bytes(), &feed)
	
	if errParse != nil {
		t.Errorf("Error parsing response: %s", errParse)
	}

	expectedLenFeed := 2
	if len(feed.Messsages) != 2 {
		t.Errorf("Expecting %d messages, found %d", expectedLenFeed, len(feed.Messsages))
	}

	if feed.Messsages[0] != "Message 1" {
		t.Errorf("Expecting 'Message 1', found %s", feed.Messsages[0])
	}
	if feed.Messsages[1] != "Message 2" {
		t.Errorf("Expecting 'Message 2', found %s", feed.Messsages[1])
	}
}

func TestGetMessagesNoFeed(t *testing.T) {
	// Make the repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ipfsApi := "<host>:<port>"
	rep := NewLedgerMock("./", ipfsApi)

	// Load tracker and api
	trk, _ := tracker.NewTracker(ctx, "./", ipfsApi)
	api := NewAPI("1234", rep, trk)
	//Create a new request
	req, errRequest := http.NewRequest("GET", "/v0/messages/user2", nil)

	if errRequest != nil {
		t.Fatal(errRequest)
	}

	//Record the response
	w := httptest.NewRecorder()
	api.r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Response code should be %d, was: %d", http.StatusBadRequest, w.Code)
	}

	var err ErrorResponse
	errParse := json.Unmarshal(w.Body.Bytes(), &err)
	
	if errParse != nil {
		t.Errorf("Error parsing response: %s", errParse)
	}

	expectedError := "Can't find feed for user: user2"
	if err.Error != expectedError {
		t.Errorf("Message should be %s, was: %s", expectedError, err.Error)
	}	
}

//TODO: Add test for get my id