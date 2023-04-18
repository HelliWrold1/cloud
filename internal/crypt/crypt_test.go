package crypt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Generate_SaltPwd(t *testing.T) {
	t.Run("GenerateSaltPwd", func(t *testing.T) {
		pwd, err := GenerateSaltPwd("admin")
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Println(pwd)
		assert.Nil(t, err)
		assert.NotNil(t, pwd)
	})

	t.Run("CheckSaltedPassword", func(t *testing.T) {
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
	})
}
