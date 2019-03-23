package importer

import (
	"errors"

	"../data"
)

type Persister interface {
	SaveAdvertiser(Advertiser)
	SaveProduct(Product) error
}

type PersisterImpl struct {
	adRepo      data.AdvertiserRepository
	productRepo data.ProductRepository
}

func NewPersister(adRepo data.AdvertiserRepository, productRepo data.ProductRepository) Persister {
	return PersisterImpl{adRepo, productRepo}
}

func (p PersisterImpl) SaveAdvertiser(ad Advertiser) {
	p.adRepo.Save(data.AdvertiserEntity{Id: -1, Name: ad.Name})
}

func (p PersisterImpl) SaveProduct(product Product) error {
	ad := p.adRepo.GetByName(product.Advertiser)

	if ad == nil {
		return ErrUnknownAdvertiser
	}

	p.productRepo.Save(data.ProductEntity{
		Sku:          product.Sku,
		Name:         product.Name,
		AdvertiserId: ad.Id,
	})

	return nil
}

var ErrUnknownAdvertiser = errors.New("Uknown advertiser")
