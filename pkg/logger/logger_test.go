package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestParseLevel(t *testing.T) {
	type given struct {
		level Level
	}

	type want struct {
		level zapcore.Level
		err   bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"debug": {
			given: given{
				level: DebugLevel,
			},
			want: want{
				level: zapcore.DebugLevel,
			},
		},
		"info": {
			given: given{
				level: InfoLevel,
			},
			want: want{
				level: zapcore.InfoLevel,
			},
		},
		"warn": {
			given: given{
				level: WarnLevel,
			},
			want: want{
				level: zapcore.WarnLevel,
			},
		},
		"error": {
			given: given{
				level: ErrorLevel,
			},
			want: want{
				level: zapcore.ErrorLevel,
			},
		},
		"invalid": {
			given: given{
				level: "trace",
			},
			want: want{
				err: true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseLevel(tt.given.level)

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.level, got)
		})
	}
}

func TestNew(t *testing.T) {
	type given struct {
		opts []Option
	}

	type want struct {
		err bool
	}

	tests := map[string]struct {
		given given
		want  want
	}{
		"default production": {
			given: given{
				opts: []Option{WithLevel(InfoLevel)},
			},
		},
		"development": {
			given: given{
				opts: []Option{WithLevel(DebugLevel), WithDevelopment(true)},
			},
		},
		"invalid level": {
			given: given{
				opts: []Option{WithLevel("trace")},
			},
			want: want{
				err: true,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			logger, err := New(tt.given.opts...)

			if tt.want.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, logger)
		})
	}
}
