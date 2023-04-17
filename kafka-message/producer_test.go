package kafka_message

import (
	"github.com/segmentio/kafka-go"
	"os"
	"testing"
)

func TestProduce(t *testing.T) {
	type args struct {
		topic   string
		message []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test with empty topic and empty message",
			args: args{
				topic:   "",
				message: []string{},
			},
			wantErr: true,
		},
		{
			name: "test with empty topic",
			args: args{
				topic:   "",
				message: []string{"hello world"},
			},
			wantErr: true,
		},
		{
			name: "test with empty message",
			args: args{
				topic:   "test-topic",
				message: []string{},
			},
			wantErr: false,
		},
		{
			name: "test with single message",
			args: args{
				topic:   "test-topic",
				message: []string{"hello world"},
			},
			wantErr: false,
		},
		{
			name: "test with multiple messages",
			args: args{
				topic:   "test-topic",
				message: []string{"hello", "world", "this", "is", "a", "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Produce(tt.args.topic, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Produce() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProduceWriter(t *testing.T) {

	address1 := os.Getenv("address1")
	address2 := os.Getenv("address2")
	type args struct {
		address  []string
		topic    string
		messages []kafka.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "single address",
			args: args{
				address: []string{address1},
				topic:   "test-topic",
				messages: []kafka.Message{
					{Value: []byte("hello")},
					{Value: []byte("world")},
				},
			},
			wantErr: false,
		},
		{
			name: "test with multiple addresses",
			args: args{
				address: []string{address1, address2},
				topic:   "test-topic",
				messages: []kafka.Message{
					{Key: []byte("Key-A"), Value: []byte("Hello World!")},
					{Key: []byte("Key-B"), Value: []byte("TWO")},
					{Key: []byte("Key-C"), Value: []byte("THREE")},
				},
			},
			wantErr: false,
		},
		{
			name: "test with same KEY",
			args: args{
				address: []string{address1, address2},
				topic:   "test-topic",
				messages: []kafka.Message{
					{Key: []byte("Key-A"), Value: []byte("testA")},
					{Key: []byte("Key-B"), Value: []byte("testB")},
					{Key: []byte("Key-C"), Value: []byte("testC")},
				},
			},
			wantErr: false,
		},
		{
			name: "test with no topic",
			args: args{
				address: []string{address1, address2},
				topic:   "",
				messages: []kafka.Message{
					{Key: []byte("Key-D"), Value: []byte("test1")},
					{Key: []byte("Key-E"), Value: []byte("test2")},
					{Key: []byte("Key-F"), Value: []byte("test3")},
				},
			},
			wantErr: true,
		},
		{
			name: "test with topic don't exist",
			args: args{
				address: []string{address1, address2},
				topic:   "no-exist-topic",
				messages: []kafka.Message{
					{Key: []byte("Key-D"), Value: []byte("test1")},
					{Key: []byte("Key-E"), Value: []byte("test2")},
					{Key: []byte("Key-F"), Value: []byte("test3")},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProduceWriter(tt.args.address, tt.args.topic, tt.args.messages); (err != nil) != tt.wantErr {
				t.Errorf("ProduceWriter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
