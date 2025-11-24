package services

import (
	"fmt"

	"github.com/greatdaveo/XChange/rates"
)

// To hold one or more rate providers
type RateService struct {
	Providers []rates.RateProvider
}

func NewRateService(providers []rates.RateProvider) *RateService {
	return &RateService{Providers: providers}
}

// To make GetRate fetch rate from the first working provider
func (rs *RateService) GetRate(from, to string) (float64, error) {
	for _, provider := range rs.Providers {
		rate, err := provider.FetchRate(from, to)
		fmt.Println("provider.GetRate: rate =", rate, "err =", err)

		if err == nil && rate > 0 {
			fmt.Println("err loop: ", err)
			return rate, nil
		}
	}

	return 0, fmt.Errorf("could not fetch exchange rate from any provider")
}
