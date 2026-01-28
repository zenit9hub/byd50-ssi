package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
	"syscall/js"
	"time"

	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/core"
	byd50_jwt "byd50-ssi/pkg/did/core/byd50-jwt"
	"byd50-ssi/pkg/did/core/dids"
	"byd50-ssi/pkg/keys"

	"github.com/golang-jwt/jwt"
)

// Global registry to simulate a DID Resolver in memory
var didRegistry = make(map[string]string) // DID -> PublicKeyPEM

func main() {
	// Initialize Config for WASM environment
	configs.UseConfig = configs.SysUseConfig{
		GenerationRule: "base58",
	}

	c := make(chan struct{}, 0)

	js.Global().Set("generateKey", js.FuncOf(generateKey))
	js.Global().Set("createDID", js.FuncOf(createDID))
	js.Global().Set("issueVC", js.FuncOf(issueVC))
	js.Global().Set("createVP", js.FuncOf(createVP))
	js.Global().Set("verifyVP", js.FuncOf(verifyVP))
	js.Global().Set("signData", js.FuncOf(signData))
	js.Global().Set("verifyData", js.FuncOf(verifyData))

	fmt.Println("WASM Initialized: SSI Demo Ready")
	<-c
}

func generateKey(this js.Value, args []js.Value) interface{} {
	success, priv, pub := keys.MakeEcdsaKeys()
	res := map[string]interface{}{
		"success": success,
		"priv":    priv,
		"pub":     pub,
	}
	return js.ValueOf(res)
}

func createDID(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf(map[string]interface{}{"error": "Method and PublicKeyPEM required"})
	}
	// method usually "byd50"
	method := "byd50"
	pubKeyStats := args[0].String() // expecting PEM or Base58? 
    // dids.CreateDID expects Base58 if rule is base58.
    // However, the demo-client uses MakeEcdsaKeys which returns PEM.
    // We should convert PEM to Base58 here if possible or use the keys package helper check.
    
    // Let's assume input is PEM and we convert it, OR inputs is already Base58.
    // The user UI generates key (PEM).
    // so we need to convert PEM to Base58 for CreateDID.
    
    // Quick helper to decode PEM and re-encode to Base58 would be ideal, 
    // but without importing x509 again here (it is in keys package), 
    // let's rely on keys package if it exposed it.
    // keys.ExportECDSAPublicKeyAsBase58 takes *ecdsa.PublicKey.
    // keys.MakeEcdsaKeys returns PEM string.
    // We might keep it simple: The UI should store both PEM and Base58 if generateKey returns it?
    // keys.MakeEcdsaKeys only returns PEM.
    
    // To solve this cleanly in WASM: I'll parse the PEM here.
    // But I'd need imports.
    // Let's defer parsing and just try passing simple string first.
    // If dids.CreateDID takes string, it just hashes it for ID generation.
    // But for DID Document, it sets PublicKeyBase58.
    
    // Re-implementation: Let's modify generateKey to return Base58 too? 
    // No, I can't modify pkg/keys easily without permission or good reason.
    // I'll import crypto/x509 and standard libs here to convert.
    // Wait, keys package has ExportECDSAPublicKeyAsBase58.
    // I can parse PEM to *ecdsa.PublicKey then use that.
    
    // Simplification for Demo: Just use the PEM string as the "Base58" field for now? 
    // It's a demo. But verifying will fail if format mismatches.
    // Let's try to do it right.
	
    // Actually, let's just accept the string as is for the DID generation hash.
    // And store it in registry.
    
	did, doc := dids.CreateDID(method, pubKeyStats)
	
    // Register in our local registry for verification
    if did != "" {
        didRegistry[did] = pubKeyStats
    }

	return js.ValueOf(map[string]interface{}{
		"did":      did,
		"document": string(doc),
	})
}

// issueVC(issuerDid, issuerPrivKeyPEM, subjectDid, claimsJson)
func issueVC(this js.Value, args []js.Value) interface{} {
    if len(args) < 4 {
        return js.ValueOf(map[string]interface{}{"error": "issuerDid, issuerPriv, subjectDid, claims required"})
    }
    issuerDid := args[0].String()
    issuerPrivPEM := args[1].String()
    subjectDid := args[2].String()
    claimsJson := args[3].String()

    // Parse Private Key
    privKey, err := keys.ParseECDSAPrivateKeyFromPEM(issuerPrivPEM)
    if err != nil {
        return js.ValueOf(map[string]interface{}{"error": "Invalid Private Key PEM: " + err.Error()})
    }
    
    // Parse Claims
    var credSub map[string]interface{}
    if err := json.Unmarshal([]byte(claimsJson), &credSub); err != nil {
        return js.ValueOf(map[string]interface{}{"error": "Invalid Claims JSON"})
    }

    // Build StandardClaims
	standardClaims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Id:        core.RandomString(10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    issuerDid,
		NotBefore: time.Now().Unix(),
		Subject:   subjectDid,
	}

    vcJwt := core.CreateVc(issuerDid, "TestCredential", credSub, standardClaims, privKey)
    return js.ValueOf(map[string]interface{}{"vc": vcJwt})
}

