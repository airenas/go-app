package goapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHidePass_Mongo(t *testing.T) {
	assert.Equal(t, "mongodb://mongo:27017", HidePass("mongodb://mongo:27017"))
	assert.Equal(t, "mongodb://l:----@mongo:27017", HidePass("mongodb://l:olia@mongo:27017"))
}
