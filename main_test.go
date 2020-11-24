package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGoodReturnsProperly(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(isGood("92.5% Online"), "good")
	assert.Equal(isGood("Normal"), "good")
	assert.Equal(isGood("nOrMaL"), "good")
	assert.Equal(isGood("1"), "bad")
	assert.Equal(isGood("2"), "bad")
	assert.Equal(isGood("Service Unavailable"), "bad")
}

func TestCrowbarReturnsValidResult(t *testing.T) {
	assert := assert.New(t)
	status, err := fetchStatus()

	assert.NotEqual(status, nil)
	assert.Equal(err, nil)
	assert.NotEmpty(status.Services)
	assert.Less(int64(1606241083), status.Time)
	assert.Less(float32(0), status.Online)
}

func TestReturnsValidRedditInstance(t *testing.T) {
	reddit, err := makeReddit()

	assert.NotEqual(t, reddit, nil)
	assert.Equal(t, err, nil)
}

func TestNormalRun(t *testing.T) {
	assert.Equal(t, run(), nil)
}
