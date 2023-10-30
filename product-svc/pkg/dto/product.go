package dto

type ProcessProductRequestDto struct {
	ProductId int64 `json:"product_id" validate:"required"`
	Quantity  int32 `json:"quantity" validate:"required"`
}
