package video

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

type UseCase struct{}

func NewService() *UseCase {
	return &UseCase{}
}

func (uc *UseCase) UploadVideo(ctx context.Context, video *multipart.FileHeader, folder string) (entity.File, error) {
	if video == nil {
		return entity.File{}, errors.New("video not found")
	}

	ext := strings.ToLower(filepath.Ext(video.Filename))
	allowed := map[string]bool{
		".mp4": true,
		".mov": true,
		".avi": true,
		".mkv": true,
	}

	if !allowed[ext] {
		return entity.File{}, fmt.Errorf("invalid file extension: %s (allowed: .mp4, .mov, .avi, .mkv)", ext)
	}

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return entity.File{}, err
	}

	randomName := uuid.New().String() + filepath.Ext(video.Filename)
	dst := filepath.Join(folder, randomName)

	f, err := video.Open()
	if err != nil {
		return entity.File{}, err
	}
	defer f.Close()

	out, err := os.Create(dst)
	if err != nil {
		return entity.File{}, err
	}
	defer out.Close()

	if _, err := io.Copy(out, f); err != nil {
		return entity.File{}, err
	}

	info := entity.File{
		OriginalName: video.Filename,
		SavedName:    randomName,
		Size:         fmt.Sprintf("%d", video.Size),
		Path:         dst,
	}

	return info, nil
}

func (uc UseCase) Delete(ctx context.Context, url string) error {
	err := os.Remove("./" + url)
	return err
}

func (uc UseCase) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]entity.File, error) {
	var links []entity.File

	for _, f := range files {
		link, err := uc.UploadVideo(ctx, f, folder)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}
