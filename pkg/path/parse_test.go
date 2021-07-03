package path

import "testing"

func TestGetNamespaced(t *testing.T) {
	type args struct {
		targetPath string
	}
	tests := []struct {
		name     string
		args     args
		wantNs   string
		wantName string
	}{{
		name: "normal case",
		args: args{
			targetPath: "/namespace/default2cqs5/pipeline/go",
		},
		wantNs: "default2cqs5",
		wantName: "go",
	}, {
		name: "not a pipeline router path",
		args: args{
			targetPath: "/namespace/default2cqs5/good/go",
		},
		wantNs: "",
		wantName: "",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNs, gotName := GetPipelineNamespaced(tt.args.targetPath)
			if gotNs != tt.wantNs {
				t.Errorf("GetPipelineNamespaced() gotNs = %v, want %v", gotNs, tt.wantNs)
			}
			if gotName != tt.wantName {
				t.Errorf("GetPipelineNamespaced() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
