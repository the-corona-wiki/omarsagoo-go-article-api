package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type keywords struct {
	gorm.Model

	Name  string `json:"name"`
	Value string `json:"value"`
	Rank  int    `json:"rank"`
	Major string `json:"major"`
}

type articleJSON struct {
	gorm.Model

	Abstract  string     `json:"abstract"`
	URL       string     `json:"url"`
	Title     string     `json:"title"`
	Byline    string     `json:"byline"`
	Published string     `json:"published"`
	Keywords  []keywords `json:"keywords"`
}

type listOfArticle struct {
	Articles []articleJSON `json:"list"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func allArticlesDB(collections []string, db *mongo.Database) map[string]interface{} {
	// var allArticles map[string]interface{}
	allArticles := make(map[string]interface{})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{})
	var article []bson.M

	for _, collection := range collections {
		coll := db.Collection(collection)
		cursor, err := coll.Find(ctx, bson.D{}, opts)
		check(err)
		err = cursor.All(ctx, &article)
		check(err)

		allArticles[collection] = article
	}
	return allArticles
}

func giveJSON(json []byte) {

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb+srv://omar:Corona1234@cluster0-3rny3.mongodb.net/test?retryWrites=true&w=majority",
	))
	check(err)

	db := client.Database("articles")

	results, err := db.ListCollectionNames(ctx, bson.D{})
	check(err)

	articles := allArticlesDB(results, db)

	json, err := json.Marshal(articles)
	check(err)

	e := echo.New()
	e.GET("/all", func(c echo.Context) error {
		return c.String(http.StatusOK, string(json))
	})
	e.Logger.Fatal(e.Start(":1323"))
	// fmt.Println(string(json))
}
