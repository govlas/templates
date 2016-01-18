package amap

// template type Amap(K,V)
type K int
type V int

type MP map[K]V

type asyncCommand struct {
	batch func(MP)
}

// Pair using for key-value operations in AsyncMap
type Pair struct {
	First  K
	Second V
}

// Empty is empty struct for channels
type Empty struct{}

// Amap provides asynchronous and thread-safe access to map.
// All funcs of this class return channels to interact with map.
type Amap struct {
	cache    MP
	commands chan *asyncCommand
	closed   bool
}

// NewAmap creates new async map.
// cache indicates a size of channel of interim values.
func NewAmap(cache int) *Amap {
	ret := new(Amap)
	ret.cache = make(MP)
	ret.commands = make(chan *asyncCommand, cache)
	go func() {
		for cmd := range ret.commands {
			cmd.batch(ret.cache)
		}
		ret.cache = nil
	}()
	return ret
}

func (am *Amap) checkClosed() {
	if am.closed {
		panic("Async map is closed")
	}
}

// Set returns channel to set the value.
func (am *Amap) Set(key K) chan<- V {
	am.checkClosed()
	ch := make(chan V, 1)
	am.commands <- &asyncCommand{batch: func(m MP) {
		m[key] = <-ch
		close(ch)
	}}
	return ch
}

func (am *Amap) SetPair() chan Pair {
	am.checkClosed()
	ch := make(chan Pair)
	am.commands <- &asyncCommand{batch: func(m MP) {
		go func(pair chan Pair) {
			for p := range pair {
				am.Set(p.First) <- p.Second
			}
		}(ch)

	}}
	return ch
}

// Get return channel to get the value.
// If channel closed then map don't have value by the key.
func (am *Amap) Get(key K) <-chan V {
	am.checkClosed()
	ch := make(chan V, 1)
	am.commands <- &asyncCommand{batch: func(m MP) {
		if value, ok := m[key]; ok {
			ch <- value
		}
		close(ch)
	}}
	return ch
}

// Delete ask to delete the element and returns channel which indicates that the element is deleted.
func (am *Amap) Delete(key K) <-chan Empty {
	am.checkClosed()
	ch := make(chan Empty, 1)

	am.commands <- &asyncCommand{batch: func(m MP) {
		delete(m, key)
		ch <- Empty{}
		close(ch)

	}}
	return ch
}

// Len returns channel to get the len of map.
func (am *Amap) Len() <-chan int {
	am.checkClosed()
	ch := make(chan int, 1)
	am.commands <- &asyncCommand{batch: func(m MP) {
		ch <- len(m)
		close(ch)
	}}
	return ch
}

// List returns channel for get key-value pairs from map
func (am *Amap) List() <-chan Pair {
	am.checkClosed()
	ch := make(chan Pair, <-am.Len())
	am.commands <- &asyncCommand{batch: func(m MP) {
		for k, v := range m {
			ch <- Pair{k, v}
		}
		close(ch)
	}}
	return ch
}

// Close closes internal channels and finish work goroutine.
func (am *Amap) Close() {
	am.closed = true
	close(am.commands)
}

// Release return channel to get the value. Value deletes from map.
// If channel closed then map don't have value by the key.
func (am *Amap) Release(key K) <-chan V {
	am.checkClosed()
	ch := make(chan V, 1)
	am.commands <- &asyncCommand{batch: func(m MP) {
		if value, ok := m[key]; ok {
			ch <- value
			delete(m, key)
		}
		close(ch)
	}}
	return ch
}

func (am *Amap) Batch(f func(MP)) <-chan Empty {
	am.checkClosed()
	ch := make(chan Empty, 1)
	am.commands <- &asyncCommand{batch: func(m MP) {
		f(m)
		ch <- Empty{}
		close(ch)
	}}
	return ch
}
