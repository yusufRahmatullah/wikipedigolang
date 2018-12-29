package repeat

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/utils"
	"github.com/gin-gonic/gin"
)

// Handler handles repeat-specific request
func Handler(c *gin.Context) {
	repeatStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(repeatStr)
	utils.HandleError(
		err,
		fmt.Sprintf("Error converting $REPEAT to an int: %q - Using default\n", err),
		func() { repeat = 5 },
	)
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!\n")
	}
	c.String(http.StatusOK, buffer.String())
}
