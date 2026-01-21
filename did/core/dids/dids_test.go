package dids

import (
	"byd50-ssi/did/configs"
	"encoding/json"
	"strings"
	"testing"
)

func TestCreateDIDRules(t *testing.T) {
	origRule := configs.UseConfig.GenerationRule
	defer func() { configs.UseConfig.GenerationRule = origRule }()

	for _, rule := range []string{"hexdigit", "uuid", "base58"} {
		configs.UseConfig.GenerationRule = rule
		did, doc := CreateDID("byd50", "publicKeyBase58")
		if did == "" || len(doc) == 0 {
			t.Fatalf("did/doc empty for rule: %s", rule)
		}
		if !strings.HasPrefix(did, "did:byd50:") {
			t.Fatalf("unexpected did prefix: %s", did)
		}
	}
}

func TestUpdateDocumentError(t *testing.T) {
	if _, err := UpdateDocument("did:byd50:123", []byte("not-json")); err == nil {
		t.Fatal("expected error for invalid document")
	}
}

func TestDocumentSerialization(t *testing.T) {
	did, doc := CreateDID("byd50", "publicKeyBase58")
	if did == "" || len(doc) == 0 {
		t.Fatal("doc creation failed")
	}

	var parsed DocumentInterface
	if err := json.Unmarshal(doc, &parsed); err != nil {
		t.Fatal(err)
	}
	if parsed.ID != did {
		t.Fatalf("unexpected doc id: %s", parsed.ID)
	}

	if _, err := UpdateDocument(did, doc); err != nil {
		t.Fatalf("update document failed: %v", err)
	}
}

func TestResolutionErrorCodeString(t *testing.T) {
	if InvalidDid.String() == "" || NotFound.String() == "" {
		t.Fatal("invalid resolution error string")
	}
}
