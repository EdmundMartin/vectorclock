package vectorclock

type VersionedBytes struct {
	Clock    *VectorClock
	Contents []byte
}

func NewVersionedBytes(contents []byte, version *VectorClock) *VersionedBytes {
	if version == nil {
		version = NewEmptyClock()
	}
	return &VersionedBytes{
		Clock:    version,
		Contents: contents,
	}
}

func (v *VersionedBytes) ToBytes() []byte {
	clockBytes := v.Clock.ToBytes()
	contentSize := len(v.Contents)
	totalSize := 2 + len(clockBytes) + contentSize
	result := make([]byte, totalSize)
	copy(result, uint16ToBytes(uint16(contentSize)))
	copy(result[2:], v.Contents)
	copy(result[2+contentSize:], clockBytes)
	return result
}

func VersionedBytesFromBytes(contents []byte) *VersionedBytes {
	v := &VersionedBytes{}
	sizeContents := readUint16(contents)
	v.Contents = contents[2 : 2+sizeContents]
	v.Clock = VectorClockFromBytes(contents[2+sizeContents:])
	return v
}

func (v *VersionedBytes) HappenedBefore(other *VersionedBytes) (int, error) {
	result, err := v.Clock.Compare(other.Clock)
	if err != nil {
		return 0, err
	}
	if result == BEFORE {
		return -1, nil
	}
	if result == AFTER {
		return 1, nil
	}
	return 0, nil
}

type VersionedBytesCollection []*VersionedBytes

func (v VersionedBytesCollection) Len() int {
	return len(v)
}

func (v VersionedBytesCollection) Less(i, j int) bool {
	res, _ := v[i].HappenedBefore(v[j])
	return res < 0
}

func (v VersionedBytesCollection) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
