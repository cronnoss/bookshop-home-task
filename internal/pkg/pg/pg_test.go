package pg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDial_EmptyDSN(t *testing.T) {
	dsn := ""
	db, err := Dial(dsn)
	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Equal(t, "no postgres DSN provided", err.Error())
}
