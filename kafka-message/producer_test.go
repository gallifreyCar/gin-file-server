package kafka_message

import "testing"

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
