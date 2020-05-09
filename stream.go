package salticidae

// #include <stdlib.h>
// #include "salticidae/stream.h"
// #include "salticidae/endian.h"
// uint16_t _salti_htole16(uint16_t x) { return htole16(x); }
// uint32_t _salti_htole32(uint32_t x) { return htole32(x); }
// uint64_t _salti_htole64(uint64_t x) { return htole64(x); }
//
// uint16_t _salti_letoh16(uint16_t x) { return le16toh(x); }
// uint32_t _salti_letoh32(uint32_t x) { return le32toh(x); }
// uint64_t _salti_letoh64(uint64_t x) { return le64toh(x); }
//
// uint16_t _salti_htobe16(uint16_t x) { return htobe16(x); }
// uint32_t _salti_htobe32(uint32_t x) { return htobe32(x); }
// uint64_t _salti_htobe64(uint64_t x) { return htobe64(x); }
//
// uint16_t _salti_betoh16(uint16_t x) { return be16toh(x); }
// uint32_t _salti_betoh32(uint32_t x) { return be32toh(x); }
// uint64_t _salti_betoh64(uint64_t x) { return be64toh(x); }
//
import "C"
import "runtime"

//// begin ByteArray def

/// CByteArray is the C pointer type for a ByteArray object.
type CByteArray = *C.bytearray_t
type byteArray struct {
	inner    CByteArray
	autoFree bool
}

// ByteArray is an array of binary data.
type ByteArray = *byteArray

// ByteArrayFromC converts an existing C pointer into a go pointer.
func ByteArrayFromC(ptr CByteArray) ByteArray {
	return &byteArray{inner: ptr}
}

func byteArraySetFinalizer(res ByteArray, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self ByteArray) { self.Free() })
	}
}

// Free the ByteArray manually. If the object is constructed with autoFree =
// true, this will immediately free the object.
func (ba ByteArray) Free() {
	C.bytearray_free(ba.inner)
	if ba.autoFree {
		runtime.SetFinalizer(ba, nil)
	}
}

//// end ByteArray def

//// begin DataStream def

// The C pointer to a DataStream object.
type CDataStream = *C.datastream_t
type dataStream struct {
	inner    CDataStream
	attached map[uintptr]interface{}
	autoFree bool
}

// Stream of binary data.
type DataStream = *dataStream

func DataStreamFromC(ptr CDataStream) DataStream {
	return &dataStream{
		inner:    ptr,
		attached: make(map[uintptr]interface{}),
	}
}

func dataStreamSetFinalizer(res DataStream, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self DataStream) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (ds DataStream) Free() {
	C.datastream_free(ds.inner)
	if ds.autoFree {
		runtime.SetFinalizer(ds, nil)
	}
}

func (ds DataStream) attach(ptr RawPtr, obj interface{}) { ds.attached[uintptr(ptr)] = obj }
func (ds DataStream) detach(ptr RawPtr)                  { delete(ds.attached, uintptr(ptr)) }

//// end DataStream def

//// begin UInt256 def

// CUInt256 is the C pointer to a UInt256 object.
type CUInt256 = *C.uint256_t
type uint256 struct {
	inner    CUInt256
	autoFree bool
}

// UInt256 is a 256-bit integer.
type UInt256 = *uint256

// UInt256FromC converts an existing C pointer into a go pointer.
func UInt256FromC(ptr CUInt256) UInt256 {
	return &uint256{inner: ptr}
}

func uint256SetFinalizer(res UInt256, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self UInt256) { self.Free() })
	}
}

// Free the underlying C pointer manually.
func (u256 UInt256) Free() {
	C.uint256_free(u256.inner)
	if u256.autoFree {
		runtime.SetFinalizer(u256, nil)
	}
}

//// end UInt256 def

//// begin ByteArray methods

// Create an empty byte array (with zero contained bytes).
func NewByteArray(autoFree bool) ByteArray {
	res := ByteArrayFromC(C.bytearray_new())
	byteArraySetFinalizer(res, autoFree)
	return res
}

