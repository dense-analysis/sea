package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEscapeSpace(t *testing.T) {
	got := escapeSpace("hello world")
	if got != "hello\\ world" {
		t.Errorf("expected 'hello\\\\ world', got %q", got)
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.toml")
	data := []byte(`listen = 8080
listen_ssl = 443
server_name = "example.com"

redirect_http = true

ssl_certificate = "/tmp/full.pem"
ssl_certificate_key = "/tmp/key.pem"
letsencrypt = true

[[custom_keywords]]
phrase = "foo bar"
dest = "google"`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	cfg, err := loadConfig(path)
	if err != nil {
		t.Fatalf("loadConfig returned error: %v", err)
	}
	if cfg.Listen != 8080 {
		t.Errorf("expected listen 8080, got %d", cfg.Listen)
	}
	if cfg.ListenSSL != 443 {
		t.Errorf("expected listen_ssl 443, got %d", cfg.ListenSSL)
	}
	if !cfg.RedirectHTTP {
		t.Error("expected redirect_http true")
	}
	if cfg.SSLCertificate != "/tmp/full.pem" || cfg.SSLCertificateKey != "/tmp/key.pem" {
		t.Errorf("unexpected certificate paths %+v", cfg)
	}
	if !cfg.LetsEncrypt {
		t.Error("expected letsencrypt true")
	}
	if cfg.ServerName != "example.com" {
		t.Errorf("expected server name example.com, got %s", cfg.ServerName)
	}
	if len(cfg.CustomKeywords) != 1 {
		t.Fatalf("expected 1 custom keyword, got %d", len(cfg.CustomKeywords))
	}
	if cfg.CustomKeywords[0].Phrase != "foo bar" || cfg.CustomKeywords[0].Dest != "google" {
		t.Errorf("unexpected custom keyword %+v", cfg.CustomKeywords[0])
	}
}

func TestGenerateNginx(t *testing.T) {
	cfg := Config{
		Listen:       8080,
		ListenSSL:    8443,
		ServerName:   "example.com",
		RedirectHTTP: true,
		CustomKeywords: []KeywordRule{{
			Phrase: "foo",
			Dest:   "google",
		}},
	}
	out, err := generateNginx(cfg)
	if err != nil {
		t.Fatalf("generateNginx returned error: %v", err)
	}
	if !strings.Contains(out, "server_name example.com;") {
		t.Errorf("generated config missing server name: %s", out)
	}
	if !strings.Contains(out, "listen 8443 ssl") {
		t.Errorf("generated config missing ssl server: %s", out)
	}
	if !strings.Contains(out, "~*(?i)^foo$") {
		t.Errorf("generated config missing custom rule: %s", out)
	}
}
