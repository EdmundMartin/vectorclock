# VectorClock

## Overview

The `vectorclock` package provides an implementation of vector clocks for tracking causality in distributed systems. It includes a `VectorClock` structure to manage versioning of events across multiple nodes and a `VersionedBytes` structure to store versioned data.

## Features

- **VectorClock**: Implements vector clocks to track causality.
- **VersionedBytes**: Associates a vector clock with byte content for versioned storage.
- **Comparison Operations**: Determine if one vector clock happened before, after, or concurrently with another.
- **Serialization and Deserialization**: Convert vector clocks and versioned bytes to and from byte arrays.

## Installation

To use this package in your Go project:

```sh
 go get github.com/EdmundMartin/vectorclock
```

## Usage

### Creating and Manipulating a VectorClock

```go
import "github.com/EdmundMartin/vectorclock"

vc := vectorclock.NewEmptyClock()
vc.IncrementVersion(1, time.Now().UnixMilli())
vc.IncrementVersion(2, time.Now().UnixMilli())
```

### Comparing Vector Clocks

```go
vc1 := vectorclock.NewEmptyClock()
vc1.IncrementVersion(1, time.Now().UnixMilli())

vc2 := vectorclock.NewEmptyClock()
vc2.IncrementVersion(2, time.Now().UnixMilli())

result, _ := vc1.Compare(vc2)
if result == vectorclock.CONCURRENTLY {
    fmt.Println("vc1 and vc2 are concurrent")
}
```

### Using VersionedBytes

```go
contents := []byte("Hello, World!")
vc := vectorclock.NewEmptyClock()
versionedData := vectorclock.NewVersionedBytes(contents, vc)

serialized := versionedData.ToBytes()
restored := vectorclock.VersionedBytesFromBytes(serialized)
```

## License

This project is licensed under the MIT License.

## Contributing

Feel free to submit issues and pull requests to improve the package.

