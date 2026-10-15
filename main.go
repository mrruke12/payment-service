package main

import (
	"context"
	"fmt"
	"sync"
	"uuid"

	"github.com/Masterminds/squirrel"
	"github.com/mrruke12/payment-service/internal/application/orderusecase"
	"github.com/mrruke12/payment-service/internal/domains/order"
	"github.com/mrruke12/payment-service/internal/infrastructure/config"
	"github.com/mrruke12/payment-service/internal/infrastructure/postgres"
)

func main() {
	config.LoadEnvironment()
	cfg, err := config.GetDBConfig()

	if err != nil {
		panic(err)
	}

	err = postgres.RunMigrations(*cfg)

	fmt.Println(err)

	db := postgres.Connect(context.Background(), *cfg)

	seed(db)

	uow := postgres.NewCreateOrderUnitOfWork(db)

	s := orderusecase.NewCreateOrderService(uow)

	items := []order.Item{
		{
			ProductID: uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c01"),
			Quantity:  1,
		},
	}
	var userID = uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c0c")

	o, err := s.CreateOrder(context.Background(), userID, items)

	if err != nil {
		panic(err)
	}

	fmt.Println(o)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				var userID = uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c0c")

				items := []order.Item{
					{
						ProductID: uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c01"),
						Quantity:  1,
					},
				}

				o, err := s.CreateOrder(context.Background(), userID, items)
				if err != nil {
					fmt.Printf("worker %d iter %d: %v\n", worker, j, err)
					continue
				}
				_ = o
			}
		}(i)
	}
	wg.Wait()

}

func seed(db postgres.Querier) {
	ctx := context.Background()

	var userID = uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c0c")

	productIDs := []uuid.UUID{
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c01"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c02"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c03"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c04"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c05"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c06"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c07"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c08"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c09"),
		uuid.MustParse("01a130d5-8ca9-7746-8b2f-e1065aad6c0a"),
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	ins := builder.Insert("users").Columns("id", "full_name").Values(userID, "test user")
	sql, args, err := ins.ToSql()
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(ctx, sql, args...); err != nil {
		return
	}

	ins = builder.Insert("products").Columns("id", "name", "quantity", "price")
	for i, id := range productIDs {
		ins = ins.Values(id, fmt.Sprintf("product %d", i+1), 10_000, 1000)
	}
	sql, args, err = ins.ToSql()
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(ctx, sql, args...); err != nil {
		return
	}
}
