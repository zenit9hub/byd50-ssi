package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetConfigFromFile(t *testing.T) {
	cfg := GetConfig()

	if cfg.DidRegistryPort != ":50051" {
		t.Fatalf("unexpected registry port: %v", cfg.DidRegistryPort)
	}
	if cfg.DidRegistrarPort != ":50052" {
		t.Fatalf("unexpected registrar port: %v", cfg.DidRegistrarPort)
	}
	if !strings.Contains(cfg.EthClientUrl, "http") {
		t.Fatalf("unexpected eth client url: %v", cfg.EthClientUrl)
	}
}

func TestRootDir(t *testing.T) {
	root := RootDir()
	if root == "" {
		t.Fatal("root dir is empty")
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("root dir not found: %v", err)
	}
}

func TestGetConfigModes(t *testing.T) {
	root := RootDir()
	if root == "" {
		t.Fatal("root dir is empty")
	}
	if _, err := os.Getwd(); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, "configs", "configs.yml")

	writeConfig := func(mode string) {
		content := fmt.Sprintf(`
system_mode: "%s"
generation_rule: "hexdigit"
rel_service:
  did-registry:
    address: "rel:50051"
    port: ":50051"
  did-sep:
    address: "rel:50052"
    port: ":50052"
    adopted_driver_list: ["byd50"]
  service_endpoint:
    address: "rel:50053"
    port: ":50053"
  relying_party:
    address: "rel:50054"
    port: ":50054"
  issuer:
    address: "rel:50055"
    port: ":50055"
  eth_client:
    raw_url: "http://rel"
    sc_address: "0x1"
dev_service:
  did-registry:
    address: "dev:50051"
    port: ":50051"
  did-sep:
    address: "dev:50052"
    port: ":50052"
    adopted_driver_list: ["byd50"]
  service_endpoint:
    address: "dev:50053"
    port: ":50053"
  relying_party:
    address: "dev:50054"
    port: ":50054"
  issuer:
    address: "dev:50055"
    port: ":50055"
  eth_client:
    raw_url: "http://dev"
    sc_address: "0x2"
local_service:
  did-registry:
    address: "local:50051"
    port: ":50051"
  did-sep:
    address: "local:50052"
    port: ":50052"
    adopted_driver_list: ["byd50"]
  service_endpoint:
    address: "local:50053"
    port: ":50053"
  relying_party:
    address: "local:50054"
    port: ":50054"
  issuer:
    address: "local:50055"
    port: ":50055"
  eth_client:
    raw_url: "http://local"
    sc_address: "0x3"
`, mode)
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	defer func() {
		_ = os.Remove(configPath)
	}()

	writeConfig("REL")
	cfg := GetConfig()
	if cfg.DidRegistryAddress != "rel:50051" {
		t.Fatalf("expected rel config, got %v", cfg.DidRegistryAddress)
	}

	writeConfig("DEV")
	cfg = GetConfig()
	if cfg.DidRegistryAddress != "dev:50051" {
		t.Fatalf("expected dev config, got %v", cfg.DidRegistryAddress)
	}

	writeConfig("LOCAL")
	cfg = GetConfig()
	if cfg.DidRegistryAddress != "local:50051" {
		t.Fatalf("expected local config, got %v", cfg.DidRegistryAddress)
	}
}
