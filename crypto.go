package salticidae

// #include <stdlib.h>
// #include "salticidae/crypto.h"
import "C"
import "runtime"

// The C pointer type for a X509 handle.
type CX509 = *C.x509_t
type x509 struct{ inner CX509 }

// The handle for a X509 certificate.
type X509 = *x509

func X509FromC(ptr CX509) X509 {
	return &x509{inner: ptr}
}

func NewX509FromPemFile(fname string, passwd *string, err *Error) X509 {
	fname_c_str := C.CString(fname)
	passwd_c_str := (*C.char)(nil)
	if passwd != nil {
		passwd_c_str = C.CString(*passwd)
	}
	res := X509FromC(C.x509_new_from_pem_file(fname_c_str, passwd_c_str, err))
	if res != nil {
		runtime.SetFinalizer(res, func(self X509) { self.free() })
	}
	C.free(RawPtr(fname_c_str))
	if passwd_c_str != nil {
		C.free(RawPtr(passwd_c_str))
	}
	return res
}

func NewX509FromDer(der ByteArray, err *Error) X509 {
	res := X509FromC(C.x509_new_from_der(der.inner, err))
	if res != nil {
		runtime.SetFinalizer(res, func(self X509) { self.free() })
	}
	return res
}

func (self X509) free() { C.x509_free(self.inner) }
func (self X509) GetPubKey() PKey {
	res := &pKey{inner: C.x509_get_pubkey(self.inner)}
	runtime.SetFinalizer(res, func(self PKey) { self.free() })
	return res
}

func (self X509) GetDer(autoFree bool) ByteArray {
	res := ByteArrayFromC(C.x509_get_der(self.inner))
	byteArraySetFinalizer(res, autoFree)
	return res
}

// The C pointer type for a PKey handle.
type CPKey = *C.pkey_t
type pKey struct{ inner CPKey }

// The handle for an OpenSSL EVP_PKEY.
type PKey = *pKey

func NewPrivKeyFromPemFile(fname string, passwd *string, err *Error) PKey {
	fname_c_str := C.CString(fname)
	passwd_c_str := (*C.char)(nil)
	if passwd != nil {
		passwd_c_str = C.CString(*passwd)
	}
	res := &pKey{inner: C.pkey_new_privkey_from_pem_file(fname_c_str, passwd_c_str, err)}
	if res != nil {
		runtime.SetFinalizer(res, func(self PKey) { self.free() })
	}
	C.free(RawPtr(fname_c_str))
	if passwd_c_str != nil {
		C.free(RawPtr(passwd_c_str))
	}
	return res
}

func NewPrivKeyFromDer(der ByteArray, err *Error) PKey {
	res := &pKey{inner: C.pkey_new_privkey_from_der(der.inner, err)}
	if res != nil {
		runtime.SetFinalizer(res, func(self PKey) { self.free() })
	}
	return res
}

func (self PKey) free() { C.pkey_free(self.inner) }
func (self PKey) GetPubKeyDer(autoFree bool) ByteArray {
	res := ByteArrayFromC(C.pkey_get_pubkey_der(self.inner))
	byteArraySetFinalizer(res, autoFree)
	return res
}

func (self PKey) GetPrivKeyDer(autoFree bool) ByteArray {
	res := ByteArrayFromC(C.pkey_get_privkey_der(self.inner))
	byteArraySetFinalizer(res, autoFree)
	return res
}
