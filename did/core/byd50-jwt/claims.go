package byd50_jwt

import (
	"encoding/json"
	"errors"
	"fmt"
)

// getInt64 Get the integer claim.
func (m MapClaims) getInt64(claim string) (int64, error) {
	var err error
	var int64Claim int64

	v, ok := m[claim]
	if !ok {
		err = errors.New(fmt.Sprintf("couldn't find claim(%v) field", claim))
	}
	switch vType := v.(type) {
	case float64:
		int64Claim = int64(vType)
	case json.Number:
		int64Claim, _ = vType.Int64()
	default:
		err = errors.New("'exp' type error")
	}
	return int64Claim, err
}

// getString Get the 'String' claim.
func (m MapClaims) getString(claim string) (string, error) {
	var err error
	var vString string

	v, ok := m[claim]
	if !ok {
		err = errors.New(fmt.Sprintf("couldn't find claim(%v) field", claim))
	}
	switch vType := v.(type) {
	case string:
		vString = vType
	case []string:
		if len(vType) > 1 {
			err = errors.New("claim has multiple string array")
			return vString, err
		}

		vString = vType[0]
	case []interface{}:
		if len(vType) > 1 {
			err = errors.New("claim has multiple interface{} array")
			return vString, err
		}
		vs, ok := vType[0].(string)
		if !ok {
			err = errors.New("type error inside of claim")
			return vString, err
		}
		vString = vs
	default:
		err = errors.New("'claim' type error")
	}
	return vString, err
}

// getStringArray Get the 'String' claim.
func (m MapClaims) getStringArray(claim string) ([]string, error) {
	var err error
	var stringArray []string

	v, ok := m[claim]
	if !ok {
		err = errors.New(fmt.Sprintf("couldn't find claim(%v) field", claim))
	}
	switch vType := v.(type) {
	case string:
		stringArray = append(stringArray, vType)
	case []string:
		stringArray = vType
	case []interface{}:
		for _, a := range vType {
			vs, ok := a.(string)
			if !ok {
				err = errors.New("type error inside of claim")
				return stringArray, err
			}
			stringArray = append(stringArray, vs)
		}
	default:
		err = errors.New("'claim' type error")
	}
	return stringArray, err
}
