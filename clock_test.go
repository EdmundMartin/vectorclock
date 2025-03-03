package vectorclock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getClock(nodes ...int) *VectorClock {
	vectorClock := NewEmptyClock()
	increment(vectorClock, nodes...)
	return vectorClock
}

func increment(clock *VectorClock, nodes ...int) {
	for _, n := range nodes {
		err := clock.IncrementVersion(n, clock.timestamp)
		if err != nil {
			return
		}
	}
	return
}

func TestVectorClock_Compare(t *testing.T) {

	res, err := getClock().Compare(getClock())
	assert.NoError(t, err)
	assert.NotEqual(t, res, CONCURRENTLY)

	res, err = getClock(1, 1, 2).Compare(getClock(1, 1, 2))
	assert.NoError(t, err)
	assert.NotEqual(t, res, CONCURRENTLY)

	res, err = getClock(1, 1, 2).Compare(getClock(1, 1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, res, BEFORE)

	// Clocks with different events should be concurrent.
	res, err = getClock(1).Compare(getClock(2))
	assert.NoError(t, err)
	assert.Equal(t, res, CONCURRENTLY)

	// Clocks with different events should be concurrent
	res, err = getClock(1, 2, 3, 3).Compare(getClock(1, 1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, res, AFTER)

	res, err = getClock(2, 2).Compare(getClock(1, 2, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, res, BEFORE)

	res, err = getClock(1, 2, 2, 3).Compare(getClock(2, 2))
	assert.NoError(t, err)
	assert.Equal(t, res, AFTER)
}

func TestVectorClock_CompareVersion(t *testing.T) {

	sameTime := time.Now().Add(-time.Minute).UnixMilli()
	clockOne := &VectorClock{
		SerialVersionID: 1,
		versionMap:      map[uint16]uint64{1: 10},
		timestamp:       sameTime,
	}

	clockTwo := &VectorClock{
		SerialVersionID: 1,
		versionMap:      map[uint16]uint64{1: 10},
		timestamp:       sameTime,
	}

	res, err := clockOne.Compare(clockTwo)
	assert.NoError(t, err)
	assert.Equal(t, BEFORE, res)

	err = clockOne.IncrementVersion(1, time.Now().UnixMilli())
	assert.NoError(t, err)

	res, err = clockOne.Compare(clockTwo)
	assert.NoError(t, err)
	assert.Equal(t, AFTER, res)
}

func TestVectorClock_ToBytes(t *testing.T) {

	clock := getClock(1, 2, 3, 4, 4, 1, 1, 1)
	recoveredClock := VectorClockFromBytes(clock.ToBytes())
	assert.Equal(t, clock.timestamp, recoveredClock.timestamp)
	assert.Equal(t, clock.versionMap, recoveredClock.versionMap)
}
