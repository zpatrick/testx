package suite_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/zpatrick/testx/suite"
)

type DBSuite struct {
	db *sql.DB
}

func (d *DBSuite) Setup(tb testing.TB) error {
	log.Println("[DBSuite] running setup")
	d.db = &sql.DB{}
	return nil
}

func (d *DBSuite) Exec(query string, args ...any) error {
	log.Println("[DBSuite] running query: ", query)
	return nil
}

func (d *DBSuite) Teardown() error {
	log.Println("[DBSuite] running teardown")

	return nil
}

type UserSuite struct {
	dbSuite *DBSuite

	UserID int
}

func (u *UserSuite) Setup(tb testing.TB) error {
	log.Println("[UserSuite] running setup")

	u.dbSuite = suite.Get[*DBSuite](tb)
	if err := u.dbSuite.Exec("INSERT INTO users ..."); err != nil {
		return err
	}

	u.UserID = 1
	return nil
}

func (u *UserSuite) Teardown() error {
	log.Println("[UserSuite] running teardown")
	return u.dbSuite.Exec("DELETE FROM users WHERE id=?", u.UserID)
}

type ProductSuite struct {
	suite.Base

	dbSuite     *DBSuite
	ProductID   int
	ProductName string
}

func (p *ProductSuite) Setup(tb testing.TB) error {
	log.Println("[ProductSuite] running setup")

	p.dbSuite = suite.Get[*DBSuite](tb)
	if err := p.dbSuite.Exec("INSERT INTO products ..."); err != nil {
		return err
	}

	p.ProductID = 1
	return nil
}

func (p *ProductSuite) Teardown() error {
	log.Println("[ProductSuite] running teardown")
	return p.dbSuite.Exec("DELETE FROM products WHERE id=?", p.ProductID)
}

func TestMain(m *testing.M) {
	// Teardowns happen in FILO order, so we
	// register the db suite first since it should close after we
	// cleanup our test user and product.
	suite.Register(&DBSuite{})
	suite.Register(&UserSuite{})
	suite.Register(&ProductSuite{ProductName: "Shampoo"})

	os.Exit(suite.Run(m))
}

func TestUserSuite_alpha(t *testing.T) {
	t.Parallel()

	u := suite.Get[*UserSuite](t)
	t.Log("user id:", u.UserID)
}

func TestUserSuite_bravo(t *testing.T) {
	t.Parallel()

	u := suite.Get[*UserSuite](t)
	t.Log("user id:", u.UserID)
}

func TestUserSuite_charlie(t *testing.T) {
	t.Parallel()

	u := suite.Get[*UserSuite](t)
	t.Log("user id:", u.UserID)
}

func TestOtherSuite_alpha(t *testing.T) {
	t.Parallel()

	p := suite.Get[*ProductSuite](t)
	t.Log("product name:", p.ProductName)
}
