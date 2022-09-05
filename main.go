package main

import (
	"github.com/gotneb/manga_api/db"
	"github.com/gotneb/manga_api/web"
)

func main() {
	section := "abcdefghijklmnopqrstuvwxyz"
	for _, v := range section {
		links := web.FetchPages(string(v))

		for _, link := range links {
			manga, err := web.FetchMangaData(link)
			if err != nil {
				panic(err)
			}
			db.AddManga(&manga)
		}
	}
	/*
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}

		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
		}
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()

		coll := client.Database("manga_api").Collection("meus_mangas")
		_, err = coll.InsertOne(
			context.TODO(),
			bson.D{
				{"title", "akame ga kill"},
				{"total_chapters", 40},
			},
		)

		if err != nil {
			panic(err)
		}
		fmt.Println("Sucess!")
	*/
	/*
		===========================================
		title := "Berserk"

		var result bson.M
		err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)
		if err == mongo.ErrNoDocuments {
			fmt.Printf("No document was found with the title %s\n", title)
			return
		}
		if err != nil {
			panic(err)
		}

		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", jsonData)
	*/
}
