package salticidae

// #include <stdlib.h>
// #include "salticidae/stream.h"
import "C"
import "runtime"

type byteArray struct {
    inner *C.bytearray_t
}
// Array of binary data.
type ByteArray = *byteArray

// Create an empty byte array (with zero contained bytes).
func NewByteArray() ByteArray {
    res := &byteArray{ inner: C.bytearray_new() }
    runtime.SetFinalizer(res, func(self ByteArray) { self.free() })
    return res
}

func (self ByteArray) free() { C.bytearray_free(self.inner) }

// Create a byte array by taking out all data from src. Notice this is a
// zero-copy operation that consumes and invalidates the data in src ("move"
// semantics) so that no more operation should be done to src after this
// function call.
func NewByteArrayMovedFromDataStream(src DataStream) ByteArray {
    res := &byteArray{ inner: C.bytearray_new_moved_from_datastream(src.inner) }
    runtime.SetFinalizer(res, func(self ByteArray) { self.free() })
    return res
}

type dataStream struct {
    inner *C.datastream_t
}

// Stream of binary data.
type DataStream = *dataStream

// Create an empty DataStream.
func NewDataStream() DataStream {
    res := &dataStream{ inner: C.datastream_new() }
    runtime.SetFinalizer(res, func(self DataStream) { self.free() })
    return res
}

// Create a DataStream with data copied from bytes.
func NewDataStreamFromBytes(bytes []byte) (res DataStream) {
    size := len(bytes)
    if size > 0 {
        base := (*C.uint8_t)(&bytes[0])
        res = &dataStream{ inner: C.datastream_new_from_bytes(base, C.size_t(size)) }
    } else {
        res = &dataStream{ inner: C.datastream_new() }
    }
    runtime.SetFinalizer(res, func(self DataStream) { self.free() })
    return
}

func (self DataStream) free() { C.datastream_free(self.inner) }

// Make a copy of the object.
func (self DataStream) Copy() DataStream {
    res := &dataStream{ inner: C.datastream_copy(self.inner) }
    runtime.SetFinalizer(res, func(self DataStream) { self.free() })
    return res
}

// TODO: datastream_data

// Empty the DataStream.
func (self DataStream) Clear() { C.datastream_clear(self.inner) }

func (self DataStream) Size() int { return int(C.datastream_size(self.inner)) }

// Write a uint8 integer to the stream (no byte order conversion).
func (self DataStream) PutU8(v uint8) bool { return bool(C.datastream_put_u8(self.inner, C.uint8_t(v))) }
// Write a uint16 integer to the stream (no byte order conversion).
func (self DataStream) PutU16(v uint16) bool { return bool(C.datastream_put_u16(self.inner, C.uint16_t(v))) }
// Write a uint32 integer to the stream (no byte order conversion).
func (self DataStream) PutU32(v uint32) bool { return bool(C.datastream_put_u32(self.inner, C.uint32_t(v))) }
// Write a uint64 integer to the stream (no byte order conversion).
func (self DataStream) PutU64(v uint64) bool { return bool(C.datastream_put_u64(self.inner, C.uint64_t(v))) }

// Write an int8 integer to the stream (no byte order conversion).
func (self DataStream) PutI8(v int8) bool { return bool(C.datastream_put_i8(self.inner, C.int8_t(v))) }
// Write an int16 integer to the stream (no byte order conversion).
func (self DataStream) PutI16(v int16) bool { return bool(C.datastream_put_i16(self.inner, C.int16_t(v))) }
// Write an int32 integer to the stream (no byte order conversion).
func (self DataStream) PutI32(v int32) bool { return bool(C.datastream_put_i32(self.inner, C.int32_t(v))) }
// Write an int64 integer to the stream (no byte order conversion).
func (self DataStream) PutI64(v int32) bool { return bool(C.datastream_put_i64(self.inner, C.int64_t(v))) }

