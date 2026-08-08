package ui

import (
	"reflect"
	"testing"
)

func TestWrapLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		maxWidth int
		want     []string
	}{
		{
			name:     "short line within limit",
			line:     "hello world",
			maxWidth: 20,
			want:     []string{"hello world"},
		},
		{
			name:     "line wrapped at spaces",
			line:     "6 goblins ambush the party!",
			maxWidth: 15,
			want:     []string{"6 goblins", "ambush the", "party!"},
		},
		{
			name:     "unbroken long word broken",
			line:     "abcdefghijklmnopqrstuvwxyz",
			maxWidth: 10,
			want:     []string{"abcdefghij", "klmnopqrst", "uvwxyz"},
		},
		{
			name:     "zero or negative width returns original",
			line:     "hello world",
			maxWidth: 0,
			want:     []string{"hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WrapLine(tt.line, tt.maxWidth)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WrapLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapText(t *testing.T) {
	text := "line one is very long and should wrap\nline two"
	got := WrapText(text, 15)
	want := []string{"line one is", "very long and", "should wrap", "line two"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WrapText() = %v, want %v", got, want)
	}
}
