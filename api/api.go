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
	"flag"
	"github.com/gin-gonic/gin"

	"github.com/Bit-Nation/BITNATION-Panthalassa/repo"
	"github.com/Bit-Nation/BITNATION-Panthalassa/tracker"
)

type API struct {
	Repo    repo.Ledger
	Tracker tracker.Tracker

	listen_address string
	r              *gin.Engine
}

func doResult(c *gin.Context, value interface{}, err error) {
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, value)
	}
}

func NewAPI(listen string, rep repo.Ledger, track tracker.Tracker) API {
	a := API{Repo: rep, Tracker: track}
	a.listen_address = listen

	a.r = gin.Default()

	// Build the router api
	a.r.GET("/sync", a.sync)

	// Messages
	a.r.GET("/get_message/:user/:seq", a.getMessage)
	a.r.GET("/get_last_seq/:user", a.getLastSeq)
	a.r.GET("/get_feed/:user", a.getFeed)

	// Profiles
	a.r.GET("/me", a.me)
	a.r.GET("/about/:user", a.about)

	a.r.POST("/set_about", a.setAbout)

	// Social actions
	a.r.GET("/follow/:user", a.follow)
	a.r.GET("/unfollow/:user", a.unFollow)
	a.r.GET("/following", a.getFollowing)

	// Publishing
	a.r.POST("/publish", a.publish)
	a.r.POST("/upload", a.upload)
	a.r.POST("/download/:id", a.download)

	return a
}

func (a *API) Run() error {
	return a.r.Run(a.listen_address)
}

func (a *API) sync(c *gin.Context) {
	var err error = nil
	//Ignore ipfs sync when running tests
	if flag.Lookup("test.v") == nil {
		err = a.Repo.Sync()
	}

	doResult(c, nil, err)
}

func (a *API) getMessage(c *gin.Context) {
	user := c.Params.ByName("user")
	seq := c.Params.ByName("seq")

	msg, err := a.Repo.GetMessage(user, seq)

	doResult(c, msg, err)
}

func (a *API) getLastSeq(c *gin.Context) {
	user := c.Params.ByName("user")

	seq, err := a.Repo.GetLastSeq(user)

	doResult(c, seq, err)
}

func (a *API) getFeed(c *gin.Context) {
	user := c.Params.ByName("user")

	feed, err := a.Repo.GetFeed(user)

	doResult(c, feed, err)
}

func (a *API) me(c *gin.Context) {
	id := a.Repo.Whoami()

	doResult(c, id, nil)
}

func (a *API) about(c *gin.Context) {
	user := c.Params.ByName("user")

	about, err := a.Repo.About(user)

	doResult(c, about, err)
}

func (a *API) setAbout(c *gin.Context) {
	about := repo.About{}

	about.Pseudo = c.PostForm("pseudo")
	about.Image = c.PostForm("image")
	about.ETHAddress = c.PostForm("eth_address")
	about.Description = c.PostForm("description")

	err := a.Repo.SetAbout(about)

	doResult(c, about, err)
}

func (a *API) follow(c *gin.Context) {
	user := c.Params.ByName("user")

	err := a.Tracker.Follow(user)

	doResult(c, nil, err)
}

func (a *API) unFollow(c *gin.Context) {
	user := c.Params.ByName("user")

	err := a.Tracker.UnFollow(user)

	doResult(c, nil, err)
}

func (a *API) getFollowing(c *gin.Context) {
	followed, err := a.Tracker.GetUsersFollowed()

	doResult(c, followed, err)
}

func (a *API) publish(c *gin.Context) {
	data := c.PostForm("data")

	err := a.Repo.Publish(data)

	doResult(c, nil, err)
}

func (a *API) upload(c *gin.Context) {
	data := c.PostForm("data_b64")

	id, err := a.Repo.AddRessource(data)

	doResult(c, id, err)

}

func (a *API) download(c *gin.Context) {
	id := c.Params.ByName("id")

	data, err := a.Repo.GetRessource(id)

	doResult(c, data, err)
}
