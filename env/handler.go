package function

import (
	"os"
	"strings"
)

// Handle a serverless request
func Handle(req []byte) string {

	response := strings.Builder{}
	for _, e := range os.Environ() {
		response.WriteString(e)
		response.WriteString("\n")
	}

	return response.String()
}
