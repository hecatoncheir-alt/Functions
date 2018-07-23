package function

import (
	"fmt"
)

func Handle(req []byte) string {
	return fmt.Sprintf("Input: %s", string(req))
}
