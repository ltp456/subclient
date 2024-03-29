package subclient

/*
#cgo LDFLAGS: -L./lib -lsubapi
#include <stdlib.h>
#include "./lib/subapi.h"
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
	err := c.freeRes(o)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c *Client) ffiSignedExtrinsic(hash, seed, address, amount, nonce, specVersion, transactionVersion, networkId, moduleIndex, callIndex string) (string, error) {
	cHash := C.CString(hash)
	cSeed := C.CString(seed)
	cAddress := C.CString(address)
	cAmount := C.CString(amount)
	cNonce := C.CString(nonce)
	cSpecVersion := C.CString(specVersion)
	cTransactionVersion := C.CString(transactionVersion)
	cNetworkId := C.CString(networkId)
	cModuleIndex := C.CString(moduleIndex)
	cCallIndex := C.CString(callIndex)
	defer C.free(unsafe.Pointer(cHash))
	defer C.free(unsafe.Pointer(cSeed))
	defer C.free(unsafe.Pointer(cAddress))
	defer C.free(unsafe.Pointer(cAmount))
	defer C.free(unsafe.Pointer(cNonce))
	defer C.free(unsafe.Pointer(cSpecVersion))
	defer C.free(unsafe.Pointer(cTransactionVersion))
	defer C.free(unsafe.Pointer(cNetworkId))
	defer C.free(unsafe.Pointer(cModuleIndex))
	defer C.free(unsafe.Pointer(cCallIndex))
	o := C.signed_extrinsic(cHash, cSeed, cAddress, cAmount, cNonce, cSpecVersion, cTransactionVersion, cNetworkId, cModuleIndex, cCallIndex)
	output := C.GoString(o)
	err := c.freeRes(o)
	if err != nil {
		return "", err
	}
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
	err := c.freeRes(o)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c *Client) ffiDecodeExtrinsic(raw string) (string, error) {
	cRaw := C.CString(raw)
	defer C.free(unsafe.Pointer(cRaw))
	o := C.decode_extrinsic(cRaw)
	output := C.GoString(o)
	err := c.freeRes(o)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (c *Client) freeRes(value *C.char) error {
	C.free_res(value)
	return nil
}
