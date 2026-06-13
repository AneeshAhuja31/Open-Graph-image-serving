package main

import (
	"encoding/json"
	
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func main() {

	InitDB()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{username}/",func(w http.ResponseWriter, r *http.Request) {
		username:= r.PathValue("username")
		collection := Client.Database("og_db").Collection("users")
		ctx := r.Context()
		var user User

		err := collection.FindOne(
			ctx,
			bson.M{"username":username},
		).Decode(&user)
		w.Header().Set("Content-Type","application/json")
		if err != nil{
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"User does not exist"}`))
			return
		}
		json.NewEncoder(w).Encode(user)
	})

	mux.HandleFunc("GET /")
	http.ListenAndServe(":8080",mux)
}