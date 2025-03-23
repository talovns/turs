package tur

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=12345 dbname=postgres port=5432 sslmode=disable search_path=tur"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Миграция схемы с ловкой ошибки
	if err := db.AutoMigrate(&Oplata{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

type Oplata struct {
	KOD_OPLATI  uint   `gorm:"primaryKey" json:"kod_oplati"`
	DATA_OPLATI string `gorm:"column:data_oplati" json:"data_oplati"`
	SUMM        int    `gorm:"column:summ" json:"summ"`
	KOD_PUTIVKI int    `gorm:"column:kod_putivki" json:"kod_putivki"`
}

func getOrders(c *gin.Context) {
	var order []Oplata
	if err := db.Table("tur.oplata").Find(&order).Error; err != nil {
		log.Printf("Error fetching orders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (Oplata) TableName() string {
	return "tur.oplata"
}

func getOrderByID(c *gin.Context) {
	id := c.Param("id")
	var order Oplata
	if err := db.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}
func createOrder(c *gin.Context) {
	var neworder Oplata
	if err := c.BindJSON(&neworder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	db.Create(&neworder)
	c.JSON(http.StatusCreated, neworder)
}

func updateOrder(c *gin.Context) {
	id := c.Param("id")
	var updateOrder Oplata
	if err := c.BindJSON(&updateOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	if err := db.Model(&Oplata{}).Where("id = ?", id).Updates(updateOrder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}
	c.JSON(http.StatusOK, updateOrder)
}

func deleteOrder(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Oplata{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "order not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}
