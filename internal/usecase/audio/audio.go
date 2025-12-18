package audio

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

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) UploadAudio(ctx context.Context, audio *multipart.FileHeader, folder string) (entity.File, error) {
	if audio == nil {
		return entity.File{}, errors.New("video not found")
	}

	contentType := audio.Header.Get("Content-type")
	if !strings.HasPrefix(contentType, "audio/") {
		return entity.File{}, fmt.Errorf("invalid file type: %s (only video files are allowed)", contentType)
	}

	ext := strings.ToLower(filepath.Ext(audio.Filename))
	allowed := map[string]bool{
		".mp3": true,
		".wav": true,
		".ogg": true,
		".m4a": true,
	}
	if !allowed[ext] {
		return entity.File{}, fmt.Errorf("invalid file extension: %s (allowed: .mp3, .wav, .ogg, .m4a)", ext)
	}

	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return entity.File{}, err
	}

	randomName := uuid.New().String() + filepath.Ext(audio.Filename)
	dst := filepath.Join(folder, randomName)

	f, err := audio.Open()
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
		OriginalName: audio.Filename,
		SavedName:    randomName,
		Size:         fmt.Sprintf("%d", audio.Size),
		Path:         dst,
	}

	return info, nil
}

func (s Service) Delete(url string) error {
	err := os.Remove("./" + url)
	return err
}

func (s Service) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string) ([]entity.File, error) {
	var links []entity.File

	for _, f := range files {
		link, err := s.UploadAudio(ctx, f, folder)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}
