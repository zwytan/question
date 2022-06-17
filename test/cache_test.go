package test

import (
	"github.com/stretchr/testify/assert"
	"question/pkg/cache"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := cache.NewCache()
	cache.SetMaxMemory("30kb")
	cache.Set("key", "value", time.Second*3)
	v, ok := cache.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "value", v.(string))
}
