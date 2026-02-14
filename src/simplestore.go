package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// Matching the API spec for hw5
type product struct {
	ProductId    int    `json:"product_id"` // capitalized for JSON marshalling
	Sku          string `json:"sku"`
	Manufacturer string `json:"manufacturer"`
	CategoryId   int    `json:"category_id"`
	Weight       int    `json:"weight"`
	SomeOtherId  int    `json:"some_other_id"`
}

// In-memory product store
// Anticipating concurrent access with mostly reads and occasional writes, so using RWMutex for better read performance
// Sync.Map could be an alternative for a more concurrent read-heavy workload, but RWMutex is simpler for this use case
type productStore struct {
	mu       sync.RWMutex
	products map[int]product
}

var store productStore

func (s *productStore) Update(key int, value product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.products[key] = value
}

func (s *productStore) Get(key int) (product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.products[key]
	return val, exists
}

func (s *productStore) Exists(key int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.products[key]
	return exists
}

func (s *productStore) Delete(key int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.products[key]
	if exists {
		delete(s.products, key)
	}
	return exists
}

func (s *productStore) seed() {
	s.Update(1, product{ProductId: 1, Sku: "sku1", Manufacturer: "manufacturer1", CategoryId: 1, Weight: 10, SomeOtherId: 100})
	s.Update(200, product{ProductId: 200, Sku: "sku2", Manufacturer: "manufacturer2", CategoryId: 2, Weight: 20, SomeOtherId: 200})
	s.Update(3000, product{ProductId: 3000, Sku: "sku3", Manufacturer: "manufacturer3", CategoryId: 3, Weight: 30, SomeOtherId: 300})
}

// Function calls from HTTP requests

// Retrieve product details from the store and return as JSON
func getProductById(c *gin.Context) {
	productId := c.Param("productId")
	// convert productId to int
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid product ID"}) // if productId is invalid, return 404 with error message (per spec, not 400)
		return
	}
	// if product is found, return 200 with product details as JSON
	product, exists := store.Get(productIdInt)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"}) // if product is not found, return 404 with error message
		return
	} else {
		c.IndentedJSON(http.StatusOK, product) // if product is found, return 200 with product details as JSON
	}

	// if an issue with the server, return 500 with error message

}

// Parse request body for product details and update the store
func updateProductDetails(c *gin.Context) {
	productId := c.Param("productId")
	// convert productId to int
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"}) // if productId is invalid, return 400 with error message
		return
	}
	// get product details from request body
	var updatedProduct product
	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product details"}) // if product details are invalid, return 400 with error message
		return
	}
	exists := store.Exists(productIdInt)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"}) // if product is not found, return 404 with error message
		return
	} else {
		store.Update(productIdInt, updatedProduct) // update product with details from request body
		c.Status(http.StatusNoContent)             // if product is found and updated return 204
	}

	// if an issue with the server, return 500 with error message
}

func main() {
	//initialize product store
	store = productStore{
		products: make(map[int]product),
	}
	store.seed() // seed with some initial data

	//set up HTTP server and routes

	router := gin.Default()

	router.GET("/products/:productId", func(c *gin.Context) {
		getProductById(c)
	})
	router.POST("/products/:productId/details", func(c *gin.Context) {
		updateProductDetails(c)
	})

	router.Run(":8080")
}
