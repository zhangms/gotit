package res

import (
	"fmt"
	"testing"
)

func TestApesTesting(t *testing.T) {

	er := ApesTesting()
	if er != nil {
		fmt.Println(er)
	}

}
