package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"vectors-lab/main/server"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

var client *weaviate.Client

func init() {
	cfg := weaviate.Config{
		Host:   "localhost:8080",
		Scheme: "http",
	}
	newClient, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	client = newClient
}

func main() {
	// CreateMovieSchema()
	// GetSchema()
	// TestInsertMoviesObjects()

	// InsertMoviesJsonFile()
	// SemanticSearch("sci-fi")

	vs := server.NewVectorServer(8809)
	vs.PrepareAndStart("/api/vectors", func(w http.ResponseWriter, r *http.Request) {
		server.EnableCors(w)
		// todo: narrow down to POST requests
		ssReq := server.SemanticSearchRequest{}
		if r.Body != nil {
			defer r.Body.Close()
		}
		extB, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(extB, &ssReq)

		log.Printf("Request received %v", ssReq)
		results := SemanticSearch(ssReq.Text)
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "failed to encode search results", http.StatusInternalServerError)
		}
	})
}

// --------------- DB related

func GetSchema() {
	schema, err := client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v \n", schema)
}

func CreateMovieSchema() {
	err := client.Schema().ClassCreator().WithClass(MoviesClassObject).Do(context.Background())
	if err != nil {
		panic(err)
	}
}

// deprecated
func TestInsertMoviesObjects() {
	// Insert objects
	objects := make([]*models.Object, len(MoviesDataObjects))
	for i, obj := range MoviesDataObjects {
		objects[i] = &models.Object{
			Class:      "Movie",
			Properties: obj,
		}
	}

	_, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Imported & vectorized %d objects into the Movie collection\n", len(MoviesDataObjects))
}

func InsertMoviesObjects(movies []Movie) {
	// Insert objects
	objects := make([]*models.Object, len(movies))
	for i, obj := range movies {
		objects[i] = &models.Object{
			Class:      "Movie",
			Properties: obj,
		}
	}

	_, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Imported & vectorized %d objects into the Movie collection\n", len(movies))
}

func InsertMoviesJsonFile() {
	fB, err := os.ReadFile("movies.json")
	if err != nil {
		panic(err)
	}

	toInsert := []Movie{}
	err = json.Unmarshal(fB, &toInsert)
	if err != nil {
		panic(err)
	}

	for _, o := range toInsert {
		fmt.Printf("inserting %s \n", o.Title)
	}
	fmt.Println()
	InsertMoviesObjects(toInsert)
}

func SemanticSearch(text2Search string) []interface{} {
	title := graphql.Field{Name: "title"}
	description := graphql.Field{Name: "description"}
	genre := graphql.Field{Name: "genre"}

	// https: //docs.weaviate.io/weaviate/api/graphql/additional-properties
	additional := graphql.Field{Name: "_additional", Fields: []graphql.Field{
		{Name: "id"},
		{Name: "distance"},
		{Name: "vector"},
	}}

	nearText := client.GraphQL().NearTextArgBuilder().
		WithConcepts([]string{text2Search})

	result, err := client.GraphQL().Get().
		WithClassName("Movie").
		WithNearText(nearText).
		WithLimit(2).
		WithFields(title, description, genre, additional).
		Do(context.Background())

	if err != nil {
		panic(err)
	}

	// Inspect the results
	if result.Errors != nil {
		errorJSON, marshalErr := json.MarshalIndent(result.Errors, "", "  ")
		if marshalErr != nil {
			fmt.Printf("GraphQL errors (fallback): %#v\n", result.Errors)
			panic("GraphQL error ocurred")
		}
		fmt.Printf("GraphQL errors:\n%s\n", string(errorJSON))
		panic("GraphQL error ocurred")
	}

	data := result.Data["Get"].(map[string]interface{})
	movies := data["Movie"].([]interface{})

	for _, movie := range movies {
		jsonData, err := json.MarshalIndent(movie, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonData))
	}

	return movies
}
