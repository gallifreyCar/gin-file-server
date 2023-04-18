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
	stop1 := make(chan bool, 100)
	done1 := make(chan int, 100)
	stop2 := make(chan bool, 100)
	done2 := make(chan int, 100)

	type args struct {
		brokers   []string
		topic     string
		partition int
		stop      chan bool
		done      chan int
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
			stop:      stop1,
			done:      done1,
		}},
		{name: "test loop receive message2", args: args{
			brokers:   []string{address2, address1},
			topic:     "no-exist-topic",
			partition: 0,
			stop:      stop2,
			done:      done2,
		}},
		{name: "test upload files", args: args{
			brokers:   []string{address2, address1},
			topic:     "upload-file-multiple",
			partition: 0,
			stop:      make(chan bool, 100),
			done:      make(chan int, 100),
		}},
		{name: "test upload single file", args: args{
			brokers:   []string{address2, address1},
			topic:     "upload-file-single",
			partition: 0,
			stop:      make(chan bool, 100),
			done:      make(chan int, 100),
		}},
		{name: "test download file", args: args{
			brokers:   []string{address2, address1},
			topic:     "download-file",
			partition: 0,
			stop:      make(chan bool, 100),
			done:      make(chan int, 100),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			go func() {
				t.Log("🚩star running")
				if err := ConsumeReader(tt.args.brokers, tt.args.topic, tt.args.partition, tt.args.stop, tt.args.done); (err != nil) != tt.wantErr {
					t.Errorf("ConsumeReader() error = %v, wantErr %v", err, tt.wantErr)
				}

			}()

			for {
				select {

				case <-tt.args.done:
					t.Logf("🎉done!!!")
					close(tt.args.done)
					close(tt.args.stop)
					return
				case <-time.After(30 * time.Second):
					tt.args.stop <- true
					t.Logf("🚗sending stop signal...")
				}
			}
		})
	}
}
