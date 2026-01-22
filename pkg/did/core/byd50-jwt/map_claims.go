package byd50_jwt

import (
	"encoding/json"
	"errors"
	// "fmt"
)

// MapClaims type that uses the map[string]interface{} for JSON decoding
// This is the default claims type if you don't supply one
type MapClaims map[string]interface{}

// GetAudience Get the audience
func (m MapClaims) GetAudience() ([]string, error) {
	var aud []string
	switch v := m["aud"].(type) {
	case string:
		aud = append(aud, v)
	case []string:
		aud = v
	case []interface{}:
		for _, a := range v {
			vs, ok := a.(string)
			if !ok {
				return aud, errors.New("'aud' type error")
			}
			aud = append(aud, vs)
		}
	}
	return aud, nil
}

// GetExpiresAt Get the expiresAt
func (m MapClaims) GetExpiresAt() (int64, error) {
	var err error
	var expiresAt int64

	exp, ok := m["exp"]
	if !ok {
		err = errors.New("claims hasn't exp field")
	}
	switch expType := exp.(type) {
	case float64:
		expiresAt = int64(expType)
	case json.Number:
		expiresAt, _ = expType.Int64()
	default:
		err = errors.New("'exp' type error")
	}
	return expiresAt, err
}

// GetIssuedAt Get the issuedAt.
func (m MapClaims) GetIssuedAt() (int64, error) {
	issuedAt, err := m.getInt64("iat")
	return issuedAt, err
}

// GetIssuer Get the issuer.
func (m MapClaims) GetIssuer() (string, error) {
	issuer, err := m.getString("iss")
	return issuer, err
}

// GetNotBefore Get the notBefore.
func (m MapClaims) GetNotBefore() (int64, error) {
	notBefore, err := m.getInt64("nbf")
	return notBefore, err
}

// GetVc Get the VC.
func (m MapClaims) GetVc() (map[string]interface{}, error) {
	var err error

	vc, ok := m["vc"].(map[string]interface{})
	if !ok {
		err = errors.New("couldn't find vp field")
	}

	return vc, err
}

// GetVcType Get the VC type.
func (m MapClaims) GetVcType() ([]string, error) {
	var err error

	vc, ok := m["vc"].(map[string]interface{})
	if !ok {
		err = errors.New("couldn't find vc field")
	}

	var vTyp []string

	v, ok := vc["type"]
	if !ok {
		err = errors.New("couldn't find 'type' field")
	}

	switch vType := v.(type) {
	case string:
		vTyp = append(vTyp, vType)
	case []string:
		vTyp = vType
	case []interface{}:
		for _, a := range vType {
			vs, ok := a.(string)
			if !ok {
				err = errors.New("type error inside of claim")
				return vTyp, err
			}
			vTyp = append(vTyp, vs)
		}
	default:
		err = errors.New("'claim' type error")
	}

	return vTyp, err
}

// GetVp Get the VP.
func (m MapClaims) GetVp() (map[string]interface{}, error) {
	var err error

	vp, ok := m["vp"].(map[string]interface{})
	if !ok {
		err = errors.New("couldn't find vp field")
	}

	return vp, err
}

// GetVpType Get the VP type.
func (m MapClaims) GetVpType() ([]string, error) {
	var err error

	vp, ok := m["vp"].(map[string]interface{})
	if !ok {
		err = errors.New("couldn't find vp field")
	}

	var vTyp []string

	v, ok := vp["type"]
	if !ok {
		err = errors.New("couldn't find 'type' field")
	}

	switch vType := v.(type) {
	case string:
		vTyp = append(vTyp, vType)
	case []string:
		vTyp = vType
	case []interface{}:
		for _, a := range vType {
			vs, ok := a.(string)
			if !ok {
				err = errors.New("type error inside of claim")
				return vTyp, err
			}
			vTyp = append(vTyp, vs)
		}
	default:
		err = errors.New("'claim' type error")
	}

	return vTyp, err
}
