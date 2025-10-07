package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTmpl(t *testing.T) {
	tests := []struct {
		name     string
		template string
		ctx      any
		want     string
	}{
		{"basic", "{{.USER}}", map[string]string{"USER": "test"}, "test"},
		{"upper", "{{.USER | upper}}", map[string]string{"USER": "test"}, "TEST"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := tmpl(strings.NewReader(tt.template), false, &buf, tt.ctx)
			if err != nil {
				t.Fatalf("tmpl() error = %v", err)
			}
			if got := buf.String(); !strings.Contains(got, tt.want) {
				t.Errorf("tmpl() = %q, want to contain %q", got, tt.want)
			}
		})
	}
}
