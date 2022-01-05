package suite_test

import (
	"log"
	"os"
	"testing"

	"github.com/zpatrick/testx/suite"
)

type UserSuite struct {
	UserID int
}

func (u *UserSuite) Setup() error {
	log.Println("[UserSuite] running setup")
	u.UserID = 1
	return nil
}

func (u *UserSuite) Teardown() error {
	log.Println("[UserSuite] running teardown")
	return nil
}

type ProductSuite struct {
	suite.Base

	ProductName string
}

func (p *ProductSuite) Setup() error {
	p.ProductName = "shampoo"
	return nil
}

func TestMain(m *testing.M) {
	suite.Register(&UserSuite{})
	suite.Register(&ProductSuite{})

	os.Exit(suite.Run(m))
}

func TestUserSuite_alpha(t *testing.T) {
	u := suite.Get[*UserSuite](t)
	t.Log("user id:", u.UserID)
}

func TestUserSuite_bravo(t *testing.T) {
	u := suite.Get[*UserSuite](t)
	t.Log("user id:", u.UserID)
}

func TestUserSuite_charlie(t *testing.T) {
	u := suite.Get[*UserSuite](t)
	t.Log("user id:", u.UserID)
}

func TestOtherSuite_alpha(t *testing.T) {
	p := suite.Get[*ProductSuite](t)
	t.Log("product name:", p.ProductName)
}
