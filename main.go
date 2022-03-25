package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// album represents data about a record album.
//you have to name your fields "Uppercase" so that they will be exportable *
//add the json:"<tags>" in lowercase
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

//you have to name your fields "Uppercase" so that they will be exportable *
//add the json:"<tags>" in lowercase
type Product struct {
	Product_id        int     `json:"product_id"`
	Name              string  `json:"name"`
	Quantity_in_stock int     `json:"quantity_in_stock"`
	Unit_price        float32 `json:"unit_price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

//{"ID": "4", "Title": "Pops", "Artist": "Lisa Stansfield - Mulligan", "Price": 20.98}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.GET("/products/", getProducts)

	router.Run("localhost:8090")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

var arrProducts = []Product{}

func getProducts(c *gin.Context) {

	//declare connection
	db, err := sql.Open("mysql", "root:Mag17615@@tcp(127.0.0.1:3306)/sql_store")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//queyr the database
	results, err := db.Query("SELECT product_id,name,quantity_in_stock,unit_price FROM products")
	if err != nil {
		panic(err.Error())
	}

	var _product Product
	//loop through the resultset
	for results.Next() {
		var product Product
		err = results.Scan(&product.Product_id, &product.Name, &product.Quantity_in_stock, &product.Unit_price)

		_product.Product_id = product.Product_id
		_product.Name = product.Name
		_product.Quantity_in_stock = product.Quantity_in_stock
		_product.Unit_price = product.Unit_price
		addProduct(_product) //append each product object to the array

		if err != nil {
			panic(err.Error())
		}

	}
	c.IndentedJSON(http.StatusOK, arrProducts)
}

//the syntax below is how you append objects to a struct
func addProduct(_product Product) {
	var my_product = new(Product)
	my_product.Product_id = _product.Product_id
	my_product.Name = _product.Name
	my_product.Quantity_in_stock = _product.Quantity_in_stock
	my_product.Unit_price = _product.Unit_price
	arrProducts = append(arrProducts, *my_product)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
