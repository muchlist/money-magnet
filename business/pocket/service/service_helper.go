package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/category/model"
)

func generateDefaultCategory(pocketID uuid.UUID) []model.Category {
	timeNow := time.Now()
	categories := []model.Category{
		/*INCOME CATEGORY*/
		{
			PocketID:     pocketID,
			CategoryName: "Salary",
			CategoryIcon: 1,
			IsIncome:     true,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Grants",
			CategoryIcon: 2,
			IsIncome:     true,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Refunds",
			CategoryIcon: 3,
			IsIncome:     true,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Sale",
			CategoryIcon: 4,
			IsIncome:     true,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Rental",
			CategoryIcon: 5,
			IsIncome:     true,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		/*EXPENSE CATEGORY*/
		{
			PocketID:     pocketID,
			CategoryName: "Baby",
			CategoryIcon: 6,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Beauty",
			CategoryIcon: 7,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Bills",
			CategoryIcon: 8,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Vehicle",
			CategoryIcon: 9,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Clothing",
			CategoryIcon: 10,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Education",
			CategoryIcon: 11,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Electronics",
			CategoryIcon: 12,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Entertainment",
			CategoryIcon: 13,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Food",
			CategoryIcon: 14,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Health",
			CategoryIcon: 15,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Home",
			CategoryIcon: 16,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Insurance",
			CategoryIcon: 17,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Shopping",
			CategoryIcon: 18,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Social",
			CategoryIcon: 19,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Sport",
			CategoryIcon: 20,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Tax",
			CategoryIcon: 21,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Telephone",
			CategoryIcon: 22,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Internet",
			CategoryIcon: 23,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Transportation",
			CategoryIcon: 24,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
		{
			PocketID:     pocketID,
			CategoryName: "Work",
			CategoryIcon: 25,
			IsIncome:     false,
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
		},
	}

	return categories
}
