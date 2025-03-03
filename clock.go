package vectorclock

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
)

const maxNumberVersions = math.MaxUint16

type VectorClock struct {
	SerialVersionID int
	versionMap      map[uint16]uint64
	timestamp       int64
}

func (v *VectorClock) String() string {
	buf := bytes.Buffer{}

	for k, v := range v.versionMap {
		buf.WriteString(fmt.Sprintf("K=%d,V=%d\n", k, v))
	}
	buf.WriteString(fmt.Sprintf("timestamp: %d\n", v.timestamp))
	return buf.String()
}

func (v *VectorClock) Clone() *VectorClock {
	clone := &VectorClock{SerialVersionID: 1,
		timestamp:  v.timestamp,
		versionMap: make(map[uint16]uint64, len(v.versionMap)),
	}
	clone.SerialVersionID = 1
	clone.timestamp = v.timestamp

	for k, v := range v.versionMap {
		clone.versionMap[k] = v
	}
	return clone
}

func (v *VectorClock) Compare(other *VectorClock) (Occurred, error) {
	if v == nil || other == nil {
		return BEFORE, errors.New("cant compare null clocks")
	}
	vBigger := false
	otherBigger := false

	for k, versionV := range v.versionMap {
		if versionOther, exists := other.versionMap[k]; exists {
			if versionV > versionOther {
				vBigger = true
			} else if versionV < versionOther {
				otherBigger = true
			}
		} else {
			vBigger = true
		}
	}

	for k := range other.versionMap {
		if _, exists := v.versionMap[k]; !exists {
			otherBigger = true
		}
	}

	if !vBigger && otherBigger {
		return BEFORE, nil
	}
	if vBigger && !otherBigger {
		return AFTER, nil
	}
	if vBigger && otherBigger {
		return CONCURRENTLY, nil
	}
	return BEFORE, nil // Equal case
}

func (v *VectorClock) Merge(clock *VectorClock) *VectorClock {
	newClock := NewEmptyClock()
	for key, val := range v.versionMap {
		newClock.versionMap[key] = val
	}
	for key, val := range clock.versionMap {
		currentVal, ok := newClock.versionMap[key]
		if !ok {
			newClock.versionMap[key] = val
		} else {
			newClock.versionMap[key] = max(val, currentVal)
		}
	}
	return newClock
}

func NewEmptyClock() *VectorClock {
	return &VectorClock{
		SerialVersionID: 1,
		versionMap:      make(map[uint16]uint64),
		timestamp:       time.Now().UnixMilli(),
	}
}

func (v *VectorClock) IncrementVersion(node int, timestampMillis int64) error {

	if node < 0 || node > math.MaxUint16 {
		return errors.New("outside range of acceptable node ids")
	}
	v.timestamp = timestampMillis

	version, ok := v.versionMap[uint16(node)]
	if !ok {
		version = 1
	} else {
		version += 1
	}
	v.versionMap[uint16(node)] = version

	if len(v.versionMap) >= maxNumberVersions {
		return errors.New("vector clock is full")
	}
	return nil
}

func (v *VectorClock) GetMaxVersion() uint64 {
	var maxVersion uint64
	for _, value := range v.versionMap {
		maxVersion = maxUint64(maxVersion, value)
	}
	return maxVersion
}

type ClockEntry struct {
	Key   uint16
	Value uint64
}

func (v *VectorClock) GetEntries() ([]*ClockEntry, error) {
	var result []*ClockEntry

	for key, val := range v.versionMap {
		clockEntry := &ClockEntry{
			Key:   key,
			Value: val,
		}
		result = append(result, clockEntry)
	}
	return result, nil
}

func (v *VectorClock) CopyFromVectorClock(vc *VectorClock) error {
	v.versionMap = map[uint16]uint64{}
	v.timestamp = vc.timestamp
	entries, err := v.GetEntries()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		v.versionMap[entry.Key] = entry.Value
	}
	return nil
}

func (v *VectorClock) Incremented(nodeId int, timeMillis int64) (*VectorClock, error) {

	outputVersionMap := make(map[uint16]uint64)
	for k, v := range v.versionMap {
		outputVersionMap[k] = v
	}

	newVc := &VectorClock{
		SerialVersionID: 1,
		versionMap:      outputVersionMap,
		timestamp:       v.timestamp,
	}
	if err := newVc.IncrementVersion(nodeId, timeMillis); err != nil {
		return nil, err
	}
	return newVc, nil
}

// TODO - Introduce a genric set type?
func (v *VectorClock) keySet() map[uint16]interface{} {
	result := map[uint16]interface{}{}

	for k, _ := range v.versionMap {
		result[k] = nil
	}
	return result
}

func uint16ToBytes(value uint16) []byte {
	short := make([]byte, 2)
	binary.BigEndian.PutUint16(short, value)
	return short
}

func uint64ToBytes(value uint64) []byte {
	long := make([]byte, 8)
	binary.BigEndian.PutUint64(long, value)
	return long
}

func readUint16(contents []byte) uint16 {
	return binary.BigEndian.Uint16(contents)
}

func readUint64(contents []byte) uint64 {
	return binary.BigEndian.Uint64(contents)
}

func VectorClockFromBytes(contents []byte) *VectorClock {

	vc := &VectorClock{
		SerialVersionID: 1,
	}

	size := readUint16(contents[0:2])
	vc.versionMap = make(map[uint16]uint64, size)

	timestamp := readUint64(contents[len(contents)-8:])
	vc.timestamp = int64(timestamp)

	offset := 2

	for offset < len(contents)-8 {
		key := readUint16(contents[offset:])
		offset += 2
		value := readUint64(contents[offset:])
		offset += 8
		vc.versionMap[key] = value
	}
	return vc
}

func (v *VectorClock) ToBytes() []byte {
	buffer := bytes.Buffer{}
	// Write the size
	buffer.Write(uint16ToBytes(uint16(len(v.versionMap))))
	for k, v := range v.versionMap {
		buffer.Write(uint16ToBytes(k))
		buffer.Write(uint64ToBytes(v))
	}
	buffer.Write(uint64ToBytes(uint64(v.timestamp)))
	return buffer.Bytes()
}

func copySet(toCopy map[uint16]interface{}) map[uint16]interface{} {
	result := make(map[uint16]interface{}, len(toCopy))
	for k, _ := range toCopy {
		result[k] = nil
	}
	return result
}

func retainSet(original map[uint16]interface{}, other map[uint16]interface{}) map[uint16]interface{} {
	output := map[uint16]interface{}{}

	for k, _ := range other {
		_, ok := original[k]
		if ok {
			output[k] = nil
		}
	}
	return output
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
