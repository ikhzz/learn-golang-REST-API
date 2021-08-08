package response

import "restTestOne/models"

type CreatePrResponse struct {
	Status string
	Data models.Product
}