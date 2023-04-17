package kafka_message

import (
	"strings"
	"testing"
)

func TestConsume(t *testing.T) {
	tests := []struct {
		name        string
		topic       string
		wantMessage string
		wantErr     bool
	}{
		{
			name:        "Test Consume Returns Correct Message",
			topic:       "test-topic",
			wantMessage: "hello world",
			wantErr:     false,
		},
		{
			name:        "Test Empty Errors",
			topic:       "",
			wantMessage: "",
			wantErr:     true,
		},
		{
			name:        "Test Consume Handles Errors",
			topic:       "nonexistent-topic",
			wantMessage: "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage, err := Consume(tt.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("Consume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(gotMessage, tt.wantMessage) {
				t.Errorf("Consume() gotMessage = %v, want %v", gotMessage, tt.wantMessage)
			}
		})
	}
}
