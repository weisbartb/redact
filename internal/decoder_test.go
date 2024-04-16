package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecoders(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		s := NewInstructionScanner(`[test,"test2","\"test3\""]`)
		s.Scan()
		require.Equal(t, opCodeSet, s.firstOp.opCode)
		require.Equal(t, opCodeString, s.firstOp.children[0].opCode)
		require.Equal(t, "test", s.firstOp.children[0].value)
		require.Equal(t, opCodeString, s.firstOp.children[1].opCode)
		require.Equal(t, "test2", s.firstOp.children[1].value)
		require.Equal(t, opCodeString, s.firstOp.children[2].opCode)
		require.Equal(t, `"test3"`, s.firstOp.children[2].value)
	})
	t.Run("string", func(t *testing.T) {
		t.Run("escape-quoted", func(t *testing.T) {
			s := NewInstructionScanner(`"\"test3\""`)
			s.Scan()
			require.Equal(t, opCodeString, s.firstOp.opCode)
			require.Equal(t, `"test3"`, s.firstOp.value)
		})
		t.Run("quoted", func(t *testing.T) {
			s := NewInstructionScanner(`"test3"`)
			s.Scan()
			require.Equal(t, opCodeString, s.firstOp.opCode)
			require.Equal(t, `test3`, s.firstOp.value)
		})
		t.Run("raw", func(t *testing.T) {
			s := NewInstructionScanner(`test3`)
			s.Scan()
			require.Equal(t, opCodeString, s.firstOp.opCode)
			require.Equal(t, `test3`, s.firstOp.value)
		})

	})
}
