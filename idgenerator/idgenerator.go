/*
Package idgenerator contains several Trace ID generators which can be
used by the Zipkin tracer. Additional third party generators can be plugged in
if they adhere to the IDGenerator interface.
*/
package idgenerator

import (
	"math/rand"
	"sync"
	"time"

	"github.com/ximply/traceid"
)

var (
	seededIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	// NewSource returns a new pseudo-random Source seeded with the given value.
	// Unlike the default Source used by top-level functions, this source is not
	// safe for concurrent use by multiple goroutines. Hence the need for a mutex.
	seededIDLock sync.Mutex
)

// IDGenerator interface can be used to provide the Zipkin Tracer with custom
// implementations to generate Trace IDs.
type IDGenerator interface {
	TraceID() traceid.TraceID                // Generates a new Trace ID
}

// NewRandom64 returns an ID Generator which can generate 64 bit trace
func NewRandom64() IDGenerator {
	return &randomID64{}
}

// NewRandom128 returns an ID Generator which can generate 128 bit trace
func NewRandom128() IDGenerator {
	return &randomID128{}
}

// NewRandomTimestamped generates 128 bit time sortable traceid's
func NewRandomTimestamped() IDGenerator {
	return &randomTimestamped{}
}

// randomID64 can generate 64 bit traceid's and 64 bit spanid's.
type randomID64 struct{}

func (r *randomID64) TraceID() (id traceid.TraceID) {
	seededIDLock.Lock()
	id = traceid.TraceID{
		Low: uint64(seededIDGen.Int63()),
	}
	seededIDLock.Unlock()
	return
}

// randomID128 can generate 128 bit traceid's
type randomID128 struct{}

func (r *randomID128) TraceID() (id traceid.TraceID) {
	seededIDLock.Lock()
	id = traceid.TraceID{
		High: uint64(seededIDGen.Int63()),
		Low:  uint64(seededIDGen.Int63()),
	}
	seededIDLock.Unlock()
	return
}

// randomTimestamped can generate 128 bit time sortable traceid's compatible
type randomTimestamped struct{}

func (t *randomTimestamped) TraceID() (id traceid.TraceID) {
	seededIDLock.Lock()
	id = traceid.TraceID{
		High: uint64(time.Now().Unix()<<32) + uint64(seededIDGen.Int31()),
		Low:  uint64(seededIDGen.Int63()),
	}
	seededIDLock.Unlock()
	return
}