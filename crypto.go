package salticidae

// #include <stdlib.h>
// #include "salticidae/crypto.h"
import "C"
import "runtime"

//// begin X509 def

// CX509 is the C pointer type for a X509 handle.
type CX509 = *C.x509_t
type x509 struct {
	inner    CX509
	autoFree bool
}

// X509 is the handle for a X509 certificate.
type X509 = *x509

// X509FromC converts an existing C pointer into a go pointer
func X509FromC(ptr CX509) X509 {
	return &x509{inner: ptr}
}

func x509SetFinalizer(res X509, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self X509) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (x509 X509) Free() {
	C.x509_free(x509.inner)
	if x509.autoFree {
		runtime.SetFinalizer(x509, nil)
	}
}

//// end X509 def

//// begin PKey def

// CPKey is the C pointer type for a PKey handle.
type CPKey = *C.pkey_t
type pKey struct {
	inner    CPKey
	autoFree bool
}

// PKey is the handle for an OpenSSL EVP_PKEY.
type PKey = *pKey

// PKeyFromC converts an existing C pointer into a go pointer
func PKeyFromC(ptr CPKey) PKey {
	return &pKey{inner: ptr}
}

func pKeySetFinalizer(res PKey, autoFree bool) {
	res.autoFree = autoFree
	if res != nil && autoFree {
		runtime.SetFinalizer(res, func(self PKey) { self.Free() })
	}
}

// Free manually frees the underlying C pointer.
func (pkey PKey) Free() {
	C.pkey_free(pkey.inner)
	if pkey.autoFree {
		runtime.SetFinalizer(pkey, nil)
	}
}

//// end PKey def

//// begin X509 methods

// NewX509FromPemFile loads X509 object from a file.
func NewX509FromPemFile(fname string, passwd *string, err *Error) (res X509) {
	fnameCStr := C.CString(fname)
	passwdCStr := (*C.char)(nil)
	if passwd != nil {
		passwdCStr = C.CString(*passwd)
	}
	res = X509FromC(C.x509_new_from_pem_file(fnameCStr, passwdCStr, err))
	x509SetFinalizer(res, true)
	C.free(rawPtr(fnameCStr))
	if passwdCStr != nil {
		C.free(rawPtr(passwdCStr))
	}
	return
}

// NewX509FromDer loads X509 object from a DER ByteArray.
func NewX509FromDer(der ByteArray, err *Error) (res X509) {
	res = X509FromC(C.x509_new_from_der(der.inner, err))
	x509SetFinalizer(res, true)
	runtime.KeepAlive(der)
	return
}

// GetPubKey returns the public key.
func (x509 X509) GetPubKey() (res PKey) {
	res = PKeyFromC(C.x509_get_pubkey(x509.inner))
	pKeySetFinalizer(res, true)
	runtime.KeepAlive(x509)
	return
}

// GetDer returns the DER format copy.
func (x509 X509) GetDer(autoFree bool) (res ByteArray) {
	res = ByteArrayFromC(C.x509_get_der(x509.inner))
	byteArraySetFinalizer(res, autoFree)
	runtime.KeepAlive(x509)
	return
}

//// end X509 methods

//// begin PKey methods

func NewPrivKeyFromPemFile(fname string, passwd *string, err *Error) (res PKey) {
	fnameCStr := C.CString(fname)
	passwdCStr := (*C.char)(nil)
	if passwd != nil {
		passwdCStr = C.CString(*passwd)
	}
	res = PKeyFromC(C.pkey_new_privkey_from_pem_file(fnameCStr, passwdCStr, err))
	pKeySetFinalizer(res, true)
	C.free(rawPtr(fnameCStr))
	if passwdCStr != nil {
		C.free(rawPtr(passwdCStr))
	}
	return
}

func NewPrivKeyFromDer(der ByteArray, err *Error) (res PKey) {
	res = PKeyFromC(C.pkey_new_privkey_from_der(der.inner, err))
	pKeySetFinalizer(res, true)
	runtime.KeepAlive(der)
	return
}

func (pkey PKey) GetPubKeyDer(autoFree bool) (res ByteArray) {
	res = ByteArrayFromC(C.pkey_get_pubkey_der(pkey.inner))
	byteArraySetFinalizer(res, autoFree)
	runtime.KeepAlive(pkey)
	return
}

func (pkey PKey) GetPrivKeyDer(autoFree bool) (res ByteArray) {
	res = ByteArrayFromC(C.pkey_get_privkey_der(pkey.inner))
	byteArraySetFinalizer(res, autoFree)
	runtime.KeepAlive(pkey)
	return
}

//// end PKey methods
