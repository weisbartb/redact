package redaction

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type userRecord struct {
	Username             string
	Email                string `redact:"~admin=star(4)"`
	Password             string `redact:"all=zero"`
	LastName             string `redact:"~admin=remove(1)"`
	FirstName            string
	InterfaceTest        any     `redact:"~admin=remove(1)"`
	InterfacePointerTest any     `redact:"~admin=remove(1)"`
	PointerTest          *string `redact:"~admin=star(1)"`
	nonExported          *string
}

type unionRecord struct {
	userRecord
}

type stackedRecord struct {
	userRecord
	Embedded userRecord
}

func strPointer(str string) *string {
	return &str
}

func TestRedactRecords(t *testing.T) {
	t.Run("*map[string]T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}
		var out = &map[string]userRecord{
			"Test": rec,
		}
		clean, err := RedactRecord(out)
		cleanDref := *clean
		require.NoError(t, err)
		require.Equal(t, "Test", cleanDref["Test"].Username)
		require.Equal(t, "test*********", cleanDref["Test"].Email)
		require.Equal(t, "", cleanDref["Test"].Password)
		require.Equal(t, "l", cleanDref["Test"].LastName)
		require.Equal(t, "fname", cleanDref["Test"].FirstName)
		require.Equal(t, "t******", *cleanDref["Test"].PointerTest)
		require.Equal(t, "t", cleanDref["Test"].InterfaceTest.(string))
		require.Equal(t, "t", *(cleanDref["Test"].InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, cleanDref["Test"].nonExported)
	})
	t.Run("map[string]T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}
		var out = map[string]userRecord{
			"Test": rec,
		}
		clean, err := RedactRecord(out)
		require.NoError(t, err)
		require.Equal(t, "Test", clean["Test"].Username)
		require.Equal(t, "test*********", clean["Test"].Email)
		require.Equal(t, "", clean["Test"].Password)
		require.Equal(t, "l", clean["Test"].LastName)
		require.Equal(t, "fname", clean["Test"].FirstName)
		require.Equal(t, "t******", *clean["Test"].PointerTest)
		require.Equal(t, "t", clean["Test"].InterfaceTest.(string))
		require.Equal(t, "t", *(clean["Test"].InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, clean["Test"].nonExported)
	})
	t.Run("[]T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}
		var out = []userRecord{
			rec,
		}
		clean, err := RedactRecord(out)
		require.NoError(t, err)
		require.Equal(t, "Test", clean[0].Username)
		require.Equal(t, "test*********", clean[0].Email)
		require.Equal(t, "", clean[0].Password)
		require.Equal(t, "l", clean[0].LastName)
		require.Equal(t, "fname", clean[0].FirstName)
		require.Equal(t, "t******", *clean[0].PointerTest)
		require.Equal(t, "t", clean[0].InterfaceTest.(string))
		require.Equal(t, "t", *(clean[0].InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, clean[0].nonExported)
	})
	t.Run("*[]T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}
		var out = []userRecord{
			rec,
		}
		refClean, err := RedactRecord(&out)
		clean := *refClean
		require.NoError(t, err)
		require.Equal(t, "Test", clean[0].Username)
		require.Equal(t, "test*********", clean[0].Email)
		require.Equal(t, "", clean[0].Password)
		require.Equal(t, "l", clean[0].LastName)
		require.Equal(t, "fname", clean[0].FirstName)
		require.Equal(t, "t******", *clean[0].PointerTest)
		require.Equal(t, "t", clean[0].InterfaceTest.(string))
		require.Equal(t, "t", *(clean[0].InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, clean[0].nonExported)
	})
	t.Run("[x]T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}
		var out = [1]userRecord{
			rec,
		}
		clean, err := RedactRecord(out)
		require.NoError(t, err)
		require.Equal(t, "Test", clean[0].Username)
		require.Equal(t, "test*********", clean[0].Email)
		require.Equal(t, "", clean[0].Password)
		require.Equal(t, "l", clean[0].LastName)
		require.Equal(t, "fname", clean[0].FirstName)
		require.Equal(t, "t******", *clean[0].PointerTest)
		require.Equal(t, "t", clean[0].InterfaceTest.(string))
		require.Equal(t, "t", *(clean[0].InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, clean[0].nonExported)
	})
	t.Run("*[x]T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}
		var out = [1]userRecord{
			rec,
		}
		refClean, err := RedactRecord(&out)
		clean := *refClean
		require.NoError(t, err)
		require.Equal(t, "Test", clean[0].Username)
		require.Equal(t, "test*********", clean[0].Email)
		require.Equal(t, "", clean[0].Password)
		require.Equal(t, "l", clean[0].LastName)
		require.Equal(t, "fname", clean[0].FirstName)
		require.Equal(t, "t******", *clean[0].PointerTest)
		require.Equal(t, "t", clean[0].InterfaceTest.(string))
		require.Equal(t, "t", *(clean[0].InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, clean[0].nonExported)
	})
	t.Run("T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}

		clean, err := RedactRecord(rec)
		require.NoError(t, err)
		require.Equal(t, "Test", clean.Username)
		require.Equal(t, "test*********", clean.Email)
		require.Equal(t, "", clean.Password)
		require.Equal(t, "l", clean.LastName)
		require.Equal(t, "fname", clean.FirstName)
		require.Equal(t, "t******", *clean.PointerTest)
		require.Equal(t, "t", clean.InterfaceTest.(string))
		require.Equal(t, "t", *(clean.InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, clean.nonExported)
	})
	t.Run("*T", func(t *testing.T) {
		rec := userRecord{
			Username:             "Test",
			Email:                "test@test.com",
			Password:             "testp",
			LastName:             "lname",
			FirstName:            "fname",
			InterfaceTest:        "test",
			PointerTest:          strPointer("testptr"),
			InterfacePointerTest: strPointer("testptr"),
			nonExported:          strPointer("tmp"),
		}

		clean, err := RedactRecord(&rec)
		require.NoError(t, err)
		require.Equal(t, "Test", clean.Username)
		require.Equal(t, "test*********", clean.Email)
		require.Equal(t, "", clean.Password)
		require.Equal(t, "l", clean.LastName)
		require.Equal(t, "fname", clean.FirstName)
		require.Equal(t, "t******", *clean.PointerTest)
		require.Equal(t, "t", clean.InterfaceTest.(string))
		require.Equal(t, "t", *(clean.InterfacePointerTest.(*string)))
		require.Equal(t, rec.nonExported, clean.nonExported)
	})

}
