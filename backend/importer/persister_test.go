package importer

import (
	"testing"
)

import "../data"

func TestSaveAdvertiser(t *testing.T) {
	wasSaved := false

	sut := PersisterImpl{
		adRepo: FakeAdRepo{
			onSave: func(ad data.AdvertiserEntity) {
				wasSaved = true

				if ad.Id != -1 {
					t.Error("Should be empty value", ad.Id)
				}
				if ad.Name != "google" {
					t.Error("Wrong name", ad.Name)
				}
			},
		},
	}

	sut.SaveAdvertiser(Advertiser{Name: "google"})

	if !wasSaved {
		t.Error("Wasn't saved")
	}
}

func TestSaveProduct_NoAd(t *testing.T) {
	sut := PersisterImpl{
		adRepo: FakeAdRepo{ad: nil},

		productRepo: FakeProdRepo{
			onSave: func(prod data.ProductEntity) error {
				t.Error("Should not save")

				return nil
			},
		},
	}

	err := sut.SaveProduct(Product{
		Sku:        "123",
		Name:       "iphone",
		Advertiser: "google",
	})

	if err != ErrUnknownAdvertiser {
		t.Error("Expected unknown advertiser error, but got", err)
	}
}

func TestSaveProduct_SomeAd(t *testing.T) {
	wasSaved := false

	sut := PersisterImpl{
		adRepo: FakeAdRepo{ad: &data.AdvertiserEntity{Id: 123, Name: "google"}},

		productRepo: FakeProdRepo{onSave: func(prod data.ProductEntity) error {
			wasSaved = true
			if prod.Sku != "123" || prod.Name != "iphone" || prod.AdvertiserId != 123 {
				t.Error("Unexpected entity saved", prod)
			}

			return nil
		}},
	}

	err := sut.SaveProduct(Product{
		Sku:        "123",
		Name:       "iphone",
		Advertiser: "google",
	})

	if err != nil {
		t.Error("Unexpected error", err)
	}
	if !wasSaved {
		t.Error("Wasn't saved")
	}
}

func TestSaveProduct_Duplicate(t *testing.T) {
	sut := PersisterImpl{
		adRepo: FakeAdRepo{ad: &data.AdvertiserEntity{Id: 123, Name: "google"}},

		productRepo: FakeProdRepo{
			onSave: func(prod data.ProductEntity) error {
				return data.ErrDuplicateProduct
			},
		},
	}

	err := sut.SaveProduct(Product{
		Sku:        "123",
		Name:       "iphone",
		Advertiser: "google",
	})

	if err != data.ErrDuplicateProduct {
		t.Error("Expected duplicate product error, but got", err)
	}
}

type FakeAdRepo struct {
	ad     *data.AdvertiserEntity
	onSave func(ad data.AdvertiserEntity)
}

func (repo FakeAdRepo) GetByName(name string) *data.AdvertiserEntity {
	return repo.ad
}

func (repo FakeAdRepo) Save(ad data.AdvertiserEntity) {
	repo.onSave(ad)
}

type FakeProdRepo struct {
	onSave func(prod data.ProductEntity) error
}

func (repo FakeProdRepo) Save(prod data.ProductEntity) error { return repo.onSave(prod) }
