package kafka_message

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
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

func TestConsumeReader(t *testing.T) {

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	address1 := os.Getenv("address1")
	address2 := os.Getenv("address2")

	type args struct {
		brokers   []string
		topic     string
		partition int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "test loop receive message", args: args{
			brokers:   []string{address2, address1},
			topic:     "test-topic",
			partition: 0,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stop := make(chan bool, 100)
			done := make(chan int, 100)

			go func() {
				t.Log("🚩star running")
				if err := ConsumeReader(tt.args.brokers, tt.args.topic, tt.args.partition, stop, done); (err != nil) != tt.wantErr {
					t.Errorf("ConsumeReader() error = %v, wantErr %v", err, tt.wantErr)
				}

			}()

			for {
				select {

				case <-done:
					t.Logf("🎉done!!!")
					close(done)
					close(stop)
					return
				case <-time.After(30 * time.Second):
					stop <- true
					t.Logf("🚗sending stop signal...")
				}
			}
		})
	}
}
