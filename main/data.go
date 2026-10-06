package main

import "github.com/weaviate/weaviate/entities/models"

var MoviesDataObjects = []map[string]interface{}{
	{
		"title":       "The Matrix",
		"description": "A computer hacker learns about the true nature of reality and his role in the war against its controllers.",
		"genre":       "Science Fiction",
	},
	{
		"title":       "Spirited Away",
		"description": "A young girl becomes trapped in a mysterious world of spirits and must find a way to save her parents and return home.",
		"genre":       "Animation",
	},
	{
		"title":       "The Lord of the Rings: The Fellowship of the Ring",
		"description": "A meek Hobbit and his companions set out on a perilous journey to destroy a powerful ring and save Middle-earth.",
		"genre":       "Fantasy",
	},
}

type Movie struct {
	Title string `json: "string"`
	Description string `json: "string"`
	Genre string `json: "string"`
}

// using bundled vectorizer to calculate the vectors
var MovieClassName = "Movie"

var MoviesClassObject = &models.Class{
	Class: MovieClassName,
	// Vectorizer: "text2vec-weaviate", // requires API Key
	Vectorizer:  "text2vec-transformers", // using local vectorizer
	Description: "Collection of movies with their descriptions and genre information",
	ModuleConfig: map[string]interface{}{
		"generative-anthropic": map[string]interface{}{
			"model": "claude-haiku-4-5",
		},
	},
}
