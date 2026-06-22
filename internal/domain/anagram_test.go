package domain_test

import (
    "testing"

    "github.com/zandgaim/anagram-counter/internal/domain"
)

func TestSignature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"Listen", "listen", "eilnst"},
        {"Silent", "silent", "eilnst"}, // Should match 'listen'
        {"Numbers", "123", "123"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := domain.Signature(tt.input)
            if result != tt.expected {
                t.Errorf("Signature(%q) = %q; want %q", tt.input, result, tt.expected)
            }
        })
    }
}