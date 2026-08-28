package valkeycompat

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/valkey-io/valkey-go/mock"
	"go.uber.org/mock/gomock"
)

func TestSetFromBuffer(t *testing.T) {
	t.Run("valid binary payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewClient(ctrl)
		adapter := NewAdapter(client)

		key := "test-set-from-buffer-binary"
		payload := []byte{0x00, 0x01, 0xFF, 0xFE, 0xFD}

		client.EXPECT().Do(gomock.Any(), mock.MatchFn(func(cmd []string) bool {
			if len(cmd) != 3 || cmd[0] != "SET" || cmd[1] != key {
				return false
			}
			if !bytes.Equal([]byte(cmd[2]), payload) {
				return false
			}
			return true
		}, "SET key payload")).Return(
			mock.Result(mock.ValkeyString("OK")),
		)

		res := adapter.SetFromBuffer(context.Background(), key, payload)
		if err := res.Err(); err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if got := res.Val(); got != "OK" {
			t.Fatalf("expected status OK, got %q", got)
		}
	})

	t.Run("empty buffer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewClient(ctrl)
		adapter := NewAdapter(client)

		key := "test-set-from-buffer-empty"
		payload := []byte{}

		client.EXPECT().Do(gomock.Any(), mock.MatchFn(func(cmd []string) bool {
			return len(cmd) == 3 && cmd[0] == "SET" && cmd[1] == key && cmd[2] == ""
		}, "SET key empty-string")).Return(mock.Result(mock.ValkeyString("OK")))

		if err := adapter.SetFromBuffer(context.Background(), key, payload).Err(); err != nil {
			t.Fatalf("expected nil error on empty buffer, got %v", err)
		}
	})

	t.Run("large payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewClient(ctrl)
		adapter := NewAdapter(client)

		key := "test-set-from-buffer-large"
		payload := []byte(strings.Repeat("a", 1<<20))

		client.EXPECT().Do(gomock.Any(), mock.MatchFn(func(cmd []string) bool {
			if len(cmd) != 3 || cmd[0] != "SET" || cmd[1] != key {
				return false
			}
			return len(cmd[2]) == len(payload) && cmd[2] == string(payload)
		}, "SET key large payload")).Return(mock.Result(mock.ValkeyString("OK")))

		if err := adapter.SetFromBuffer(context.Background(), key, payload).Err(); err != nil {
			t.Fatalf("expected nil error for large payload, got %v", err)
		}
	})

	t.Run("context cancellation propagates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewClient(ctrl)
		adapter := NewAdapter(client)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		client.EXPECT().Do(gomock.Any(), gomock.Any()).Return(mock.ErrorResult(context.Canceled))

		err := adapter.SetFromBuffer(ctx, "test-set-from-buffer-cancel", []byte("payload")).Err()
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}
