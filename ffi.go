package subclient

/*
#cgo LDFLAGS: -L./lib -lccapsapi
#include <stdlib.h>
#include "./lib/ccapsapi.h"
*/
import "C"
import "unsafe"

func (c *Client) ffiPalletInfo(palletName, callName, metadata string) (string, error) {
	cPalletName := C.CString(palletName)
	cCallName := C.CString(callName)
	cMetadata := C.CString(metadata)
	defer C.free(unsafe.Pointer(cPalletName))
	defer C.free(unsafe.Pointer(cCallName))
	defer C.free(unsafe.Pointer(cMetadata))
	o := C.pallet_info(cMetadata, cPalletName, cCallName)
	output := C.GoString(o)
	return output, nil
}

func (c *Client) ffiSignedExtrinsic(hash, seed, address, amount, nonce, specVersion, transactionVersion string) (string, error) {
	cHash := C.CString(hash)
	cSeed := C.CString(seed)
	cAddress := C.CString(address)
	cAmount := C.CString(amount)
	cNonce := C.CString(nonce)
	cSpecVersion := C.CString(specVersion)
	cTransactionVersion := C.CString(transactionVersion)
	defer C.free(unsafe.Pointer(cHash))
	defer C.free(unsafe.Pointer(cSeed))
	defer C.free(unsafe.Pointer(cAddress))
	defer C.free(unsafe.Pointer(cAmount))
	defer C.free(unsafe.Pointer(cNonce))
	defer C.free(unsafe.Pointer(cSpecVersion))
	defer C.free(unsafe.Pointer(cTransactionVersion))
	o := C.signed_extrinsic(cHash, cSeed, cAddress, cAmount, cNonce, cSpecVersion, cTransactionVersion)
	output := C.GoString(o)
	return output, nil
}

func (c *Client) ffiDynamicDecodeStorage(palletName, storageEntry, raw, metadata string) (string, error) {
	cPalletName := C.CString(palletName)
	cStorageEntry := C.CString(storageEntry)
	cRaw := C.CString(raw)
	cMetadata := C.CString(metadata)

	defer C.free(unsafe.Pointer(cPalletName))
	defer C.free(unsafe.Pointer(cStorageEntry))
	defer C.free(unsafe.Pointer(cRaw))
	defer C.free(unsafe.Pointer(cMetadata))
	o := C.dynamic_decode_storage(cPalletName, cStorageEntry, cRaw, cMetadata)
	output := C.GoString(o)
	return output, nil
}

func (c *Client) ffiDecodeExtrinsic(raw string) (string, error) {
	cRaw := C.CString(raw)
	defer C.free(unsafe.Pointer(cRaw))
	o := C.decode_extrinsic(cRaw)
	output := C.GoString(o)
	return output, nil
}
