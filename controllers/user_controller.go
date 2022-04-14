package controllers

import (
	"echo-mongo-api/configs"
	"echo-mongo-api/models"
	"echo-mongo-api/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

var stateCollection *mongo.Collection = configs.GetCollection(configs.DB, "state")

// var validate = validator.New()

func MakeRequest(URL string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Header_Key", "Header_Value")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Err is", err)
	}
	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	response := string(resBody)

	return response
}

func InsertTable(c echo.Context) error {
	url := "https://data.covid19india.org/v4/min/data.min.json"

	res := MakeRequest(url)                // Making the request
	resBytes := []byte(res)                // Converting the string "res" into byte array
	var jsonRes map[string]interface{}     // declaring a map for key names as string and values as interface{}
	_ = json.Unmarshal(resBytes, &jsonRes) // Unmarshalling

	state := c.Param("stateId")
	details_map := jsonRes[state].(map[string]interface{}) // type the interface again to a map with key string type and value as interface
	total := details_map["total"].(map[string]interface{})
	confirmed := total["confirmed"].(float64)
	deceased := total["deceased"].(float64)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newState := models.State{
		Id:        primitive.NewObjectID(),
		State:     state,
		Confirmed: confirmed,
		Deceased:  deceased,
	}
	result, err := stateCollection.InsertOne(ctx, newState)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	return c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &echo.Map{"data": result}})
}

// func UpdateState(c echo.Context) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	stateId := c.Param("stateId")
// 	var state models.State
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(stateId)
// 	err := stateCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&state)

// 	update := bson.M{"state": state.State, "confirmed": state.Confirmed, "deceased": state.Deceased}

// 	result, err := stateCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
// 	}

// 	//get updated user details
// 	var updatedState models.State
// 	if result.MatchedCount == 1 {
// 		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedState)
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
// 		}
// 	}

// 	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": updatedState}})

// }

func GetAState(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	stateId := c.Param("stateId")
	var state models.State
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(stateId)

	err := stateCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}
	f := c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": state}})
	return f
}
