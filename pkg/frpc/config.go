package frpc

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/fatedier/frp/pkg/config"
	"gopkg.in/ini.v1"
)

var configTemplate = `[common]
server_addr = {{ .ServerAddr }}
server_port = {{ .ServerPort }}
authentication_method = token
token = "{{ .Token }}"
authenticate_new_work_conns = true
authenticate_heartbeats = true
admin_addr = {{ .AdminAddr }}
admin_port = {{ .AdminPort }}
includes = ./conf.d/*.ini
log_level = {{ .LogLevel }}
tls_enable = true
tls_cert_file = ./tls/frps/client.crt
tls_key_file = ./tls/frps/client.key
tls_trusted_ca_file = ./tls/frps/ca.crt
meta_client_id = "{{ index .Metas "client_id" }}"
meta_client_signature = "{{ index .Metas "client_signature" }}"
`

// ParseFrpcCommonConfig parses the common section of a frpc configuration file.
// The source can be a string, []byte, or io.Reader.
func ParseFrpcCommonConfig(source interface{}) (config.ClientCommonConf, error) {
	f, err := ini.LoadSources(ini.LoadOptions{
		Insensitive:         false,
		InsensitiveSections: false,
		InsensitiveKeys:     false,
		IgnoreInlineComment: true,
		AllowBooleanKeys:    true,
	}, source)
	if err != nil {
		return config.ClientCommonConf{}, fmt.Errorf("failed to parse configuration file: %s", err)
	}

	s, err := f.GetSection("common")
	if err != nil {
		return config.ClientCommonConf{}, fmt.Errorf("invalid configuration file, not found [common] section")
	}

	common := config.ClientCommonConf{}
	err = s.MapTo(&common)
	if err != nil {
		return config.ClientCommonConf{}, fmt.Errorf("failed to map common config: %s", err)
	}

	return common, nil
}

func SaveTLSFiles(caCert, cert, key string, path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create tls directory: %s", err)
	}

	decodedCaCert, err := base64.StdEncoding.DecodeString(caCert)
	if err != nil {
		return fmt.Errorf("failed to decode ca.crt: %s", err)
	}

	if err := os.WriteFile(filepath.Join(path, "ca.crt"), decodedCaCert, 0600); err != nil {
		return fmt.Errorf("failed to write ca.crt: %s", err)
	}

	decodedCert, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return fmt.Errorf("failed to decode client.crt: %s", err)
	}

	if err := os.WriteFile(filepath.Join(path, "client.crt"), decodedCert, 0600); err != nil {
		return fmt.Errorf("failed to write client.crt: %s", err)
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return fmt.Errorf("failed to decode client.key: %s", err)
	}

	if err := os.WriteFile(filepath.Join(path, "client.key"), decodedKey, 0600); err != nil {
		return fmt.Errorf("failed to write client.key: %s", err)
	}

	return nil
}

func SaveFrpcCommonConfig(common config.ClientCommonConf, path string) error {
	// parse template
	tmpl, err := template.New("frpc.ini").Parse(configTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %s", err)
	}

	// render template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, common)
	if err != nil {
		return fmt.Errorf("failed to render template: %s", err)
	}

	// save to file
	err = os.WriteFile(path, buf.Bytes(), 0600)
	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}
