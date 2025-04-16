package controllers

import (
	"context"
	"golang-restraunt-management/database"
	"golang-restraunt-management/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
) 


type InvoiceViewFormat struct {
	Invoice_id           string 
	Payment_method       string
	Order_id             string
	Payment_status       *string
	Payment_due          interface{}
	Table_number         interface{}
	Payment_due_date     time.Time
	Order_details        interface{}
}

var InvoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoices")

func GetInvoices() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result,err := InvoiceCollection.Find(context.TODO(), bson.M{})

		defer cancel()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch invoices"})
			return
		}
		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
			c.JSON(500, gin.H{"error": "Failed to decode invoices"})
			return
		}
		c.JSON(200, allInvoices)
		defer cancel()
		
	}
}

func GetInvoice() gin.HandlerFunc{	
	return func(c *gin.Context){
		ctx,cancel = context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoice_id")

		var invoice models.Invoice
		if invoiceId == "" {
			c.JSON(400, gin.H{"error": "invoice_id is required"})
			return
		}
		err := InvoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch invoice"})
			return
		}
		var invoiceView InvoiceViewFormat

		allOrderItems,err := ItemsByOrder(invoice.Order_id)

		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date
		invoiceView.Payment_method = ""



		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}
		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]
		c.JSON(200, invoiceView)
		defer cancel()

	}
}

func CreateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}

func UpdateInvoice() gin.HandlerFunc{
	return func(c *gin.Context){
		// Write your code here
	}
}