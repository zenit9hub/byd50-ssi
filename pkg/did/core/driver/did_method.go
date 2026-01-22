package driver

import (
	"sync"
)

var didMethods = map[string]func() DidMethod{}
var didMethodLock = new(sync.RWMutex)

// DidMethod Implement to add new methods for register or resolve 'did'.
type DidMethod interface {
	CreateDid(pbKeyBase58 string) (string, error)          // Returns created did or error
	ResolveDid(did string) (string, string, string, error) // Returns did document or error
	Method() string                                        // returns the method identifier for this method (example: 'byd50')
}

// RegisterDidMethod Register the "method" name and a factory function for signing method.
// This is typically done during init() in the method's implementation
func RegisterDidMethod(method string, f func() DidMethod) {
	didMethodLock.Lock()
	defer didMethodLock.Unlock()

	didMethods[method] = f
}

// GetDidMethod Get DidMethod as "method" string
func GetDidMethod(method string) (methodFunc DidMethod) {
	didMethodLock.RLock()
	defer didMethodLock.RUnlock()

	if methodF, ok := didMethods[method]; ok {
		methodFunc = methodF()
	}
	return
}
