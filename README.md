# zeroformatter

[![Build Status](https://travis-ci.org/shamaton/zeroformatter.svg?branch=master)](https://travis-ci.org/shamaton/zeroformatter)

golang version zeroformatter

## Usage
### Installation
```sh
go get github.com/shamaton/zeroformatter
```

### How to use
under construction...

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
