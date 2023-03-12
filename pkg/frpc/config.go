package frpc

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

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

// func SaveFrpcCommonConfig(common config.ClientCommonConf, path string) error {
// 	f := ini.Empty()
// 	s, err := f.NewSection("common")
// 	if err != nil {
// 		return fmt.Errorf("failed to create [common] section: %s", err)
// 	}

// 	err = s.ReflectFrom(&common)
// 	if err != nil {
// 		return fmt.Errorf("failed to reflect common config: %s", err)
// 	}

// 	return f.SaveTo(path)
// }

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
	err = os.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}

	return nil
}
