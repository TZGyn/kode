package tools

import (
	"context"
	"testing"

	"charm.land/fantasy"
	"github.com/stretchr/testify/require"
)

func TestViewTool(t *testing.T) {
	t.Parallel()

	t.Run("Run View Tool", func(t *testing.T) {
		t.Parallel()

		tool := NewViewTool()

		_, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "test",
			Name:  "test",
			Input: "{\"file_path\":\"./view.md\",\"offset\":0,\"limit\":20}"},
		)

		require.NoError(t, err)

		// t.Log(response.Data)
	})
}
