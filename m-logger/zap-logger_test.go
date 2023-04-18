package m_logger

import (
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestInitZapLogger(t *testing.T) {
	type args struct {
		fileName string
		prefix   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{name: "Test zap logger", args: args{
			fileName: "test_zap_logger",
			prefix:   "zap-logger_test/TestInitZapLogger",
		}, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogger, gotErr, gotClose := InitZapLogger(tt.args.fileName, tt.args.prefix)

			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("InitZapLogger() gotErr = %v, want %v", gotErr, tt.wantErr)
			}

			gotLogger.Info("Test Info")
			gotLogger.Log(zap.InfoLevel, gotLogger.Level().String())

			gotLogger.Debug("Test Debug")
			gotLogger.Log(zap.DebugLevel, gotLogger.Level().String())

			gotLogger.Debug("Test Debug")
			gotLogger.Log(zap.WarnLevel, gotLogger.Level().String())

			gotLogger.Error("Test Error")
			gotLogger.Log(zap.ErrorLevel, gotLogger.Level().String())

			err := gotLogger.Sync()
			if err != nil {
				t.Errorf(err.Error())
			}
			defer gotClose()
		})
	}
}
