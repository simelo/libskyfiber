package main

import (
	cli "github.com/skycoin/skycoin/src/cli"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cli_GetWalletOutputsFromFile
func SKY_cli_GetWalletOutputsFromFile(_c C.WebRpcClient__Handle, _walletFile string, _arg2 *C.ReadableOutputSet__Handle) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_BAD_HANDLE
		return
	}
	walletFile := _walletFile
	__arg2, ____return_err := cli.GetWalletOutputsFromFile(c, walletFile)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerReadableOutputSetHandle(__arg2)
	}
	return
}

//export SKY_cli_GetWalletOutputs
func SKY_cli_GetWalletOutputs(_c C.WebRpcClient__Handle, _wlt *C.Wallet__Handle, _arg2 *C.ReadableOutputSet__Handle) (____error_code uint32) {
	____error_code = SKY_OK
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	checkAPIReady()
	c, okc := lookupWebRpcClientHandle(_c)
	if !okc {
		____error_code = SKY_BAD_HANDLE
		return
	}
	wlt, okwlt := lookupWalletHandle(*_wlt)
	if !okwlt {
		____error_code = SKY_BAD_HANDLE
		return
	}
	__arg2, ____return_err := cli.GetWalletOutputs(c, wlt)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
		*_arg2 = registerReadableOutputSetHandle(__arg2)
	}
	return
}
