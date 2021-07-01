package server

import "testing"

func TestRandTokener_Generate(t *testing.T) {
	tests := []struct {
		name string
		r    RandTokener
		want int
	}{
		{
			"generate",
			RandTokener{},
			16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RandTokener{}
			if got := r.Generate(); len(got) != tt.want {
				t.Errorf("RandTokener.Generate() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