// NewByteArrayMovedFromDataStream creates a byte array by taking out all data
// from src. Notice this is a zero-copy operation that consumes and invalidates
// the data in src ("move" semantics) so that no more operation should be done
// to src after this function call. Also notice unlike copying, the entire
// DataStream buffer is moved (including the possibily consumed part).
func NewByteArrayMovedFromDataStream(src DataStream, autoFree bool) (res ByteArray) {
	res = ByteArrayFromC(C.bytearray_new_moved_from_datastream(src.inner))
	byteArraySetFinalizer(res, autoFree)
	return
}

// NewByteArrayCopiedFromDataStream creates a byte array by copying the
// remaining data from src.
func NewByteArrayCopiedFromDataStream(src DataStream, autoFree bool) (res ByteArray) {
	res = ByteArrayFromC(C.bytearray_new_copied_from_datastream(src.inner))
	byteArraySetFinalizer(res, autoFree)
	return
}

// NewByteArrayFromHex creates a byte array from the hex string.
func NewByteArrayFromHex(hex string) (res ByteArray, autoFree bool) {
	cStr := C.CString(hex)
	res = ByteArrayFromC(C.bytearray_new_from_hex(cStr))
	C.free(RawPtr(cStr))
	byteArraySetFinalizer(res, autoFree)
	return
}

// NewByteArrayFromBytes creates a byte array from a byte slice.
func NewByteArrayFromBytes(bytes []byte, autoFree bool) (res ByteArray) {
	size := len(bytes)
	if size > 0 {
		base := (*C.uint8_t)(&bytes[0])
		res = ByteArrayFromC(C.bytearray_new_from_bytes(base, C.size_t(size)))
	} else {
		res = ByteArrayFromC(C.bytearray_new())
	}
	byteArraySetFinalizer(res, autoFree)
	return
}

// GetHash gets the Sha256 hash of the given ByteArray content.
func (ba ByteArray) GetHash(autoFree bool) (res UInt256) {
	s := NewDataStreamFromByteArray(ba, false)
	res = s.GetHash(autoFree)
	s.Free()
	return
}

//// end ByteArray methods

//// begin DataStream methods

// NewDataStream creates an empty DataStream.
func NewDataStream(autoFree bool) DataStream {
	res := DataStreamFromC(C.datastream_new())
	dataStreamSetFinalizer(res, autoFree)
	return res
}

// NewDataStreamFromBytes creates a DataStream with data copied from bytes.
func NewDataStreamFromBytes(bytes []byte, autoFree bool) (res DataStream) {
	size := len(bytes)
	if size > 0 {
		base := (*C.uint8_t)(&bytes[0])
		res = DataStreamFromC(C.datastream_new_from_bytes(base, C.size_t(size)))
	} else {
		res = DataStreamFromC(C.datastream_new())
	}
	dataStreamSetFinalizer(res, autoFree)
	return
}

// NewDataStreamMovedFromByteArray creates a DataStream with content moved from
// a ByteArray.
func NewDataStreamMovedFromByteArray(bytes ByteArray, autoFree bool) (res DataStream) {
	res = DataStreamFromC(C.datastream_new_moved_from_bytearray(bytes.inner))
	dataStreamSetFinalizer(res, autoFree)
	return
}

// NewDataStreamFromByteArray creates a DataStream with content copied from a
// ByteArray.
func NewDataStreamFromByteArray(bytes ByteArray, autoFree bool) (res DataStream) {
	res = DataStreamFromC(C.datastream_new_from_bytearray(bytes.inner))
	dataStreamSetFinalizer(res, autoFree)
	return
}

// Copy the object.
func (ds DataStream) Copy(autoFree bool) (res DataStream) {
	res = DataStreamFromC(C.datastream_copy(ds.inner))
	dataStreamSetFinalizer(res, autoFree)
	return
}

// GetHash gets the Sha256 hash of the given DataStream content (without
// consuming the stream).
func (ds DataStream) GetHash(autoFree bool) UInt256 {
	res := UInt256FromC(C.datastream_get_hash(ds.inner))
	uint256SetFinalizer(res, autoFree)
	return res
}

// GetHex gets hexadicemal string representation of the given DataStream
// content (without consuming the stream).
func (ds DataStream) GetHex() string {
	tmp := C.datastream_get_hex(ds.inner)
	res := C.GoString(tmp)
	C.free(RawPtr(tmp))
	return res
}