// Write arbitrary bytes to the stream.
func (self DataStream) PutData(bytes []byte) bool {
    size := len(bytes)
    if size > 0 {
        base := (*C.uint8_t)(&bytes[0])
        return bool(C.datastream_put_data(self.inner, base, C.size_t(size)))
    } else { return true }
}

// Parse a uint8 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetU8(succ *bool) uint8 { return uint8(C.datastream_get_u8(self.inner, (*C.bool)(succ))) }
// Parse a uint16 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetU16(succ *bool) uint16 { return uint16(C.datastream_get_u16(self.inner, (*C.bool)(succ))) }
// Parse a uint32 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetU32(succ *bool) uint32 { return uint32(C.datastream_get_u32(self.inner, (*C.bool)(succ))) }
// Parse a uint64 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetU64(succ *bool) uint64 { return uint64(C.datastream_get_u64(self.inner, (*C.bool)(succ))) }

// Parse an int8 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetI8(succ *bool) int8 { return int8(C.datastream_get_i8(self.inner, (*C.bool)(succ))) }
// Parse an int16 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetI16(succ *bool) int16 { return int16(C.datastream_get_i16(self.inner, (*C.bool)(succ))) }
// Parse an int32 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetI32(succ *bool) int32 { return int32(C.datastream_get_i32(self.inner, (*C.bool)(succ))) }
// Parse an int64 integer by consuming the stream (no byte order conversion).
func (self DataStream) GetI64(succ *bool) int64 { return int64(C.datastream_get_i64(self.inner, (*C.bool)(succ))) }

// The handle returned by GetDataInPlace. The Go slice returned by Get() is
// valid only during the lifetime of the handle.
type dataStreamBytes struct {
    bytes []byte
    ds DataStream
}

type DataStreamBytes = *dataStreamBytes

func (self DataStreamBytes) Get() []byte { return self.bytes }

// Get the given length of preceeding bytes from the stream as a byte slice by
// consuming the stream. Notice this function does not copy the bytes, so the
// slice is only valid during the lifetime of DataStreamBytes handle.
func (self DataStream) GetDataInPlace(length int) DataStreamBytes {
    base := C.datastream_get_data_inplace(self.inner, C.size_t(length))
    return &dataStreamBytes{
        bytes: C.GoBytes(rawptr_t(base), C.int(length)),
        ds: self,
    }
}

type uint256 struct {
    inner *C.uint256_t
}

// 256-bit integer.
type UInt256 = *uint256

// Create a 256-bit integer.
func NewUInt256() UInt256 {
    res := &uint256{ inner: C.uint256_new() }
    runtime.SetFinalizer(res, func(self UInt256) { self.free() })
    return res
}

func (self UInt256) free() { C.uint256_free(self.inner) }

func (self UInt256) IsNull() bool { return bool(C.uint256_is_null(self.inner)) }

// Check if two 256-bit integers are equal.
func (self UInt256) IsEq(other UInt256) bool { return bool(C.uint256_is_eq(self.inner, other.inner)) }

// Write the integer to the given DataStream.
func (self UInt256) Serialize(s DataStream) { C.uint256_serialize(self.inner, s.inner) }

// Parse the integer from the given DataStream.
func (self UInt256) Unserialize(s DataStream) { C.uint256_unserialize(self.inner, s.inner) }

// Get the Sha256 hash of the given DataStream content (without consuming the
// stream).
func (self DataStream) GetHash() UInt256 {
    res := &uint256{ inner: C.datastream_get_hash(self.inner) }
    runtime.SetFinalizer(res, func(self UInt256) { self.free() })
    return res
}

// Get hexadicemal string representation of the given DataStream content
// (without consuming the stream).
func (self DataStream) GetHex() string {
    tmp := C.datastream_get_hex(self.inner)
    res := C.GoString(tmp)
    C.free(rawptr_t(tmp))
    return res
}

// Get hexadicemal string representation of the 256-bit integer.
func (self UInt256) GetHex() string {
    s := NewDataStream()
    self.Serialize(s)
    res := s.GetHex()
    return res
}
