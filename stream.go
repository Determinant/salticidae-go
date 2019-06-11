package salticidae

// #include <stdlib.h>
// #include "salticidae/stream.h"
import "C"

type ByteArray = *C.struct_bytearray_t

func NewByteArray() ByteArray { return C.bytearray_new() }
func (self ByteArray) Free() { C.bytearray_free(self) }

func NewByteArrayMovedFromDataStream(src DataStream) ByteArray {
    return C.bytearray_new_moved_from_datastream(src)
}

type DataStream = *C.struct_datastream_t

func NewDataStream() DataStream { return C.datastream_new() }
func NewDataStreamFromBytes(bytes []byte) DataStream {
    size := len(bytes)
    if size > 0 {
        base := (*C.uint8_t)(&bytes[0])
        return C.datastream_new_from_bytes(base, C.size_t(size))
    } else {
        return C.datastream_new()
    }
}

func (self DataStream) Free() { C.datastream_free(self) }

func (self DataStream) AssignByCopy(src DataStream) {
    C.datastream_assign_by_copy(self, src)
}

func (self DataStream) AssignByMove(src DataStream) {
    C.datastream_assign_by_move(self, src)
}

// TODO: datastream_data

func (self DataStream) Clear() { C.datastream_clear(self) }

func (self DataStream) Size() int { return int(C.datastream_size(self)) }

func (self DataStream) PutU8(v uint8) { C.datastream_put_u8(self, C.uint8_t(v)) }
func (self DataStream) PutU16(v uint16) { C.datastream_put_u16(self, C.uint16_t(v)) }
func (self DataStream) PutU32(v uint32) { C.datastream_put_u32(self, C.uint32_t(v)) }
func (self DataStream) PutU64(v uint32) { C.datastream_put_u64(self, C.uint64_t(v)) }

func (self DataStream) PutI8(v int8) { C.datastream_put_i8(self, C.int8_t(v)) }
func (self DataStream) PutI16(v int16) { C.datastream_put_i16(self, C.int16_t(v)) }
func (self DataStream) PutI32(v int32) { C.datastream_put_i32(self, C.int32_t(v)) }
func (self DataStream) PutI64(v int32) { C.datastream_put_i64(self, C.int64_t(v)) }

func (self DataStream) PutData(bytes []byte) {
    size := len(bytes)
    if size > 0 {
        base := (*C.uint8_t)(&bytes[0])
        C.datastream_put_data(self, base, C.size_t(size))
    }
}

func (self DataStream) GetU8() uint8 { return uint8(C.datastream_get_u8(self)) }
func (self DataStream) GetU16() uint16 { return uint16(C.datastream_get_u16(self)) }
func (self DataStream) GetU32() uint32 { return uint32(C.datastream_get_u32(self)) }
func (self DataStream) GetU64() uint64 { return uint64(C.datastream_get_u64(self)) }

func (self DataStream) GetI8() int8 { return int8(C.datastream_get_i8(self)) }
func (self DataStream) GetI16() int16 { return int16(C.datastream_get_i16(self)) }
func (self DataStream) GetI32() int32 { return int32(C.datastream_get_i32(self)) }
func (self DataStream) GetI64() int64 { return int64(C.datastream_get_i64(self)) }


func (self DataStream) GetDataInPlace(length int) []byte {
    base := C.datastream_get_data_inplace(self, C.size_t(length))
    return C.GoBytes(rawptr_t(base), C.int(length))
}

type UInt256 = *C.uint256_t

func NewUInt256() UInt256 { return C.uint256_new() }
func (self UInt256) Free() { C.uint256_free(self) }
func (self UInt256) UInt256IsNull() bool { return bool(C.uint256_is_null(self)) }
func (self UInt256) UInt256IsEq(other UInt256) bool { return bool(C.uint256_is_eq(self, other)) }
func (self UInt256) Serialize(s DataStream) { C.uint256_serialize(self, s) }
func (self UInt256) Unserialize(s DataStream) { C.uint256_unserialize(self, s) }
func (self UInt256) IsEq(other UInt256) bool { return bool(C.uint256_is_eq(self, other)) }

func (self DataStream) GetHash() UInt256 {
    return C.datastream_get_hash(self)
}

func (self DataStream) GetHex() string {
    tmp := C.datastream_get_hex(self)
    res := C.GoString(tmp)
    C.free(rawptr_t(tmp))
    return res
}

func (self UInt256) GetHex() string {
    s := NewDataStream()
    self.Serialize(s)
    res := s.GetHex()
    s.Free()
    return res
}
