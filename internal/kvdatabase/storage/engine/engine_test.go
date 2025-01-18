package engine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEngine_New(t *testing.T) {
	t.Parallel()

	engine := NewEngine()
	require.NotNil(t, engine)
	assert.NotNil(t, engine.data)
}

func TestEngine_Set(t *testing.T) {
	t.Parallel()

	engine := &Engine{
		data: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}

	tests := map[string]struct {
		key           string
		value         string
		expectedValue string
	}{
		"set not existing key": {
			key:           "key4",
			value:         "new_value4",
			expectedValue: "new_value4",
		},
		"set existing key": {
			key:           "key1",
			value:         "new_value1",
			expectedValue: "new_value1",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine.Set(test.key, test.value)
			value, found := engine.Get(test.key)
			assert.Equal(t, test.expectedValue, value)
			assert.True(t, found)
		})
	}

}

func TestEngine_Get(t *testing.T) {
	t.Parallel()

	engine := &Engine{
		data: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}

	tests := map[string]struct {
		key           string
		expectedValue string
		found         bool
	}{
		"get not existing key": {
			key:           "key5",
			found:         false,
			expectedValue: "",
		},
		"get existing key": {
			key:           "key3",
			found:         true,
			expectedValue: "value3",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value, found := engine.Get(test.key)
			assert.Equal(t, test.expectedValue, value)
			assert.Equal(t, test.found, found)
		})
	}
}

func TestCache_Del(t *testing.T) {
	t.Parallel()

	engine := &Engine{
		data: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}
	tests := map[string]struct {
		key string
	}{
		"del not existing key": {
			key: "key_3",
		},
		"del existing key": {
			key: "key_1",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine.Del(test.key)
			value, found := engine.Get(test.key)
			assert.Equal(t, "", value)
			assert.False(t, found)
		})
	}
}
