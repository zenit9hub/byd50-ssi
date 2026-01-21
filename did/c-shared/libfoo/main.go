package main

import "C"
import (
	foo "byd50-ssi/did/c-shared/foo"
)

//export createKeyPair
func createKeyPair() *C.char {
	foo.CreateKeyPairForAndr()
	return C.CString("success")
}

//export getPbKey
func getPbKey() *C.char {
	return C.CString(foo.GetPublicKeyBase58())
}

//export getPvKey
func getPvKey() *C.char {
	return C.CString(foo.GetPrivateKeyBase58())
}

//export createVp
func createVp(str1 *C.char, str2 *C.char, str3 *C.char, str4 *C.char, str5 *C.char) *C.char {
	return C.CString(foo.CreateVpForAndr(C.GoString(str1), C.GoString(str2), C.GoString(str3), C.GoString(str4), C.GoString(str5)))
}

//export claimsGetExp
func claimsGetExp(vpJwt *C.char) C.long {
	return C.long(foo.ClaimsGetExp(C.GoString(vpJwt)))
}

//export claimsGetIat
func claimsGetIat(vpJwt *C.char) C.long {
	return C.long(foo.ClaimsGetIat(C.GoString(vpJwt)))
}

//export claimsGetInt64
func claimsGetInt64(vpJwt *C.char, claim *C.char) C.long {
	return C.long(foo.ClaimsGetInt64(C.GoString(vpJwt), C.GoString(claim)))
}

func main() {}