// TODO: datastream_data

// Clear the DataStream.
func (ds DataStream) Clear() { C.datastream_clear(ds.inner) }

// Size returns the size of the DataStream.
func (ds DataStream) Size() int { return int(C.datastream_size(ds.inner)) }

// PutU8 writes a uint8 integer to the stream (no byte order conversion).
func (ds DataStream) PutU8(v uint8) bool {
	return bool(C.datastream_put_u8(ds.inner, C.uint8_t(v)))
}

// PutU16 writes a uint16 integer to the stream (no byte order conversion).
func (ds DataStream) PutU16(v uint16) bool {
	return bool(C.datastream_put_u16(ds.inner, C.uint16_t(v)))
}

// PutU32 writes a uint32 integer to the stream (no byte order conversion).
func (ds DataStream) PutU32(v uint32) bool {
	return bool(C.datastream_put_u32(ds.inner, C.uint32_t(v)))
}

// PutU64 writes a uint64 integer to the stream (no byte order conversion).
func (ds DataStream) PutU64(v uint64) bool {
	return bool(C.datastream_put_u64(ds.inner, C.uint64_t(v)))
}

// PutI8 writes an int8 integer to the stream (no byte order conversion).
func (ds DataStream) PutI8(v int8) bool {
	return bool(C.datastream_put_i8(ds.inner, C.int8_t(v)))
}

// PutI16 writes an int16 integer to the stream (no byte order conversion).
func (ds DataStream) PutI16(v int16) bool {
	return bool(C.datastream_put_i16(ds.inner, C.int16_t(v)))
}

// PutI32 writes an int32 integer to the stream (no byte order conversion).
func (ds DataStream) PutI32(v int32) bool {
	return bool(C.datastream_put_i32(ds.inner, C.int32_t(v)))
}

// PutI64 writes an int64 integer to the stream (no byte order conversion).
func (ds DataStream) PutI64(v int32) bool {
	return bool(C.datastream_put_i64(ds.inner, C.int64_t(v)))
}

// PutData writes arbitrary bytes to the stream.
func (ds DataStream) PutData(bytes []byte) bool {
	size := len(bytes)
	if size > 0 {
		base := (*C.uint8_t)(&bytes[0])
		return bool(C.datastream_put_data(ds.inner, base, C.size_t(size)))
	}
	return true
}

// GetU8 parses a uint8 integer by consuming the stream (no byte order
// conversion).
func (ds DataStream) GetU8(succ *bool) uint8 {
	return uint8(C.datastream_get_u8(ds.inner, (*C.bool)(succ)))
}

// GetU16 parses a uint16 integer by consuming the stream (no byte order
// conversion).
func (ds DataStream) GetU16(succ *bool) uint16 {
	return uint16(C.datastream_get_u16(ds.inner, (*C.bool)(succ)))
}

// GetU32 parses a uint32 integer by consuming the stream (no byte order
// conversion).
func (ds DataStream) GetU32(succ *bool) uint32 {
	return uint32(C.datastream_get_u32(ds.inner, (*C.bool)(succ)))
}

// GetU64 parses a uint64 integer by consuming the stream (no byte order
// conversion).
func (ds DataStream) GetU64(succ *bool) uint64 {
	return uint64(C.datastream_get_u64(ds.inner, (*C.bool)(succ)))
}

// GetI8 parses an int8 integer by consuming the stream (no byte order
// conversion).
func (ds DataStream) GetI8(succ *bool) int8 {
	return int8(C.datastream_get_i8(ds.inner, (*C.bool)(succ)))
}

// GetI16 parses an int16 integer by consuming the stream (no byte order
// conversion).
func (ds DataStream) GetI16(succ *bool) int16 {
	return int16(C.datastream_get_i16(ds.inner, (*C.bool)(succ)))
}

// GetI32 parses an int32 integer by consuming the stream (no byte order conversion).
func (ds DataStream) GetI32(succ *bool) int32 {
	return int32(C.datastream_get_i32(ds.inner, (*C.bool)(succ)))
}

// GetI64 parses an int64 integer by consuming the stream (no byte order conversion).
func (ds DataStream) GetI64(succ *bool) int64 {
	return int64(C.datastream_get_i64(ds.inner, (*C.bool)(succ)))
}