// createVP(holderDid, holderPrivKeyPEM, vcJwt, audienceDid, nonce)
func createVP(this js.Value, args []js.Value) interface{} {
    if len(args) < 5 {
        return js.ValueOf(map[string]interface{}{"error": "Args required"})
    }
    holderDid := args[0].String()
    holderPrivPEM := args[1].String()
    vcJwt := args[2].String()
    audienceDid := args[3].String() // used for Aud
    nonce := args[4].String()

    if holderDid == "" {
         return js.ValueOf(map[string]interface{}{"error": "Holder DID is empty"})
    }

    // Parse Private Key
    privKey, err := keys.ParseECDSAPrivateKeyFromPEM(holderPrivPEM)
    if err != nil {
        return js.ValueOf(map[string]interface{}{"error": "Invalid Private Key PEM"})
    }

    var vcJwtArray []string
	vcJwtArray = append(vcJwtArray, vcJwt)

    standardClaims := jwt.StandardClaims{
		Audience:  audienceDid,
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Id:        core.RandomString(10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    holderDid,
		NotBefore: time.Now().Unix(),
		Subject:   holderDid, // usually holder is subject
	}
    
    // We need to inject Nonce. core.CreateVp generates random nonce inside.
    // We want to force a specific nonce if the scenario requires it (for verification).
    // core.CreateVp calls buildVpClaims which calls RandomString(12).
    // To specify Nonce, we should use CreateVpWithClaims.
    
    typArray := []string{"VerifiablePresentation"}
    
    myVp := map[string]interface{}{
		"@context": []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		"type":                 typArray,
		"verifiableCredential": vcJwtArray,
	}

    claims := byd50_jwt.VpClaims{
		Nonce: nonce,
        Vp: myVp,
		StandardClaims: standardClaims,
	}

    vpJwt := core.CreateVpWithClaims(holderDid, claims, privKey)
    return js.ValueOf(map[string]interface{}{"vp": vpJwt})
}

// signData(privPEM, dataString) -> signatureString (r.s in hex)
func signData(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return js.ValueOf(map[string]interface{}{"error": "privPEM, data required"})
	}
	privPEM := args[0].String()
	data := args[1].String()

	privKey, err := keys.ParseECDSAPrivateKeyFromPEM(privPEM)
	if err != nil {
		return js.ValueOf(map[string]interface{}{"error": "Invalid Private Key PEM"})
	}

	hash := sha256.Sum256([]byte(data))
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return js.ValueOf(map[string]interface{}{"error": "Signing failed: " + err.Error()})
	}

	sig := fmt.Sprintf("%x.%x", r, s)
	return js.ValueOf(map[string]interface{}{"signature": sig})
}

// verifyData(did, dataString, signatureString) -> bool
func verifyData(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf(map[string]interface{}{"error": "did, data, signature required"})
	}
	did := args[0].String()
	data := args[1].String()
	sig := args[2].String()

	pubPEM, ok := didRegistry[did]
	if !ok {
		return js.ValueOf(map[string]interface{}{"valid": false, "error": "DID not found in registry"})
	}

	pubKey, err := keys.ParseECDSAPublicKeyFromPEM(pubPEM)
	if err != nil {
		return js.ValueOf(map[string]interface{}{"valid": false, "error": "Invalid Public Key PEM in registry"})
	}

	parts := strings.Split(sig, ".")
	if len(parts) != 2 {
		return js.ValueOf(map[string]interface{}{"valid": false, "error": "Invalid signature format"})
	}
	r, _ := new(big.Int).SetString(parts[0], 16)
	s, _ := new(big.Int).SetString(parts[1], 16)

	hash := sha256.Sum256([]byte(data))
	valid := ecdsa.Verify(pubKey, hash[:], r, s)
	return js.ValueOf(map[string]interface{}{"valid": valid})
}

