package database

import (
	"database/sql"
	"testing"

	"github.com/gabrielmq/gointensivo/internal/order/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	DB *sql.DB
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.NoError(err)

	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	s.DB = db

}

func (s *OrderRepositoryTestSuite) TearDownTest() {
	s.DB.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (s *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	s.NoError(err)
	s.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(s.DB)
	err = repo.Save(order)
	s.NoError(err)

	var orderResult entity.Order
	err = s.DB.QueryRow("Select id, price, tax, final_price from orders where id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	s.NoError(err)
	s.Equal(order.ID, orderResult.ID)
	s.Equal(order.Price, orderResult.Price)
	s.Equal(order.Tax, orderResult.Tax)
	s.Equal(order.FinalPrice, orderResult.FinalPrice)
}
