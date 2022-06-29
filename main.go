package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"test-24h/config/cache"
	"test-24h/config/db"
	"test-24h/utils"
)

type DuLieu1 struct {
	A string
	B int
	C string
}

type DuLieu2 struct {
	D int
	E string
	G int
}

type DuLieu3 struct {
	A string
	E string
	G int
}

var wg sync.WaitGroup

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/test-24h", myHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	wg.Add(3)
	c1 := make(chan string)
	c2 := make(chan string)

	res := make(map[string][]interface{})
	res["json"] = []interface{}{}

	db := db.Connect()
	defer db.Close()
	var dl1s []DuLieu1
	go func() {
		dl1s = getDataFromMySQL(db)
		res["json"] = append(res["json"], dl1s)
		c1 <- "Done"
		wg.Done()
	}()

	cache := cache.Connect()
	defer cache.Close()
	addDataRedis(cache)
	var dl2s []DuLieu2
	go func() {
		dl2s = getDataRedis(cache, []string{"id1234"})
		select {
		case <-c1:
			res["json"] = append(res["json"], dl2s)
			c2 <- "Done"
		}
		wg.Done()
	}()

	var dl3s []DuLieu3
	go func() {
		dl3s = getDataFromFile()
		select {
		case <-c2:
			res["json"] = append(res["json"], dl3s)
		}
		wg.Done()
	}()

	wg.Wait()
	utils.RespondSuccess(w, 200, res)
}

func getDataFromFile() []DuLieu3 {
	csvFile, err := os.Open("data.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	dl3s := []DuLieu3{}

	for _, line := range csvLines {
		tmpG, _ := strconv.Atoi(line[2])
		dl3 := DuLieu3{A: line[0], E: line[1], G: tmpG}
		dl3s = append(dl3s, dl3)
	}

	return dl3s
}

func getDataRedis(cache *redis.Client, keys []string) []DuLieu2 {
	var result []DuLieu2
	var res DuLieu2
	for _, key := range keys {
		val, err := cache.Get(key).Result()
		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal([]byte(val), &res)
		result = append(result, res)
	}

	return result
}

func addDataRedis(re *redis.Client) {
	json, err := json.Marshal(DuLieu2{D: 1, E: "E", G: 2})
	if err != nil {
		fmt.Println(err)
	}

	err = re.Set("id1234", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func getDataFromMySQL(db *sql.DB) []DuLieu1 {
	rs, err := db.Query("SELECT A, B, C FROM DuLieu1")
	if err != nil {
		panic(err.Error())
	}

	var mySqlData []DuLieu1
	for rs.Next() {
		var dl1 DuLieu1
		err = rs.Scan(&dl1.A, &dl1.B, &dl1.C)
		if err != nil {
			panic(err.Error())
		}
		mySqlData = append(mySqlData, dl1)
	}

	return mySqlData
}
