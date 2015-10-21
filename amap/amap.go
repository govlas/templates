package amap

// template type Amap(K,V)
type K int
type V int

const (
	asyncCommandGet int = iota
	asyncCommandSet
	asyncCommandDelete
	asyncCommandLen
	asyncCommandList
	asyncCommandSetPair
)

type asyncCommand struct {
	data   chan V
	del    chan Empty
	length chan int
	pair   chan Pair

	typ int
	key K
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
	cache    map[K]V
	commands chan *asyncCommand
	closed   bool
}

// NewAmap creates new async map.
// cache indicates a size of channel of interim values.
func NewAmap(cache int) *Amap {
	ret := new(Amap)
	ret.cache = make(map[K]V)
	ret.commands = make(chan *asyncCommand, cache)
	go func() {
		for cmd := range ret.commands {
			switch cmd.typ {
			case asyncCommandGet:
				if value, ok := ret.cache[cmd.key]; ok {
					cmd.data <- value
				}
				close(cmd.data)

			case asyncCommandSet:
				ret.cache[cmd.key] = <-cmd.data
				close(cmd.data)

			case asyncCommandDelete:
				delete(ret.cache, cmd.key)
				cmd.del <- Empty{}
				close(cmd.del)

			case asyncCommandLen:
				cmd.length <- len(ret.cache)
				close(cmd.length)

			case asyncCommandList:
				for k, v := range ret.cache {
					cmd.pair <- Pair{k, v}
				}
				close(cmd.pair)
			case asyncCommandSetPair:
				go func(pair chan Pair) {
					for p := range pair {
						ret.Set(p.First) <- p.Second
					}
				}(cmd.pair)
			}
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
	am.commands <- &asyncCommand{data: ch, typ: asyncCommandSet, key: key}
	return ch
}

func (am *Amap) SetPair() chan Pair {
	am.checkClosed()
	ch := make(chan Pair)
	am.commands <- &asyncCommand{pair: ch, typ: asyncCommandSetPair}
	return ch
}

// Get return channel to get the value.
// If channel closed then map don't have value by the key.
func (am *Amap) Get(key K) <-chan V {
	am.checkClosed()
	ch := make(chan V, 1)
	am.commands <- &asyncCommand{data: ch, typ: asyncCommandGet, key: key}
	return ch
}

// Delete ask to delete the element and returns channel which indicates that the element is deleted.
func (am *Amap) Delete(key K) <-chan Empty {
	am.checkClosed()
	ch := make(chan Empty, 1)

	am.commands <- &asyncCommand{del: ch, typ: asyncCommandDelete, key: key}
	return ch
}

// Len returns channel to get the len of map.
func (am *Amap) Len() <-chan int {
	am.checkClosed()
	ch := make(chan int, 1)
	am.commands <- &asyncCommand{length: ch, typ: asyncCommandLen}
	return ch
}

// List returns channel for get key-value pairs from map
func (am *Amap) List() <-chan Pair {
	am.checkClosed()
	ch := make(chan Pair, <-am.Len())
	am.commands <- &asyncCommand{pair: ch, typ: asyncCommandList}
	return ch
}

// Close closes internal channels and finish work goroutine.
func (am *Amap) Close() {
	am.closed = true
	close(am.commands)
}
