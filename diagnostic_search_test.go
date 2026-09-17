package main

import (
	"reflect"
	"testing"
)

func TestSelectDiagnostics(t *testing.T) {
	tests := []struct {
		name       string
		candidates []string
		indexes    []int
		want       []string
	}{
		{name: "release evidence outranks build noise", candidates: []string{"build timeout", "release rollback"}, indexes: []int{1, 0}, want: []string{"release rollback", "build timeout"}},
		{name: "invalid result index is ignored", candidates: []string{"compiler diagnostic"}, indexes: []int{3, 0}, want: []string{"compiler diagnostic"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectDiagnostics(tt.candidates, tt.indexes); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("selected diagnostics = %v, want %v", got, tt.want)
			}
		})
	}
}
