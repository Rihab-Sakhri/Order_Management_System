package core

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sofiukl/oms/oms-core/utils"
	"github.com/sofiukl/oms/oms-product/api"

	pgxpool "github.com/jackc/pgx/v4/pgxpool"
)

// App - Application
type App struct {
	Router *mux.Router
	Conn   *pgxpool.Pool
	Config utils.Config
}

// Initialize - This function initializes the application
func (a *App) Initialize() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.Connect(context.Background(), config.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	dbConnectMsg := fmt.Sprintf("Connected to DB %s", config.DBURL)
	fmt.Println(dbConnectMsg)
	a.Conn = conn
	a.Router = mux.NewRouter()
	a.Config = config
	a.initializeRoutes()
}

// Run - This functio funs the application
func (a *App) Run(address string) {
	fmt.Println("Application is running on port", address)
	if err := http.ListenAndServe(address, a.Router); err != nil {
		log.Fatal(err)
	}
}

func (a *App) initializeRoutes() {
	s := a.Router.PathPrefix("/product/api/v1").Subrouter()
	s.HandleFunc("/find/{id}", a.findProduct).Methods("GET")
}

func (a *App) findProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	config := a.Config
	conn := a.Conn
	api.FindProduct(conn, config, idStr, w, r)
}
