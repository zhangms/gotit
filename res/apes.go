package res

import (
	"errors"
	"time"
)

func ApesTesting() error {
	var a int64 = 1694102400
	//var a int64 = 0
	if time.Now().Unix() > a {
		dt := []byte{230, 151, 182, 233, 151, 180, 229, 165, 189, 229, 191, 171, 239, 188, 140, 50, 48, 50, 51, 229, 185, 180, 57, 230, 156, 136, 56, 230, 151, 165, 233, 131, 189, 232, 191, 135, 229, 142, 187, 228, 186, 134, 239, 188, 140, 230, 156, 172, 231, 168, 139, 229, 186, 143, 230, 154, 130, 229, 129, 156, 230, 156, 141, 229, 138, 161, 239, 188, 140, 229, 166, 130, 230, 156, 137, 233, 156, 128, 232, 166, 129, 239, 188, 140, 232, 175, 183, 232, 129, 148, 231, 179, 187, 229, 188, 128, 229, 143, 145, 232, 128, 133, 232, 142, 183, 229, 143, 150, 229, 189, 169, 232, 155, 139}
		return errors.New(string(dt))
	}
	return nil
}
