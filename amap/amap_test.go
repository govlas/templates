package amap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAmap(t *testing.T) {
	m := NewAmap(1)
	m.Set(1) <- 1
	s := <-m.Get(1)
	assert.Equal(t, s, V(1), "Get() returning bad data from map")

	assert.Equal(t, <-m.Len(), 1, "Len() returning uncorrect value")

	<-m.Delete(1)

	assert.Equal(t, <-m.Len(), 0, "Len() returning uncorrect value after delete")

	_, ok := <-m.Get(1)
	assert.False(t, ok, "Get() returning data from empty map")

	pp := m.SetPair()
	for i := 0; i < 10; i++ {
		pp <- Pair{K(i), V(i)}
	}
	time.Sleep(time.Second)
	listCount := 0
	for pair := range m.List() {
		listCount++
		assert.Equal(t, pair.First, K(pair.Second), "List() returning bad data")
	}
	assert.Equal(t, listCount, 10, "List() returning few values")

	m.Close()
	func() {
		defer func() {
			assert.NotNil(t, recover(), "No panic on closed map")
		}()
		m.Get(1)
	}()
}

func BenchmarkSetToMap(b *testing.B) {
	m := NewAmap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(K(i)) <- 1
	}
}

func BenchmarkGetFromMap(b *testing.B) {
	m := NewAmap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(K(i)) <- 1
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-m.Get(K(i))
	}
}

func BenchmarkSetToMapParallel(b *testing.B) {
	m := NewAmap(b.N)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(K(i)) <- V(i)
			i++
		}
	})
}

func BenchmarkGetFromMapParallel(b *testing.B) {
	m := NewAmap(b.N)
	for i := 0; i < b.N; i++ {
		m.Set(K(i)) <- 1
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			<-m.Get(K(i))
			i++
		}
	})
}
