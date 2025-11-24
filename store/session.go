package store

import "github.com/greatdaveo/SendlyPay/services"

// Auth Access Token from True Layer
var AccessTokenStore = make(map[string]services.OAuthToken)
