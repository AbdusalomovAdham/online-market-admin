package product

import (
	"context"
	"main/internal/entity"
	"mime/multipart"
)

type Service struct {
	repo Repository
	auth Auth
	file File
}

func NewService(repo Repository, auth Auth, file File) Service {
	return Service{repo: repo, auth: auth, file: file}
}

func (s Service) CreateProduct(ctx context.Context, data Create, authHeader string) (int64, error) {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, data, isValidToken.Id)
}

func (s Service) MultipleUpload(ctx context.Context, files []*multipart.FileHeader, folder string, startID *int32) ([]entity.File, error) {
	return s.file.MultipleUpload(ctx, files, folder, startID)
}

func (s Service) GetById(ctx context.Context, id int64) (GetById, error) {
	data, err := s.repo.GetById(ctx, id)
	if err != nil {
		return GetById{}, err
	}

	return data, nil
}

func (s Service) GetList(ctx context.Context, filter entity.Filter) ([]Get, int, error) {
	data, count, err := s.repo.GetList(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, count, nil
}

func (s Service) UpdateProduct(ctx context.Context, productId int, data Update, authHeader string, images []*multipart.FileHeader) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	oldProduct, err := s.repo.GetByIdDetail(ctx, int64(productId))
	if err != nil {
		return err
	}

	if oldProduct.Images == nil {
		oldProduct.Images = &[]entity.File{}
	}

	if data.RemoveImages != nil {
		removeMap := make(map[int32]bool)
		for _, id := range *data.RemoveImages {
			removeMap[id] = true
		}

		var newImages []entity.File

		for _, image := range *oldProduct.Images {
			if removeMap[image.Id] {
				err := s.file.Delete(ctx, image.Path)
				if err != nil {
					return err
				}
				continue
			}
			newImages = append(newImages, image)
		}
		oldProduct.Images = &newImages
	}

	var maxId int32 = 0

	if oldProduct.Images != nil {
		for _, img := range *oldProduct.Images {
			if img.Id > maxId {
				maxId = img.Id
			}
		}
	}

	if len(images) > 0 {
		startId := maxId + 1
		imgFile, err := s.file.MultipleUpload(ctx, images, "../media/products", &startId)
		if err != nil {
			return err
		}
		*oldProduct.Images = append(*oldProduct.Images, imgFile...)
	}

	if oldProduct.Images == nil {
		data.Images = &[]entity.File{}
	} else {
		data.Images = oldProduct.Images
	}

	return s.repo.UpdateProduct(ctx, productId, data, isValidToken.Id)
}

func (s Service) AdminDeleteProduct(ctx context.Context, productId int64, authHeader string) error {
	isValidToken, err := s.auth.IsValidToken(ctx, authHeader)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, productId, isValidToken.Id)
}
