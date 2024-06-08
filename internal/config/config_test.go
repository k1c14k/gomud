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
                    address:	0.0.0.0:2323
                mudlib_config:
                    mudlib_path:	mudlib/
                invalid_yaml`,
			},
			want:    defaultConfig,
			wantErr: true,
		}, {
			name: "Load config with valid yaml",
			args: args{
				yamlString: `
            server_config:
                address:	0.0.0.1:2323
            mudlib_config:
                mudlib_path:	mudlib/`,
			},
			want: Config{
				ServerConfig: struct {
					Address string `yaml:"address"`
				}{
					Address: "0.0.0.1:2323",
				},
				MudlibConfig: struct {
					MudlibPath string `yaml:"mudlib_path"`
				}{
					MudlibPath: "mudlib/",
				},
			},
			wantErr: false,
		}, {
			name: "Load partial config and use default values",
			args: args{
				yamlString: `
            server_config:
                address: 127.0.0.1:2323`,
			},
			want: Config{
				ServerConfig: ServerConfig{
					Address: "127.0.0.1:2323",
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
			if got := loadConfigFromYamlString(tt.args.yamlString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfigFromYamlString() = %v, want %v", got, tt.want)
			}
		})
	}
}
