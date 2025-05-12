package dtos

type CalculateTotalRequestDTO struct {
	DiscountTypesIds []int            `json:"discountTypesIds"`
	TaxTypesIds      []int            `json:"taxTypesIds"`
	ItemsDTO         []BillingItemDTO `json:"itemsDTO"`
}
