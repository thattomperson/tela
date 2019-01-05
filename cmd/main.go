package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ttp.sh/tela/internal/eddn"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

var db *gorm.DB

type System struct {
	gorm.Model
	Name string

	// Stations []Station
}

type Station struct {
	gorm.Model
	Name     string
	SystemID uint

	System System
}

type Commodity struct {
	gorm.Model
	Name      string
	MeanPrice float64
}

type StationCommodity struct {
	gorm.Model

	CommodityID uint
	StationID   uint

	Commodity Commodity
	Station   Station

	BuyPrice      float64
	Stock         int
	SellPrice     float64
	Demand        int
	StockBracket  int
	DemandBracket int
}

func main() {
	var err error
	db, err = gorm.Open("mysql", "root@tcp(127.0.0.1:3306)/elsa?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.AutoMigrate(
		&System{},
		&Station{},
		&Commodity{},
		&StationCommodity{},
	)

	go func() {
		log.Fatal(eddn.Listen("tcp://eddn.edcd.io:9500", handler))
	}()

	r := mux.NewRouter()

	r.Methods("GET").Path("/").
		Handler(appHandler(indexHandler))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
	appengine.Main()
}

func handler(i interface{}) {
	switch v := i.(type) {
	case eddn.Commodity:
		commodity(v)
	}
}

func commodity(c eddn.Commodity) {
	// Get first matched record
	var system System
	var station Station

	db.Where(System{Name: c.Message.SystemName}).FirstOrCreate(&system)
	db.Where(Station{Name: c.Message.StationName, SystemID: system.ID}).FirstOrCreate(&station)
	for _, a := range c.Message.Commodities {
		var com Commodity
		var sc StationCommodity
		db.Where(Commodity{Name: a.Name}).Attrs(Commodity{MeanPrice: a.MeanPrice}).FirstOrCreate(&com)
		if com.MeanPrice != a.MeanPrice {
			com.MeanPrice = a.MeanPrice
			db.Save(&com)
		}
		db.Where(StationCommodity{CommodityID: com.ID, StationID: station.ID}).FirstOrInit(&sc)

		sc.BuyPrice = a.BuyPrice
		sc.Stock = a.Stock
		sc.SellPrice = a.SellPrice
		sc.Demand = a.Demand
		sc.StockBracket = a.StockBracket
		sc.DemandBracket = a.DemandBracket

		db.Save(&sc)
	}

	log.Printf("Updated %d in %s/%s\n", len(c.Message.Commodities), system.Name, station.Name)

}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintln(w, "<h1>Hello World</h1>")
	return nil
	// return detailTmpl.Execute(w, r, book)
}
