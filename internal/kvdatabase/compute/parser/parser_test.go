package parser

import (
	"github.com/PechatnovVladimir/concurrency_go/internal/kvdatabase/storage/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser_New(t *testing.T) {
	t.Parallel()

	e := engine.NewEngine()
	cp := NewCommandParser(e)

	require.NotNil(t, cp)
	assert.NotNil(t, cp.engine)
}

func TestCommandParser_Set(t *testing.T) {
	t.Parallel()

	e := engine.NewEngine()
	cp := NewCommandParser(e)

	res, ok := cp.Execute("SET key1 value1")

	assert.Equal(t, res, "OK")
	assert.Equal(t, true, ok)

	value, ok := cp.Execute("GET key1")

	assert.Equal(t, value, "value1")
	assert.Equal(t, true, ok)
}
