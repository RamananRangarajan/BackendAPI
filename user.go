package main

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	user := new(User)
	_ = json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("appointy_instAPI").Collection("userinfo")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	h := sha1.New()
	h.Write([]byte(user.Password))
	user.PasswordHash = base64.URLEncoding.EncodeToString(h.Sum(nil))
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func GetUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("appointy_instAPI").Collection("userinfo")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if id == user.ID {
		fmt.Println("User Found")
		json.NewEncoder(response).Encode(user)
	} else {
		fmt.Println("User Not Found")
		return
	}
}

func userPostinfo(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var userinfo []User
	collection := client.Database("appointy_instAPI").Collection("userinfo")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		userinfo = append(userinfo, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(userinfo)
}
