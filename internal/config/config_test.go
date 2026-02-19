package config

import (
	"reflect"
	"testing"
)

func Test_loadConfigFromYamlString(t *testing.T) {
	type args struct {
		yamlString string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "Load config with empty string",
			args: args{
				yamlString: "",
			},
			want:    defaultConfig,
			wantErr: false,
		},
		{
			name: "Load config with invalid yaml",
			args: args{
				yamlString: `
                server_config:
                    host: 0.0.0.0
                    port: 2323
                mudlib_config:
                    mudlib_path: mudlib/
                invalid_yaml`,
			},
			want:    defaultConfig,
			wantErr: true,
		}, {
			name: "Load config with valid yaml",
			args: args{
				yamlString: `
            server_config:
                host: 0.0.0.1
                port: 2323
            mudlib_config:
                mudlib_path: mudlib/`,
			},
			want: Config{
				ServerConfig: ServerConfig{
					Host: "0.0.0.1",
					Port: 2323,
				},
				MudlibConfig: MudlibConfig{
					MudlibPath: "mudlib/",
				},
			},
			wantErr: false,
		}, {
			name: "Load partial config and use default values",
			args: args{
				yamlString: `
            server_config:
                host: 127.0.0.1`,
			},
			want: Config{
				ServerConfig: ServerConfig{
					Host: "127.0.0.1",
					Port: 2323,
				},
				MudlibConfig: MudlibConfig{
					MudlibPath: "mudlib/",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("loadConfigFromYamlString() recovered from panic, wantErr %v", tt.wantErr)
				}
			}()
			if got := loadConfigFromYamlString(tt.args.yamlString); !reflect.DeepEqual(got, &tt.want) {
				t.Errorf("loadConfigFromYamlString() = %v, want %v", got, tt.want)
			}
		})
	}
}
