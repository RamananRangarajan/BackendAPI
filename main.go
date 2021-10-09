package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Password     string             `json:"password,omitempty" bson:"-"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	PasswordHash string             `json:"passwordhash,omitempty" bson:"passwordhash,omitempty"`
}

type Post struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName        string             `json:"username,omitempty" bson:"username,omitempty"`
	Password        string             `json:"password,omitempty" bson:"-"`
	Caption         string             `json:"caption" bson:"caption"`
	ImageURL        string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	PostedTimeStamp string             `json:"postedtimestamp,omitempty" bson:"postedtimestamp,omitempty"`
}

func GetPostList(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	var gallery []Post
	collection_user := client.Database("appointy").Collection("community")
	ctx_user, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection_user.FindOne(ctx_user, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	if id == user.ID {
		collection_post := client.Database("appointy").Collection("gallery")
		ctx_post, _ := context.WithTimeout(context.Background(), 30*time.Second)
		cursor, err := collection_post.Find(ctx_post, Post{UserName: user.Name})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		defer cursor.Close(ctx_post)
		for cursor.Next(ctx_post) {
			fmt.Println("For loop entered")
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
	} else {
		fmt.Println("User Not Found")
		return
	}
}

func main() {
	fmt.Println("Application Started")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/posts", CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPost).Methods("GET")
	router.HandleFunc("/posts/users/{id}", GetPostList).Methods("GET")
	router.HandleFunc("/community", GetCommunity).Methods("GET")
	router.HandleFunc("/gallery", GetGallery).Methods("GET")
	http.ListenAndServe(":3000", router)
}
