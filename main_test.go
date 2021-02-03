package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGoodReturnsProperly(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("good", isGood("92.5% Online"))
	assert.Equal("good", isGood("Normal"))
	assert.Equal("good", isGood("nOrMaL"))
	assert.Equal("bad", isGood("1"))
	assert.Equal("bad", isGood("2"))
	assert.Equal("bad", isGood("Service Unavailable"))
}

/*
func TestCrowbarReturnsValidResult(t *testing.T) {
	assert := assert.New(t)
	status, err := fetchStatus()

	assert.NotEqual(nil, status)
	assert.Equal(nil, err)
	assert.NotEmpty(status.Services)
	assert.Less(int64(1606241083), status.Time)
	assert.Less(float32(0), status.Online)
}
*/

func TestReturnsValidRedditInstance(t *testing.T) {
	reddit, err := makeReddit()

	assert.NotEqual(t, nil, reddit)
	assert.Equal(t, nil, err)
}

/*
func TestNormalRun(t *testing.T) {
	assert.Equal(t, nil, run())
}
*/
