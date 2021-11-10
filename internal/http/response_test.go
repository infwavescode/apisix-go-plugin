package http

import (
	"net/http"
	"testing"

	"github.com/api7/ext-plugin-proto/go/A6"
	hrc "github.com/api7/ext-plugin-proto/go/A6/HTTPReqCall"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"

	"github.com/finecloud/apisix-oauth2-plugin/internal/util"
)

func getStopAction(t *testing.T, b *flatbuffers.Builder) *hrc.Stop {
	buf := b.FinishedBytes()
	res := hrc.GetRootAsResp(buf, 0)
	tab := &flatbuffers.Table{}
	if res.Action(tab) {
		assert.Equal(t, hrc.ActionStop, res.ActionType())
		stop := &hrc.Stop{}
		stop.Init(tab.Bytes, tab.Pos)
		return stop
	}
	return nil
}

func TestFetchChanges(t *testing.T) {
	r := CreateResponse()
	defer ReuseResponse(r)
	r.Write([]byte("hello"))
	h := r.Header()
	h.Set("foo", "bar")
	h.Add("foo", "baz")
	h.Add("cat", "dog")
	r.Write([]byte(" world"))
	assert.Equal(t, "dog", h.Get("cat"))
	builder := util.GetBuilder()
	assert.True(t, r.FetchChanges(1, builder))

	stop := getStopAction(t, builder)
	assert.Equal(t, uint16(200), stop.Status())
	assert.Equal(t, []byte("hello world"), stop.BodyBytes())

	res := http.Header{}
	assert.Equal(t, 3, stop.HeadersLength())
	for i := 0; i < stop.HeadersLength(); i++ {
		e := &A6.TextEntry{}
		stop.Headers(e, i)
		res.Add(string(e.Name()), string(e.Value()))
	}
	assert.Equal(t, h, res)
}

func TestFetchChangesEmptyResponse(t *testing.T) {
	r := CreateResponse()
	builder := util.GetBuilder()
	assert.False(t, r.FetchChanges(1, builder))
}

func TestFetchChangesStatusOnly(t *testing.T) {
	r := CreateResponse()
	r.WriteHeader(400)
	builder := util.GetBuilder()
	assert.True(t, r.FetchChanges(1, builder))

	stop := getStopAction(t, builder)
	assert.Equal(t, uint16(400), stop.Status())
}

func TestWriteHeaderTwice(t *testing.T) {
	r := CreateResponse()
	r.WriteHeader(400)
	r.WriteHeader(503)
	builder := util.GetBuilder()
	assert.True(t, r.FetchChanges(1, builder))

	stop := getStopAction(t, builder)
	assert.Equal(t, uint16(400), stop.Status())
}
