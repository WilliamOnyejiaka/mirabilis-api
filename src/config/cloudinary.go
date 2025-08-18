package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

func ConnectCloudinary() (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromParams(Envs("cloudinaryCloudName"), Envs("cloudinaryApiKey"), Envs("cloudinaryApiSecret"))
}
