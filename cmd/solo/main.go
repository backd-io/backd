package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"log"

	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/structs"
	"github.com/globalsign/mgo"
)

func main() {
	var (
		mongo    *db.Mongo
		names    []string
		err      error
		hostname string
	)

	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	println(hostname)

	mongo, err = db.NewMongo("mongodb://localhost:27017")
	er(err)

	err = mongo.Insert("ttt", "ttt", map[string]interface{}{
		"this":  "is",
		"some":  3,
		"tests": true,
	})
	er(err)

	now := time.Now()
	names, err = mongo.Session().DB("ttt").CollectionNames()
	fmt.Println(names)
	er(err)
	fmt.Println("fmt.CollectionNames:time:", time.Since(now).String())

	appTableName := db.NewXID().String()

	now = time.Now()
	var colInfo mgo.CollectionInfo
	colInfo.Validator = structs.ApplicationValidator()
	err = mongo.Session().DB("ttt").C(appTableName).Create(&colInfo)
	er(err)
	fmt.Println("fmt.CreateCol:time:", time.Since(now).String())

	var (
		app structs.Application
		d   []byte
	)

	now = time.Now()
	app.ID = db.NewXID().String()
	app.Name = "This is a test"
	app.Description = "some random description"
	app.SetCreate(db.NewXID().String(), db.NewXID().String())

	d, err = json.Marshal(app)
	fmt.Println("app:", string(d))

	err = mongo.Insert("ttt", appTableName, &app)
	er(err)
	fmt.Println("fmt.Insert:time:", time.Since(now).String())

	userTable := db.NewXID().String()

	now = time.Now()
	var colInfo2 mgo.CollectionInfo
	colInfo2.Validator = structs.UserValidator()
	err = mongo.Session().DB("ttt").C(userTable).Create(&colInfo2)
	er(err)
	fmt.Println("fmt.CreateCol:time:", time.Since(now).String())

	now = time.Now()
	var user structs.User
	user.Name = "Antonio"
	user.Email = "antoniofernadezvara+t@gmail.com"
	user.Username = user.Email
	user.ID = db.NewXID().String()
	user.SetPassword("kkdelavaca$$11")
	user.SetCreate(db.NewXID().String(), db.NewXID().String())

	d, err = json.Marshal(user)
	fmt.Println("user:", string(d))

	err = mongo.Insert("ttt", userTable, &user)
	er(err)
	fmt.Println("fmt.Insert:time:", time.Since(now).String())

	users := make(map[int]structs.User)

	now = time.Now()
	now2 := time.Now()
	for a := 0; a < 10; a++ {

		var user structs.User
		user.Name = fmt.Sprintf("Antonio %d", a)
		user.Email = fmt.Sprintf("antoniofernadezvara+%d@gmail.com", a)
		user.Username = user.Email
		user.ID = db.NewXID().String()
		//user.SetPassword("kkdelavaca$$11")
		user.SetCreate(db.NewXID().String(), db.NewXID().String())

		users[a] = user

		if a < 10 {
			fmt.Println("added", a, "items", time.Since(now2).String())
		}

		if a%1000 == 0 {
			fmt.Println("added", a, "items", time.Since(now2).String())
			now2 = time.Now()
		}

	}

	var u structs.User
	for b := 0; b < 10; b++ {
		u.Name = fmt.Sprintf("A%d", b)
		u.Email = fmt.Sprintf("antonio+%d@gmail.com", b)
		u.Username = user.Email
		u.ID = db.NewXID().String()
		u.SetCreate(db.NewXID().String(), db.NewXID().String())

		_ = mongo.Insert("ttt", "fastTest", &u)

	}

	now3 := time.Now()
	var ids []string
	err = mongo.Session().DB("ttt").C("fastTest").Find(map[string]interface{}{}).Distinct("n", &ids)

	fmt.Println("err:", err)
	fmt.Println("ids:", ids)

	fmt.Println("since:", time.Since(now3).String())

	var users2 []structs.User

	query := map[string]interface{}{
		"n": map[string]interface{}{
			"$in": ids,
		},
	}

	err = mongo.GetMany("ttt", "fastTest", query, []string{}, &users2, 0, 0)
	fmt.Println("err:GetMany", err)
	fmt.Println("users2:GetMany", users2)

	var user2 structs.User
	err = mongo.GetOneByID("ttt", "ttt", "nonexistent", &user2)

	fmt.Println("err:getOneByID:", err)

	var user3 structs.User
	err = mongo.GetOne("ttt", "ttt", map[string]interface{}{"_id": "nonexistent"}, &user3)

	fmt.Println("err:getOne:", err)

	fmt.Println("fmt.Create1M:time:", time.Since(now).String())

	fmt.Println("Sleeping for 20 seconds....")
	time.Sleep(20 * time.Second)

	// var colInfo mgo.CollectionInfo
	// colInfo.ValidationLevel = "strict"
	// colInfo.Validator = map[string]interface{}{
	// 	"$jsonSchema": map[string]interface{}{
	// 		"bsonType": "object",
	// 		"required": []string{"this", "some", "tests"},
	// 		"properties": map[string]interface{}{
	// 			"this": map[string]interface{}{
	// 				"bsonType":  "string",
	// 				"maxLength": 20,
	// 				"pattern":   "^[a-zA-Z0-9_.+-]+$",
	// 			},
	// 			"some": map[string]interface{}{
	// 				"bsonType": "int",
	// 			},
	// 			"password": map[string]interface{}{
	// 				"bsonType": "bool",
	// 			},
	// 		},
	// 	},
	// }

	// err = mongo.Session().DB("ttt").C("ttt4").Create(&colInfo)
	// er(err)

	// err = mongo.Insert("ttt", "ttt1", map[string]interface{}{
	// 	"this":  "is",
	// 	"some":  3,
	// 	"tests": true,
	// })
	// fmt.Println("good", err)
	// er(err)

	// // err = mongo.Insert("ttt", "ttt1", map[string]interface{}{
	// // 	"this":  "is",
	// // 	"some":  "3",
	// // 	"tests": true,
	// // })
	// // fmt.Println("bad", err)
	// //er(err)

	// var colInfo2 mgo.CollectionInfo
	// colInfo2.ValidationLevel = "strict"
	// colInfo2.Validator = map[string]interface{}{
	// 	"$jsonSchema": map[string]interface{}{
	// 		"bsonType":             "object",
	// 		"additionalProperties": true,
	// 	},
	// }

	// err = mongo.Session().DB("ttt").C("ttt5").Create(&colInfo2)
	// er(err)
	// err = mongo.Insert("ttt", "ttt2", map[string]interface{}{
	// 	"this":  "is",
	// 	"some":  3,
	// 	"tests": true,
	// })
	// fmt.Println("good", err)
	// er(err)

	// err = mongo.Insert("ttt", "ttt2", map[string]interface{}{
	// 	"this":  "is",
	// 	"some":  "3",
	// 	"tests": true,
	// })
	// fmt.Println("bad", err)
	// er(err)

	// var (
	// 	m    map[string]bool
	// 	arr  []string
	// 	m2   map[int]bool
	// 	arr2 []int
	// 	pos  []string
	// 	pos2 []int
	// )

	// pos = []string{"34", "345", "997", "9753"}
	// pos2 = []int{34, 345, 997, 9753}
	// m = make(map[string]bool)
	// m2 = make(map[int]bool)

	// for a := 0; a < 100000; a++ {
	// 	m[fmt.Sprintf("%d", a)] = true
	// 	arr = append(arr, fmt.Sprintf("%d", a))
	// 	m2[a] = true
	// 	arr2 = append(arr2, a)
	// }

	// now = time.Now()
	// for _, i := range pos {
	// 	now = time.Now()
	// 	if ok, _ := m[i]; ok {
	// 		fmt.Println("map", i, ok, time.Since(now))
	// 	}
	// 	now = time.Now()
	// 	ok := exists(i, arr)
	// 	fmt.Println("array", i, ok, time.Since(now))
	// 	now = time.Now()
	// }

	// for _, i := range pos2 {
	// 	if ok, _ := m2[i]; ok {
	// 		fmt.Println("map2", i, ok, time.Since(now))
	// 	}
	// 	now = time.Now()
	// 	ok := existsInt(i, arr2)
	// 	fmt.Println("array", i, ok, time.Since(now))
	// }

	// textA := "AAAABBBBCCCCDDDDEEEE"
	// textB := "11112222333344445555"
	// var c string

	// now = time.Now()
	// for a := 0; a < 10; a++ {
	// 	c = fmt.Sprintf("%s/%s", textA, textB)
	// }
	// fmt.Println("fmt.Sprintf:time:", time.Since(now).String())

	// now = time.Now()
	// for a := 0; a < 10; a++ {
	// 	c = textA + "/" + textB
	// }
	// fmt.Println("sum:time:", time.Since(now).String())

	// fmt.Println("C:", c)

}

func existsInt(item int, items []int) (e bool) {

	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false

}

func exists(item string, items []string) (e bool) {

	for _, i := range items {
		if i == item {
			return true
		}
	}

	return false

}

func er(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}
