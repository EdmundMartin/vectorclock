package vectorclock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewVersionedBytes(t *testing.T) {

	contents := []byte("HelloWorld")

	versionedContents := NewVersionedBytes(contents, nil)

	nowMilli := time.Now().UnixMilli()
	err := versionedContents.Clock.IncrementVersion(1, nowMilli)
	if err != nil {
		t.Error(err)
	}
	versionAsBytes := versionedContents.ToBytes()
	res := VersionedBytesFromBytes(versionAsBytes)

	assert.Equal(t, res.Clock, versionedContents.Clock)
	assert.Equal(t, res.Contents, versionedContents.Contents)
}
