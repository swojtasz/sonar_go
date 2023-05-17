package main

import (
    "github.com/labstack/echo/v4"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Product struct {
    gorm.Model
    Name        string
    Description string
    Price       float32
}

func setCORSHeader(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set("Access-Control-Allow-Origin", "*")
        return next(c)
    }
}

func migrate(db *gorm.DB) {
    db.AutoMigrate(&Product{})
}

type ProductController struct {
    db *gorm.DB
}

func (pc *ProductController) GetProducts(c echo.Context) error {
    var products []Product
    pc.db.Find(&products)
    return c.JSON(200, products)
}

func (pc *ProductController) GetProduct(c echo.Context) error {
    id := c.Param("id")
    var product Product
    pc.db.First(&product, id)
    return c.JSON(200, product)
}

func (pc *ProductController) CreateProduct(c echo.Context) error {
    var product Product
    c.Bind(&product)
    pc.db.Create(&product)
    return c.JSON(201, product)
}

func (pc *ProductController) UpdateProduct(c echo.Context) error {
    id := c.Param("id")
    var product Product
    pc.db.First(&product, id)
    c.Bind(&product)
    pc.db.Save(&product)
    return c.JSON(200, product)
}

func (pc *ProductController) DeleteProduct(c echo.Context) error {
    id := c.Param("id")
    var product Product
    pc.db.Delete(&product, id)
    return c.NoContent(204)
}

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    migrate(db)
    pc := &ProductController{db: db}
    e := echo.New()
	e.Use(setCORSHeader)
    e.GET("/products", pc.GetProducts)
    e.GET("/products/:id", pc.GetProduct)
    e.POST("/products", pc.CreateProduct)
    e.PUT("/products/:id", pc.UpdateProduct)
    e.DELETE("/products/:id", pc.DeleteProduct)
    e.Logger.Fatal(e.Start(":3015"))
}
