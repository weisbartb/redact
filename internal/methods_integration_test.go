package internal_test

import (
	"github.com/stretchr/testify/require"
	"github.com/weisbartb/redact/internal"
	"testing"
)

func TestInstructionScanner_StarEvaluator(t *testing.T) {
	t.Run("inverse", func(t *testing.T) {
		t.Run("no argument parsing", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("~[admin,csr]=star")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			var tStr2 = "test"

			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "****", tStr)
			ok, err = eval(&tStr2, "csr")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "test", tStr2)
		})
		t.Run("single group", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("~admin=star")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			var tStr2 = "test"

			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "****", tStr)
			ok, err = eval(&tStr2, "admin")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "test", tStr2)
		})
		t.Run("empty argument parsing", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("~[admin,csr]=star()")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "****", tStr)
		})
	})
	t.Run("mixed", func(t *testing.T) {
		t.Run("allow admin but not csr", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("[~admin,csr]=star")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			var tStr2 = "test"
			var tStr3 = "test"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "****", tStr)
			ok, err = eval(&tStr2, "csr")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "****", tStr2)
			ok, err = eval(&tStr3, "admin")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "test", tStr3)
		})
	})
	t.Run("offsets", func(t *testing.T) {
		t.Run("last 4", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("[~admin,csr]=star(-4)")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "555-555-1234"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "********1234", tStr)
		})
		t.Run("first 3", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("[~admin,csr]=star(3)")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "555-555-1234"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "555*********", tStr)
		})
	})

}
func TestInstructionScanner_RemoveEvaluator(t *testing.T) {
	t.Run("inverse", func(t *testing.T) {
		t.Run("no argument parsing", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("~[admin,csr]=remove")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			var tStr2 = "test"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "", tStr)
			ok, err = eval(&tStr2, "csr")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "test", tStr2)
		})
		t.Run("single group", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("~admin=remove")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			var tStr2 = "test"

			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "", tStr)
			ok, err = eval(&tStr2, "admin")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "test", tStr2)
		})
		t.Run("empty argument parsing", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("~[admin,csr]=remove()")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "", tStr)
		})
	})
	t.Run("mixed", func(t *testing.T) {
		t.Run("allow admin but not csr", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("[~admin,csr]=remove")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "test"
			var tStr2 = "test"
			var tStr3 = "test"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "", tStr)
			ok, err = eval(&tStr2, "csr")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "", tStr2)
			ok, err = eval(&tStr2, "admin")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "test", tStr3)
		})
	})
	t.Run("offsets", func(t *testing.T) {
		t.Run("last 4", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("[~admin,csr]=remove(-4)")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "555-555-1234"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "1234", tStr)
		})
		t.Run("first 3", func(t *testing.T) {
			scanner := internal.NewInstructionScanner("[~admin,csr]=remove(3)")
			scanner.Scan()
			eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
				"star":   internal.MethodStar,
				"remove": internal.MethodRemove,
			})
			require.NoError(t, err)
			var tStr = "555-555-1234"
			ok, err := eval(&tStr, "user")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, "555", tStr)
		})
	})

}
func TestInstructionScanner_ZeroEvaluator(t *testing.T) {
	t.Run("inverse", func(t *testing.T) {
		scanner := internal.NewInstructionScanner("[~admin,csr]=zero()")
		scanner.Scan()
		eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
			"zero": internal.MethodZero,
		})
		require.NoError(t, err)
		var tStr = "555-555-1234"
		ok, err := eval(&tStr, "user")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "", tStr)

	})

}
func TestInstructionScanner_RedactEvaluator(t *testing.T) {
	t.Run("inverse", func(t *testing.T) {
		scanner := internal.NewInstructionScanner(`[~admin,csr]=redact(*,"-")`)
		scanner.Scan()
		eval, err := scanner.GetEvaluator(map[string]internal.RawMethod{
			"redact": internal.MethodRedact,
		})
		require.NoError(t, err)
		var tStr = "555-555-1234"
		ok, err := eval(&tStr, "user")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "***-***-****", tStr)

	})

}
