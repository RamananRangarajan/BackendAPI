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

func CreatePost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	user := new(User)
	post := new(Post)
	_ = json.NewDecoder(request.Body).Decode(&post)
	post.PostedTimeStamp = time.Now().Format(time.RFC850)
	h := sha1.New()
	h.Write([]byte(post.Password))
	var HashedPassword string = base64.URLEncoding.EncodeToString(h.Sum(nil))
	collection_user := client.Database("appointy_instAPI").Collection("userinfo")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection_user.FindOne(ctx, User{Name: post.UserName, PasswordHash: HashedPassword}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if post.UserName == user.Name && HashedPassword == user.PasswordHash {
		collection_post := client.Database("appointy_instAPI").Collection("gallery")
		result, _ := collection_post.InsertOne(ctx, post)
		json.NewEncoder(response).Encode(result)
	} else {
		fmt.Println("Invalid UserName or Password")
		return
	}
}

func GetPost(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Post
	collection := client.Database("appointy_instAPI").Collection("gallery")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

func GetGallery(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var gallery []Post
	collection := client.Database("appointy_instAPI").Collection("gallery")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		gallery = append(gallery, post)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(gallery)
}
