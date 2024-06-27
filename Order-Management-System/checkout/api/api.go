package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	pgx "github.com/jackc/pgx/v4"
	pgxpool "github.com/jackc/pgx/v4/pgxpool"

	"github.com/sofiukl/oms/oms-core/models"
	"github.com/sofiukl/oms/oms-core/utils"

	"github.com/mitchellh/mapstructure"
)

// CheckoutProduct implements the checkout process. Steps as follows -
// 1. Finds the products against cart [call to cart microservice]
// 2. Finds the product details [call to product microservice]
// 3. Allocates the quantity to the customer after successfull payment.
// Uses the concept of reserve quantity to race condition [assumes dummy payment microservice]
// Please not this implementation acruire lock also for handling race condition
func CheckoutProduct(conn *pgxpool.Pool, config utils.Config, body models.CheckoutModel, lock *sync.RWMutex) {

	// begin transaction eparate this as util
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback(context.Background())

	amount := body.Amount
	product, err := findCartDetails(body.CartID)
	if err != nil {
		log.Println(err)
		return
	}

	msg, checkoutErr := processTransaction(tx, product, amount, lock)
	if checkoutErr != nil {
		log.Println(checkoutErr)
		return
	}
	log.Println(msg)

	// commit transaction
	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func processTransaction(tx pgx.Tx, product models.ProductModel, amount float64, lock *sync.RWMutex) (string, error) {

	id := product.ID
	qty := product.Quantity

	lock.Lock()
	defer lock.Unlock()
	prod, err := findProduct(id)
	log.Printf("%+v", prod)

	if err != nil {
		log.Println(err)
		return "Fail to checkout at this moment", fmt.Errorf("Fail to checkout st this moment")
	}

	// check quantity is avilable or not
	if prod.AvailQty-prod.ReserveQty < qty {
		log.Println("out of stock")
		return "The product is out of stock at this moment", nil
	}

	updateReserveQty(tx, id, qty, "increase")
	err = payment(id, amount)
	if err != nil {
		updateReserveQty(tx, id, qty, "decrease")
		return "Your payment is not successfull", fmt.Errorf("Your payment is not successfull")
	}
	updateAvailQty(tx, id, qty)
	updateReserveQty(tx, id, qty, "decrease")
	return "Yup! you successfully bought the product", nil

}

func updateReserveQty(tx pgx.Tx, prodID string, qty int, opType string) error {
	// query need to separated in other file
	var qry string
	if opType == "increase" {
		qry = fmt.Sprintf("update product set reserve_qty = reserve_qty + %d where id =$1", qty)
	} else {
		qry = fmt.Sprintf("update product set reserve_qty = reserve_qty - %d where id =$1", qty)
	}
	_, err := tx.Exec(context.Background(), qry, prodID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return fmt.Errorf("update quantity failed")
	}

	return nil
}

func updateAvailQty(tx pgx.Tx, prodID string, qty int) error {
	// query need to separated in other file
	qry := fmt.Sprintf("update product set avail_qty = avail_qty - %d where id =$1", qty)
	_, err := tx.Exec(context.Background(), qry, prodID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return fmt.Errorf("update avail quantity failed")
	}
	return nil
}

// dummy stuff for payment process
func payment(prodID string, amount float64) error {
	time.Sleep(2000)
	return nil
}

func findProduct(id string) (*models.Product, error) {
	var p models.Product
	var gresp models.GenericResponse
	link := fmt.Sprintf("http://localhost:3004/product/api/v1/find/%s", id) // this should be in env file
	resp, err := http.Get(link)
	if err != nil {
		log.Println("Fail to fetch product details from product service")
		return &p, fmt.Errorf("Fail to fetch product at this moment")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &gresp); err != nil {
		log.Println(err)
		return &p, fmt.Errorf("Fail to fetch product at this moment")
	}
	mapstructure.Decode(gresp.Result, &p)
	return &p, nil
}

func findCartDetails(id string) (models.ProductModel, error) {
	var gresp models.GenericResponse
	var cart models.CartModel
	var prod models.ProductModel

	link := fmt.Sprintf("http://localhost:3006/cart/api/v1/find/%s", id) // this should be in env file
	resp, err := http.Get(link)
	if err != nil {
		log.Println("Fail to fetch product details from product service")
		return prod, fmt.Errorf("Fail to fetch cart details")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &gresp); err != nil {
		log.Println(err)
		return prod, fmt.Errorf("Fail to fetch cart details")
	}
	mapstructure.Decode(gresp.Result, &cart)
	return cart.Products[0], nil // assumes only 1 product can be in cart for simplicity
}
