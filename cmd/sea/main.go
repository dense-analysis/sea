package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/pelletier/go-toml/v2"
)

// Config represents the TOML configuration structure.
type Config struct {
	Listen         int           `toml:"listen"`
	ServerName     string        `toml:"server_name"`
	CustomKeywords []KeywordRule `toml:"custom_keywords"`
}

// KeywordRule maps a phrase to a destination.
type KeywordRule struct {
	Phrase string `toml:"phrase"`
	Dest   string `toml:"dest"`
}

//go:embed nginx.conf.tmpl
var templateFS embed.FS

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "config.toml", "path to TOML configuration")
	flag.Parse()

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	out, err := generateNginx(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating nginx config: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(out)
}

func loadConfig(path string) (Config, error) {
	cfg := Config{Listen: 80, ServerName: "search.local"}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// generateNginx assembles the nginx configuration using heuristics
// and any custom keyword rules from the configuration file.
func generateNginx(cfg Config) (string, error) {
	tmpl, err := template.New("nginx.conf.tmpl").Funcs(template.FuncMap{
		"escape": escapeSpace,
	}).ParseFS(templateFS, "nginx.conf.tmpl")
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, cfg); err != nil {
		return "", err
	}
	return b.String(), nil
}

func escapeSpace(s string) string {
	return strings.ReplaceAll(s, " ", "\\ ")
}
