package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"fmt"
	"net/smtp"

	"github.com/bxcodec/faker/v4"
	"github.com/gofiber/fiber/v2"
)

func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.GetFullname()
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}

type OrderRequest struct {
	Code      string           `json:"code"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Email     string           `json:"email"`
	Address   string           `json:"address"`
	Country   string           `json:"country"`
	City      string           `json:"city"`
	Zip       string           `json:"zip"`
	Products  []map[string]int `json:"products"`
}

func CreateOrder(c *fiber.Ctx) error {
	var request OrderRequest
	var link models.Link
	var order models.Order

	if err := c.BodyParser(&request); err != nil {
		fmt.Println("WTF ?")
		panic(err)
	}

	tx := database.DB.Begin()

	database.DB.Preload("User").Where("code=?", request.Code).Find(&link)

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid code",
		})
	}

	order = models.Order{
		TransactionId:   faker.CCNumber(),
		Code:            request.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	if e := tx.Create(&order).Error; e != nil {
		tx.Rollback()
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error creating order",
		})
	}

	for _, requestedProduct := range request.Products {
		var product models.Product
		database.DB.Where("id=?", requestedProduct["product_id"]).Find(&product)
		if product.Id == 0 {
			tx.Rollback()
			panic("Which product ?")
		}
		total := product.Price * float64(requestedProduct["quantity"])
		orderItem := models.OrderItem{}
		orderItem.OrderId = order.Id
		orderItem.ProductTitle = product.Title
		orderItem.Price = total
		orderItem.Quantity = uint(requestedProduct["quantity"])
		orderItem.AdminRevenue = 0.9 * total
		orderItem.AmbassadorRevenue = 0.1 * total

		if errr := tx.Create(&orderItem).Error; errr != nil {
			tx.Rollback()
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "error creating order item",
			})
		}
	}
	tx.Commit()
	return c.JSON(order)
}

type OrderComplete struct {
	TransactionId string `json:"transaction_id"`
}

func CompleteOrder(c *fiber.Ctx) error {
	var data OrderComplete
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("Send proper body")
		panic(err)
	}

	var order models.Order
	database.DB.Preload("OrderItems").Where("transaction_id=?", data.TransactionId).Find(&order)

	if order.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Cheating karta hai tu !!",
		})
	}

	order.Complete = true

	database.DB.Save(&order)

	go func(o models.Order) {
		ambassadorRevenue := 0.0
		adminRevenue := 0.0

		for _, orderItem := range o.OrderItems {
			ambassadorRevenue += orderItem.AmbassadorRevenue
			adminRevenue += orderItem.AdminRevenue
		}

		var user models.User

		database.DB.Where("id=?", o.UserId).Find(&user)

		database.Cache.ZIncrBy(context.Background(), "rankings", ambassadorRevenue, user.GetFullname())

		username := "131057e5ec981d"
		password := "6991da74d1ad96"
		host := "smtp.mailtrap.io"
		addr := "smtp.mailtrap.io:2525"
		auth := smtp.PlainAuth("", username, password, host)

		ambassadorMessage := []byte(fmt.Sprintf("You have earned $%f from %s", ambassadorRevenue, order.Code))

		ambassadordMailError := smtp.SendMail(addr, auth, "noreply@ambassadorgo.com", []string{o.AmbassadorEmail}, ambassadorMessage)

		if ambassadordMailError != nil {
			fmt.Println("Error sending mail to ambassador")
			panic(ambassadordMailError)
		} else {
			fmt.Println("Mail sent successfully to ambassador")
		}

		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed", o.Id, adminRevenue))

		adminEmailError := smtp.SendMail(addr, auth, "noreply@ambassadorgo.com", []string{"admin@ambassadorgo.com"}, adminMessage)

		if adminEmailError != nil {
			fmt.Println("Erro sending mail to admin")
			panic(adminEmailError)
		} else {
			fmt.Println("Mail sent successfully to admin")
		}
	}(order)

	return c.JSON(fiber.Map{
		"message": "Order Completed",
	})
}