type dataStreamBytes struct {
	bytes []byte
	ds    DataStream
}

// DataStreamBytes is the handle returned by GetDataInPlace. The Go slice
// returned by Get() is valid only during the lifetime of the handle.
type DataStreamBytes = *dataStreamBytes

// Get the underlying bytes.
func (dsb DataStreamBytes) Get() []byte { return dsb.bytes }

// Release the handle.
func (dsb DataStreamBytes) Release() { dsb.ds.detach(RawPtr(dsb)) }

// GetDataInPlace gets the given length of preceeding bytes from the stream as
// a byte slice by consuming the stream. Notice this function does not copy the
// bytes, so the slice is only valid during the lifetime of DataStreamBytes
// handle.
func (ds DataStream) GetDataInPlace(length int) DataStreamBytes {
	base := C.datastream_get_data_inplace(ds.inner, C.size_t(length))
	res := &dataStreamBytes{
		bytes: C.GoBytes(RawPtr(base), C.int(length)),
		ds:    ds,
	}
	ds.attach(RawPtr(res), res)
	return res
}

//// end DataStream methods

//// begin UInt256 methods

// NewUInt256 creates a 256-bit integer.
func NewUInt256(autoFree bool) UInt256 {
	res := UInt256FromC(C.uint256_new())
	uint256SetFinalizer(res, autoFree)
	return res
}

// NewUInt256FromByteArray creates a 256-bit from the ByteArray.
func NewUInt256FromByteArray(bytes ByteArray, autoFree bool) (res UInt256) {
	res = UInt256FromC(C.uint256_new_from_bytearray(bytes.inner))
	uint256SetFinalizer(res, autoFree)
	return
}

// IsNull checks if the UInt256 is empty.
func (u256 UInt256) IsNull() bool { return bool(C.uint256_is_null(u256.inner)) }

// IsEq checks if two 256-bit integers are equal.
func (u256 UInt256) IsEq(other UInt256) bool { return bool(C.uint256_is_eq(u256.inner, other.inner)) }

// Serialize writes the integer to the given DataStream.
func (u256 UInt256) Serialize(s DataStream) { C.uint256_serialize(u256.inner, s.inner) }

// Unserialize parses the integer from the given DataStream.
func (u256 UInt256) Unserialize(s DataStream) { C.uint256_unserialize(u256.inner, s.inner) }

// GetHex gets hexadicemal string representation of the 256-bit integer.
func (u256 UInt256) GetHex() (res string) {
	s := NewDataStream(false)
	u256.Serialize(s)
	res = s.GetHex()
	s.Free()
	return
}

/// end UInt256 methods

func ToLittleEndianU16(x uint16) uint16 { return uint16(C._salti_htole16(C.uint16_t(x))) }
func ToLittleEndianU32(x uint32) uint32 { return uint32(C._salti_htole32(C.uint32_t(x))) }
func ToLittleEndianU64(x uint64) uint64 { return uint64(C._salti_htole64(C.uint64_t(x))) }

func FromLittleEndianU16(x uint16) uint16 { return uint16(C._salti_letoh16(C.uint16_t(x))) }
func FromLittleEndianU32(x uint32) uint32 { return uint32(C._salti_letoh32(C.uint32_t(x))) }
func FromLittleEndianU64(x uint64) uint64 { return uint64(C._salti_letoh64(C.uint64_t(x))) }

func ToBigEndianU16(x uint16) uint16 { return uint16(C._salti_htobe16(C.uint16_t(x))) }
func ToBigEndianU32(x uint32) uint32 { return uint32(C._salti_htobe32(C.uint32_t(x))) }
func ToBigEndianU64(x uint64) uint64 { return uint64(C._salti_htobe64(C.uint64_t(x))) }

func FromBigEndianU16(x uint16) uint16 { return uint16(C._salti_betoh16(C.uint16_t(x))) }
func FromBigEndianU32(x uint32) uint32 { return uint32(C._salti_betoh32(C.uint32_t(x))) }
func FromBigEndianU64(x uint64) uint64 { return uint64(C._salti_betoh64(C.uint64_t(x))) }