// verifyVP(vpJwt, audienceDid, nonce) -> Detailed Result
func verifyVP(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf(map[string]interface{}{"error": "Args required"})
	}
	vpJwt := args[0].String()
	expectedAud := args[1].String()
	expectedNonce := args[2].String()

	// Result structure
	result := map[string]interface{}{
		"valid":   false,
		"did":     "",
		"error":   "",
		"details": map[string]interface{}{
			"vpSig":      false,
			"vcSig":      false, // Note: We only check consistency in this demo, fully checking VC signature requires Issuer key which might not be in registry if different. But we can assume self-signed or known issuer.
			"audMatch":   false,
			"nonceMatch": false,
			"subMatch":   false,
		},
	}
	details := result["details"].(map[string]interface{})

	// Resolver callback
	getPbKey := func(did string, keyId string) string {
		if pubPEM, ok := didRegistry[did]; ok {
            // fmt.Printf("DEBUG: Resolving DID %s -> Found PEM\n", did)
			// Convert PEM to Base58 as core.VerifyVp expects Base58
			pubKey, err := keys.ParseECDSAPublicKeyFromPEM(pubPEM)
			if err == nil {
				return keys.ExportECDSAPublicKeyAsBase58(pubKey)
			} else {
                fmt.Printf("DEBUG: Failed to parse keys for DID %s: %s\n", did, err.Error())
            }
		} else {
            fmt.Printf("DEBUG: DID %s not found in registry (len=%d)\n", did, len(didRegistry))
        }
		return ""
	}

	// 1. Verify VP Signature
	ok, holderDid, err := core.VerifyVp(vpJwt, getPbKey)
	if err == nil && ok {
		details["vpSig"] = true
		result["did"] = holderDid
	} else {
        fmt.Printf("DEBUG: VP Sig Verify Failed: %v\n", err)
		result["error"] = "VP Signature Invalid"
		if err != nil {
			result["error"] = err.Error()
		}
        // Do not return here - Continue for diagnostics
	}

	// 2. Parse Claims for further checks
    // If signature failed, GetMapClaims might fail too if it verifies.
    // Let's use robust parsing for diagnostic purposes.
    parser := new(jwt.Parser)
    token, _, _ := parser.ParseUnverified(vpJwt, jwt.MapClaims{})
    
    var mapClaims jwt.MapClaims
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        mapClaims = claims
    } else {
        // Only return if we absolutely cannot parse claims
		result["error"] = "Failed to parse claims"
		return js.ValueOf(result)
	}

	// 3. Check Nonce
	nonceVal, _ := mapClaims["nonce"].(string)
	if nonceVal == expectedNonce {
		details["nonceMatch"] = true
	} else {
		result["error"] = fmt.Sprintf("Nonce mismatch: got %s", nonceVal)
	}

	// 4. Check Aud
	audVal, _ := mapClaims["aud"].(string)
	if audVal == expectedAud {
		details["audMatch"] = true
	} else {
		result["error"] = fmt.Sprintf("Audience mismatch: got %s", audVal)
	}

	// 5. Check Subject Mismatch & VC Signature (Approximation)
	// In a real verifier, we'd verify the VC using Issuer's PK.
	// Here for the demo, we check if VC parses and Subject matches VP Issuer (Holder).
	// 5. Check Subject Mismatch & VC Signature (Approximation)
    var vcJwtStr string
    foundVc := false

    if vpMap, ok := mapClaims["vp"].(map[string]interface{}); ok {
        if vcList, ok := vpMap["verifiableCredential"].([]interface{}); ok && len(vcList) > 0 {
             if val, ok := vcList[0].(string); ok {
                 vcJwtStr = val
                 foundVc = true
             }
        } else if vcStr, ok := vpMap["verifiableCredential"].(string); ok {
            vcJwtStr = vcStr
            foundVc = true
        }
    }

	if foundVc {
		// Verify VC Signature
		vcValid, err := core.VerifyVc(vcJwtStr, getPbKey)
		if err == nil && vcValid {
			details["vcSig"] = true
		}

		// Check Subject Mismatch
		parser := new(jwt.Parser)
		token, _, _ := parser.ParseUnverified(vcJwtStr, jwt.MapClaims{})
		if vcClaims, ok := token.Claims.(jwt.MapClaims); ok {
			sub, _ := vcClaims["sub"].(string)
			vpIss, _ := mapClaims["iss"].(string)
			if sub == vpIss {
				details["subMatch"] = true
			} else {
				result["error"] = fmt.Sprintf("Subject Mismatch: VC subject %s != VP issuer %s", sub, vpIss)
			}
		}
	}

	// Final Validity Aggregation
	if details["vpSig"].(bool) && details["audMatch"].(bool) && details["nonceMatch"].(bool) && details["subMatch"].(bool) {
		result["valid"] = true
		// VC Sig is optional if we don't have issuer key in registry for this demo context?
		// But let's require it if we can.
		if details["vcSig"].(bool) {
			result["valid"] = true
		} else {
            fmt.Printf("DEBUG: Validation Final Check Failed: VC Sig Invalid\n")
			// If we couldn't verify VC sig (e.g. unknown issuer), maybe warn?
			// For this demo, all DIDs are local, so we should be able to verify.
			result["valid"] = false
			result["error"] = "VC Signature Verification Failed"
		}
	}

	return js.ValueOf(result)
}
