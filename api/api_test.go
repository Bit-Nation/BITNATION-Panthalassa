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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Bit-Nation/BITNATION-Panthalassa/repo"
	"github.com/Bit-Nation/BITNATION-Panthalassa/tracker"	
	"github.com/gin-gonic/gin"
)

func TestSync(t *testing.T) {
	// Make the repo
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	ipfsApi := "<host>:<port>"
	rep := repo.NewLedger("./", ipfsApi)

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
