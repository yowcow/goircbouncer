package command

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPong(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Pong(buf, "hogefuga")

	assert.Nil(t, err)
	assert.Equal(t, "PONG :hogefuga\r\n", buf.String())
}

func TestNick(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Nick(buf, "hogefuga")

	assert.Nil(t, err)
	assert.Equal(t, "NICK hogefuga\r\n", buf.String())
}

func TestUser(t *testing.T) {
	buf := new(bytes.Buffer)
	err := User(buf, "hogefuga", 0, "Hoge Fuga")

	assert.Nil(t, err)
	assert.Equal(t, "USER hogefuga 0 * :Hoge Fuga\r\n", buf.String())
}

func TestJoin(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Join(buf, []string{"#hoge", "#fuga", "#foo", "#bar"})

	assert.Nil(t, err)
	assert.Equal(t, "JOIN #hoge,#fuga,#foo,#bar\r\n", buf.String())
}

func TestNames(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Names(buf, []string{"#hoge", "#fuga", "#foo", "#bar"})

	assert.Nil(t, err)
	assert.Equal(t, "NAMES :#hoge,#fuga,#foo,#bar\r\n", buf.String())
}
