package services

import (
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"mirabilis-api/src/config"
	"sync"
)

type CloudinaryService struct{}

func NewCloudinaryService() *CloudinaryService {
	return &CloudinaryService{}
}

func (this *CloudinaryService) UploadFile(ctx *gin.Context, file multipart.File, folder string, resourceType string) <-chan map[string]any {
	resultChan := make(chan map[string]any)

	go func() {
		defer close(resultChan)

		uploadParams := uploader.UploadParams{Folder: folder, ResourceType: resourceType}
		cld, err := config.ConnectCloudinary()
		if err != nil {
			log.Println(err.Error())
			resultChan <- map[string]any{"error": true, "data": nil}
			return
		}

		result, err := cld.Upload.Upload(ctx, file, uploadParams)
		if err != nil {
			log.Println(err.Error())
			resultChan <- map[string]any{"error": true, "data": nil}
			return
		}

		resultChan <- map[string]any{"error": false, "data": result}
	}()
	return resultChan
}

func (this *CloudinaryService) UploadFiles(ctx *gin.Context, files []multipart.File, folder, resourceType string) <-chan []map[string]any {
	resultChan := make(chan []map[string]any)

	go func() {
		defer close(resultChan)
		var wg sync.WaitGroup
		results := make([]map[string]any, len(files))
		mu := sync.Mutex{}

		for i, file := range files {
			wg.Add(1)

			go func(index int, f multipart.File) {
				defer wg.Done()

				uploadParams := uploader.UploadParams{
					Folder:       folder,
					ResourceType: resourceType,
				}

				cld, err := config.ConnectCloudinary()
				if err != nil {
					log.Println(err.Error())
					mu.Lock()
					results[index] = map[string]any{"error": true, "data": nil}
					mu.Unlock()
					return
				}

				result, err := cld.Upload.Upload(ctx, f, uploadParams)
				if err != nil {
					log.Println(err.Error())
					mu.Lock()
					results[index] = map[string]any{"error": true, "data": nil}
					mu.Unlock()
					return
				}

				mu.Lock()
				results[index] = map[string]any{"error": false, "data": result}
				mu.Unlock()
			}(i, file)
		}

		wg.Wait()
		resultChan <- results
	}()

	return resultChan
}
