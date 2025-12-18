package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"main/internal/entity"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type UseCase struct {
}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (uc *UseCase) Upload(ctx context.Context, image *multipart.FileHeader, folder string) (entity.File, error) {
	if image == nil {
		return entity.File{}, errors.New("image not found")
	}

	mimeType := image.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		return entity.File{}, fmt.Errorf("invalid file type: expected image, got %s", mimeType)
	}

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return entity.File{}, err
	}

	randomName := uuid.New().String() + filepath.Ext(image.Filename)
	dst := filepath.Join(folder, randomName)

	src, err := image.Open()
	if err != nil {
		return entity.File{}, err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return entity.File{}, err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return entity.File{}, err
	}

	info := entity.File{
		OriginalName: image.Filename,
		SavedName:    randomName,
		Size:         fmt.Sprintf("%d", image.Size),
		Path:         dst,
	}

	return info, nil
}

func (uc *UseCase) Delete(ctx context.Context, url string) error {
	err := os.Remove("./" + url)
	return err
}

func (uc *UseCase) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]entity.File, error) {
	var links []entity.File
	count := len(files)
	for i, f := range files {
		link, err := uc.Upload(ctx, f, folder)
		if err != nil {
			return nil, err
		}
		if count > i {
			link.Id = int32(i + 1)
			links = append(links, link)
		}
	}

	return links, nil
}
