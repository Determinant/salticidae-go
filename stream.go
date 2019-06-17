package salticidae

// #include <stdlib.h>
// #include "salticidae/stream.h"
import "C"
import "runtime"

type byteArray struct {
    inner *C.struct_bytearray_t
}
type ByteArray = *byteArray

func NewByteArray() ByteArray {
    res := &byteArray{ inner: C.bytearray_new() }
    runtime.SetFinalizer(res, func(self ByteArray) { self.free() })
    return res
}

func (self ByteArray) free() { C.bytearray_free(self.inner) }

func NewByteArrayMovedFromDataStream(src DataStream) ByteArray {
    res := &byteArray{ inner: C.bytearray_new_moved_from_datastream(src.inner) }
    runtime.SetFinalizer(res, func(self ByteArray) { self.free() })
    return res
}

type dataStream struct {
    inner *C.struct_datastream_t
}

type DataStream = *dataStream

func NewDataStream() DataStream {
    res := &dataStream{ inner: C.datastream_new() }
    runtime.SetFinalizer(res, func(self DataStream) { self.free() })
    return res
}

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

func (self DataStream) AssignByCopy(src DataStream) {
    C.datastream_assign_by_copy(self.inner, src.inner)
}

func (self DataStream) AssignByMove(src DataStream) {
    C.datastream_assign_by_move(self.inner, src.inner)
}

// TODO: datastream_data

func (self DataStream) Clear() { C.datastream_clear(self.inner) }

func (self DataStream) Size() int { return int(C.datastream_size(self.inner)) }

func (self DataStream) PutU8(v uint8) bool { return bool(C.datastream_put_u8(self.inner, C.uint8_t(v))) }
func (self DataStream) PutU16(v uint16) bool { return bool(C.datastream_put_u16(self.inner, C.uint16_t(v))) }
func (self DataStream) PutU32(v uint32) bool { return bool(C.datastream_put_u32(self.inner, C.uint32_t(v))) }
func (self DataStream) PutU64(v uint32) bool { return bool(C.datastream_put_u64(self.inner, C.uint64_t(v))) }

func (self DataStream) PutI8(v int8) bool { return bool(C.datastream_put_i8(self.inner, C.int8_t(v))) }
func (self DataStream) PutI16(v int16) bool { return bool(C.datastream_put_i16(self.inner, C.int16_t(v))) }
func (self DataStream) PutI32(v int32) bool { return bool(C.datastream_put_i32(self.inner, C.int32_t(v))) }
func (self DataStream) PutI64(v int32) bool { return bool(C.datastream_put_i64(self.inner, C.int64_t(v))) }

func (self DataStream) PutData(bytes []byte) bool {
    size := len(bytes)
    if size > 0 {
        base := (*C.uint8_t)(&bytes[0])
        return bool(C.datastream_put_data(self.inner, base, C.size_t(size)))
    } else { return true }
}

func (self DataStream) GetU8(succ *bool) uint8 { return uint8(C.datastream_get_u8(self.inner, (*C.bool)(succ))) }
func (self DataStream) GetU16(succ *bool) uint16 { return uint16(C.datastream_get_u16(self.inner, (*C.bool)(succ))) }
func (self DataStream) GetU32(succ *bool) uint32 { return uint32(C.datastream_get_u32(self.inner, (*C.bool)(succ))) }
func (self DataStream) GetU64(succ *bool) uint64 { return uint64(C.datastream_get_u64(self.inner, (*C.bool)(succ))) }

func (self DataStream) GetI8(succ *bool) int8 { return int8(C.datastream_get_i8(self.inner, (*C.bool)(succ))) }
func (self DataStream) GetI16(succ *bool) int16 { return int16(C.datastream_get_i16(self.inner, (*C.bool)(succ))) }
func (self DataStream) GetI32(succ *bool) int32 { return int32(C.datastream_get_i32(self.inner, (*C.bool)(succ))) }
func (self DataStream) GetI64(succ *bool) int64 { return int64(C.datastream_get_i64(self.inner, (*C.bool)(succ))) }


func (self DataStream) GetDataInPlace(length int) []byte {
    base := C.datastream_get_data_inplace(self.inner, C.size_t(length))
    return C.GoBytes(rawptr_t(base), C.int(length))
}

type uint256 struct {
    inner *C.uint256_t
}

type UInt256 = *uint256

func NewUInt256() UInt256 {
    res := &uint256{ inner: C.uint256_new() }
    runtime.SetFinalizer(res, func(self UInt256) { self.free() })
    return res
}

func (self UInt256) free() { C.uint256_free(self.inner) }
func (self UInt256) UInt256IsNull() bool { return bool(C.uint256_is_null(self.inner)) }
func (self UInt256) UInt256IsEq(other UInt256) bool { return bool(C.uint256_is_eq(self.inner, other.inner)) }
func (self UInt256) Serialize(s DataStream) { C.uint256_serialize(self.inner, s.inner) }
func (self UInt256) Unserialize(s DataStream) { C.uint256_unserialize(self.inner, s.inner) }
func (self UInt256) IsEq(other UInt256) bool { return bool(C.uint256_is_eq(self.inner, other.inner)) }

func (self DataStream) GetHash() UInt256 {
    res := &uint256{ inner: C.datastream_get_hash(self.inner) }
    runtime.SetFinalizer(res, func(self UInt256) { self.free() })
    return res
}

func (self DataStream) GetHex() string {
    tmp := C.datastream_get_hex(self.inner)
    res := C.GoString(tmp)
    C.free(rawptr_t(tmp))
    return res
}

func (self UInt256) GetHex() string {
    s := NewDataStream()
    self.Serialize(s)
    res := s.GetHex()
    return res
}
