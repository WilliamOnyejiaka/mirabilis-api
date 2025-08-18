package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strconv"
)

func ValidateFile(allowedTypes []string, maxFileSize int64, formName string, required bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, multipart, err := ctx.Request.FormFile(formName)

		if err != nil {
			if required {
				log.Println("File error:", err.Error())
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":   true,
					"message": "File missing",
				})
				return
			}
			//* If not required and no file, just continue
			ctx.Next()
			return
		}

		fileType := multipart.Header.Get("Content-Type")

		if !slices.Contains(allowedTypes, fileType) {
			defer file.Close()
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid image type",
			})
			return
		}

		//* Validate file size (eg, max 4MB)
		if multipart.Size > maxFileSize {
			//* Convert to string
			maxFileSizeStr := strconv.FormatInt(maxFileSize, 10) // base 10

			message := fmt.Sprintf("File is too large. Max size is %sMB.", maxFileSizeStr)

			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": message,
			})
			return
		}

		defer file.Close()

		ctx.Next()
	}
}
