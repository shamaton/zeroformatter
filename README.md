# zeroformatter

[![Build Status](https://travis-ci.org/shamaton/zeroformatter.svg?branch=master)](https://travis-ci.org/shamaton/zeroformatter)

golang version [zeroformatter](https://github.com/neuecc/ZeroFormatter)

## Usage
### Installation
```sh
go get github.com/shamaton/zeroformatter
```

### How to use
#### use simply
```go
package main;

import (
  "github.com/shamaton/zeroformatter"
  "log"
)

func main() {
	type Struct struct {
		String string
	}
	h := Struct{String: "zeroformatter"}

	d, err := zeroformatter.Serialize(h)
	if err != nil {
		log.Fatal(err)
	}
	r := Struct{}
	err = zeroformatter.Deserialize(&r, d)
	if err != nil {
		log.Fatal(err)
	}
}
```

#### delay
```go
package main;

import (
  "github.com/shamaton/zeroformatter"
  "log"
)

func how_to_use(b []byte) {
	type Struct struct {
		String string
	}
	
	r := Struct{}
	dds, _ := zeroformatter.DelayDeserialize(&r, b)
	
	// by element
	if err := dds.DeserializeByElement(&r.String); err != nil {
		log.Fatal(err)
	}
	
	// or by index
	if err := dds.DeserializeByIndex(0); err != nil {
	  log.Fatal(err)
	}
}
```

## Supported type 

### Primitive
| C# | Go |
| ---- | ---- |
| Int16 | int16 |
| Int32 | int32, int |
| Int64 | int64 |
| UInt16 | uint16 |
| UInt32 | uint32, uint |
| UInt64 | uint64 |
| Single | float32 |
| Double | float64 |
| Boolean | bool |
| Byte | uint8 |
| SByte | int8 |
| TimeSpan | time.Duration |
| DateTime | time.Time |
| String | string |

### Extension within golang
As these types can not convert with primitive type, I defined parent classes in golang.
These are only wrapping. please see codes.

| C# | Go |
| ---- | ---- |
| Char | zeroformatter.Char(rune) |
| DateTimeOffset | zeroformatter.DateTimeOffset(time.Time) |

### Array/Slice

| C# | Go |
| ---- | ---- |
| T[], List<T> | []T, [N]T |

### Map

| C# | Go |
| ---- | ---- |
| Dictionary<K, V> | map[K]V |

### Object

| C# | Go |
| ---- | ---- |
| Struct | struct |

## Not supported

`type?` is not supported, because golang doen't allow null in primitve types.

## License

This library is under the MIT License.
