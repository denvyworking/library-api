package service

import (
	"context"
	"fmt"
	"leti/pkg/models"
	"leti/pkg/repository/fake"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateBook(t *testing.T) {
	names := []string{"451 градус по Фаренгейту", "I am a bad test!", "В аптеке"}
	prices := []int{132, 100, -132}
	testCases := []struct {
		TestName string
		Name     *string
		Price    *int
		Id       int
		Correct  bool
	}{
		{
			TestName: "test1|ok",
			Name:     &names[0],
			Price:    &prices[0],
			Id:       1,
			Correct:  true,
		},
		{
			TestName: "test2|wrong ID",
			Name:     &names[1],
			Price:    &prices[1],
			Id:       -1,
			Correct:  false,
		},
		{
			TestName: "test3|nil",
			Name:     nil,
			Price:    nil,
			Id:       1,
			Correct:  true,
		},
		{
			TestName: "test4|wrong price",
			Name:     &names[2],
			Price:    &prices[2],
			Id:       1,
			Correct:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			var fakeDB = &fake.FakeRepo{}
			svc := NewService(fakeDB)
			id, _ := svc.CreateBook(context.Background(), models.Book{
				Name:      "Онегин",
				Author_id: 1,
				Genre_id:  1,
				Price:     100,
			})
			if tc.Id == -1 {
				err := svc.UpdateBook(context.Background(), tc.Id, models.BookUpdate{
					Name:  tc.Name,
					Price: tc.Price,
				})
				require.Error(t, err)
				require.Contains(t, err.Error(), fmt.Sprintf("book with id %d not found", tc.Id))
			} else {
				err := svc.UpdateBook(context.Background(), id, models.BookUpdate{
					Name:  tc.Name,
					Price: tc.Price,
				})
				if tc.Correct == false {
					require.Error(t, err)
					require.Contains(t, err.Error(), "price must be non-negative")
				} else {
					require.NoError(t, err)

				}
			}
		})
	}

}
