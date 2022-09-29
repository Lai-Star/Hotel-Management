package controllers

import (
	"context"
	"go-hotel/database"
	"go-hotel/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.Opencollection(database.Client, "orderItem")

//get all the order items func
func GetOrderItems(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	result, err := orderItemCollection.Find(context.TODO(), bson.M{})
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered items"})
		return
	}
	var allOrderItems []bson.M
	if err = result.All(ctx, &allOrderItems); err != nil {
		log.Fatal(err)
		return
	}
	c.JSON(http.StatusOK, allOrderItems)

}

//get the single order item by ID
func GetOrderItemsByOrder(c *gin.Context) {
	orderId := c.Param("order_id")
	allorderItems, err := ItemByOrder(orderId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items by order ID"})
		return
	}
	c.JSON(http.StatusOK, allorderItems)

}

func ItemByOrder(id string) (orderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchstage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookupstage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindstage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookuporderstage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindorderstage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookuptablestage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "order.table_id"}, {"as", "table"}}}}
	unwindtablestage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectstage := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"amount", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity", 1},
		}}}

	groupstage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$push", "$$ROOT"}}}}}}

	projectstagetwo := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"total_number", "$_id.table_number"},
			{"order_items", 1},
		}}}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchstage,
		lookupstage,
		unwindstage,
		lookuporderstage,
		unwindorderstage,
		lookuptablestage,
		unwindtablestage,
		projectstage,
		groupstage,
		projectstagetwo})

	if err != nil {
		panic(err)
	}
	if err = result.All(ctx, &orderItems); err != nil {
		panic(err)
	}
	defer cancel()

	return orderItems, err
}

//get order item
func GetOrderItem(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var orderitem models.OrderItem
	orderitemId := c.Param("order_item_id")
	err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": orderitemId}).Decode(&orderitem)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered item"})
		return
	}
	c.JSON(http.StatusOK, orderitem)

}

//update orderitem based on ID
func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var orderitem models.OrderItem
		orderitemId := c.Param("order_item_id")
		filter := bson.M{"order_item_id": orderitemId}
		var updateObj primitive.D

		if orderitem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price", *&orderitem.Unit_price})
		}
		if orderitem.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", *orderitem.Quantity})
		}
		if orderitem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", *orderitem.Food_id})
		}
		orderitem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderitem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			msg := "Order item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

//create order item
func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemtobeInserted := []interface{}{}

		order.Table_id = orderItemPack.Table_id

		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id
			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemtobeInserted = append(orderItemtobeInserted, orderItem)

		}

		insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemtobeInserted)
		if err != nil {
			log.Fatal(err)
		}
		defer cancel()

		c.JSON(http.StatusOK, insertedOrderItems)

	}

}
