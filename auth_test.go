package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidAuth(t *testing.T) {
	auth := NewBasicAuth("roger", "rabbit")
	assert.Equal(t, auth.Login, "roger")
	assert.Equal(t, auth.Password, "rabbit")

	req := new(http.Request)
	req.Header = make(map[string][]string)
	req.SetBasicAuth("roger", "rabbit")
	assert.Equal(t, true, auth.ValidAuth(req))
	req.SetBasicAuth("gerard", "bouchard")
	assert.Equal(t, false, auth.ValidAuth(req))
	req.SetBasicAuth("roger", "")
	assert.Equal(t, false, auth.ValidAuth(req))
}
