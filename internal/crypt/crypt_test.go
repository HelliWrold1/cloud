package crypt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSaltPwd(t *testing.T) {
	t.Run("GenerateSaltPwd", func(t *testing.T) {
		pwd, err := GenerateSaltPwd("admin")
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Println(pwd)
		// success
		err = CheckSaltPwd("admin", pwd)
		if err != nil {
			t.Fatal(err.Error())
		}
		// failed
		err = CheckSaltPwd("adm", pwd)
		if err != nil {
			assert.NotNil(t, err)
		}

		err = CheckSaltPwd("admin", "$2a$10$84IgaDinyGP0bAiGdfiHXux/5rWTdbVT/N3p7TsPU9sBV/yXlVXoy")
		assert.Nil(t, err)
	})
}
