package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	//"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	//"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

type Product struct {
	Id         int    // ID
	Model      string // Модел
	Company    string // Компания
	Price      int    // Цена покупка
	Salesprice int    // Цена продажа
	Qty        int    // кол-во в складе
	Created    time.Time
}
type Filters struct {
	Id         int    `json:"id"         form:"id"`
	Model      string `json:"model"      form:"model"`
	Company    string `json:"company"    form:"company"`
	Price      int    `json:"price"      form:"price"`
	SalesPrice int    `json:"salesprice" form:"salesprice"`
	PriceMin   int    `json:"pricemin"   form:"pricemin"`
	PriceMax   int    `json:"pricemax"   form:"pricemax"`
	Page       int    `json:"page"       form:"page"`
	Qty        int    `json:"qty"        form:"qty"`
	Login      string `json:"login"      form:"login"`
	Password   int    `json:"password"   form:"password"`
}

func DeleteByID(c *gin.Context) {
	var buff Filters
	err := c.Bind(&buff)
	products := []Product{}
	if err != nil {
		log.Print(err.Error())
		return
	}
	database.Delete(products, buff.Id).Where(buff.Id)
	//database.Create(products)
	err = database.Delete(products, buff.Id).Where(buff.Id).Error
	//err = database.Raw("DELETE FROM products WHERE id = $1", buff.Id).Scan(&products).Error        так тоже можно удалить
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Удалилься": products})
}

//  <------------------------------------------Пагинация Страница------------------------------------------------------------------------------------------->
func paginationPage(c *gin.Context) {
	var limit int = 10
	var buff Filters
	err := c.Bind(&buff)

	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}

	offset := (buff.Page - 1) * limit

	if err := database.Raw("SELECT * FROM products ORDER BY id LIMIT $1 OFFSET $2", limit, offset).Scan(&products).Error; err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"Result":        products,
		"Страница":      buff.Page,
		"Кол-во данные": len(products),
	})
}

//  <------------------------------------------Поиск от MinPrice до MaxPrice----------------------------------------------------------------------------------------->
func priceFromTo(c *gin.Context) {

	var filter Filters
	err := c.Bind(&filter)

	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT id, model, company, price, salesprice, qty FROM products WHERE price BETWEEN $1 and $2 ", filter.PriceMin, filter.PriceMax).Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"Цена": ""})
	c.JSON(http.StatusOK, gin.H{"От": filter.PriceMin})
	c.JSON(http.StatusOK, gin.H{"До": filter.PriceMax})
	c.JSON(http.StatusOK, gin.H{"Результата Поиска": products})
}

//  <------------------------------------------Поиск по Цену------------------------------------------------------------------------------------------------>
func getByPrice(c *gin.Context) {

	var filter Filters

	err := c.Bind(&filter)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("filter: ", filter)

	products := []Product{}

	if strconv.Itoa(filter.Price) != "" {
		err = database.Raw("SELECT * FROM products WHERE price = $1", strconv.Itoa(filter.Price)).Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		err = database.Raw("SELECT * FROM products").Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	log.Printf("%+v", products)
	c.JSON(http.StatusOK, gin.H{"response": products})
}

//  <------------------------------------------Поиск по Компанию--------------------------------------------------------------------------------------------->
func getByCompany(c *gin.Context) {

	var filter Filters

	err := c.Bind(&filter)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("filter: ", filter)

	products := []Product{}

	if filter.Company != "" {
		err = database.Raw("SELECT * FROM products WHERE company ilike $1", fmt.Sprint("%", filter.Company, "%")).Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {

		err = database.Raw("SELECT * FROM products").Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": products})
}

//  <------------------------------------------Поиск по Моделу----------------------------------------------------------------------------------------------->

//---------------------------------------------Поиск по ID--------------------------------------------------------------------------------------------------
func GetById(c *gin.Context) {

	var filter Filters

	err := c.Bind(&filter)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("filter: ", filter)
	products := []Product{}

	if strconv.Itoa(filter.Id) != "" {
		err = database.Raw("SELECT * FROM products where id = $1", strconv.Itoa(filter.Id)).Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		err = database.Raw("SELECT * FROM products").Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": products})
}

//---------------------------------------------Sort By ID ASC------------------------------------------------------------------------------------------------
func sortByID(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}

	err = database.Raw("SELECT * FROM products ORDER BY id").Scan(&products).Error
	//err = database.Limit(10).Offset(10).Find(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByID": products})
}

//---------------------------------------------Sort By ID DESC------------------------------------------------------------------------------------------------
func sortByIdDesc(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY id DESC").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByIdDesc": products})
}

//---------------------------------------------Sort By Model ASC------------------------------------------------------------------------------------------------
func sortByModel(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY model").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByModel": products})
}

//---------------------------------------------Sort By Model DESC------------------------------------------------------------------------------------------------
func sortByModelDesc(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY model DESC").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByModelDesc": products})
}

//---------------------------------------------Sort By Company ASC------------------------------------------------------------------------------------------------
func sortByCompany(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY company").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByCompany": products})
}

//---------------------------------------------Sort By Company DESC------------------------------------------------------------------------------------------------
func sortByCompanyDesc(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY company DESC").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByCompanyDesc": products})
}

//---------------------------------------------Sort By Price ASC------------------------------------------------------------------------------------------------
func sortByPrice(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY price").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByPrice": products})
}

//---------------------------------------------Sort By Price DESC------------------------------------------------------------------------------------------------
func sortByPriceDesc(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY price DESC").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByPriceDesc": products})
}

//---------------------------------------------Sort By Qty ASC------------------------------------------------------------------------------------------------
func sortByQty(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY qty").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByQty": products})
}

//---------------------------------------------Sort By Qty DESC------------------------------------------------------------------------------------------------
func sortByQtyDesc(c *gin.Context) {
	var sort Filters
	err := c.Bind(&sort)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products ORDER BY qty DESC").Scan(&products).Error
	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"sortByQtyDesc": products})
}
func editProduct(c *gin.Context) {
	var temp Filters
	err := c.Bind(&temp)
	if err != nil {
		log.Print(err.Error())
		return
	}
	products := []Product{}
	err = database.Raw("SELECT * FROM products WHERE id= $1", temp.Id).Scan(&products).Error

	if err != nil {
		log.Print(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"UPDATE": products})
}
func updateAndAddProduct(c *gin.Context) {
	var temp Filters
	err := c.Bind(&temp)
	if err != nil {
		log.Println("ERROR PAGE:", err.Error())
		return
	}
	products := []Product{}
	err = database.Exec("UPDATE products SET model = $1, company = $2, price = $3, salesprice = $4, qty = $5 where id = $6", temp.Model, temp.Company, temp.Price, temp.SalesPrice, temp.Qty, temp.Id).Error

	if err != nil {
		log.Println("ERROR DB:", err.Error())
		return
	}
	c.JSON(404, gin.H{"Изменено": products})
}

// var database *gorm.DB
// type Pro duct struct {
// 	Id         int `json:"id" gorm:"column:id"`            // ID
// 	Model      string `json:"model" gorm:"column:model"`   // Модел
// 	Company    string `json:"company" gorm:"column:company"` // Компания
// 	Price      int  `json:"price" gorm:"column:price"`  // Цена покупка
// 	Salesprice int  `json:"salesprice" gorm:"column:salesprice"`  // Цена продажа
// 	Qty        int   `json:"qty" gorm:"column:qty"` // кол-во в складе
// 	Created    time.Time
// }
// func CreateProduct(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var add Product
// 	_ = json.NewDecoder(r.Body).Decode(&add)
// 	add.Id = strconv.Itoa(rand.Intn(10000000))
// 	products := []Product{}
// 	products = append(products, add)
// 	json.NewEncoder(w).Encode(add)
// }

func CreateProduct(c *gin.Context) {
	var add Filters
	err := c.ShouldBindJSON(&add) //
	//err := c.Bind(&add)
	// Filter:= Product{Model: add.Model,Company: add.Company,Price: add.Price,Salesprice: add.SalesPrice,Qty: add.Qty}
	products := []Product{}
	if err != nil {
		log.Println("bind err:", err.Error())
		return
	}
	log.Println("add: ", add)
	//   err = database.Create(&Filter).Error
	err = database.Raw("INSERT INTO products (model, company, price, salesprice, qty) VALUES (?,?,?,?,?)", add.Model, add.Company, add.Price, add.SalesPrice, add.Qty).Scan(&products).Error
	if err != nil {
		log.Print("db err:", err.Error())
		return
	}
	log.Print("db err:", err.Error())
	c.JSON(http.StatusOK, gin.H{"добавлена": products})
}
func logginAndPassword(c *gin.Context) {
	var template Filters
	var str string = "Неправильный Логин || Парол"
	err := c.Bind(&template)
	if err != nil {
		log.Println("err Bind:", err.Error())
		return
	}
	products := []Product{}
	if template.Password == 85857584 && template.Login == "HLab2" {
		err = database.Raw("SELECT * FROM products").Scan(&products).Error
		if err != nil {
			log.Println("DB ERROR:", err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"Все наши продукты": products})
	} else {
		c.JSON(http.StatusOK, gin.H{"Error": str})
	}
}
func GetByModel(c *gin.Context) {

	var filter Filters

	err := c.Bind(&filter)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("filter: ", filter)

	products := []Product{}

	if filter.Model != "" {
		err = database.Raw("SELECT * FROM products WHERE model ilike $1", fmt.Sprint("%", filter.Model, "%")).Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		err = database.Raw("select *from products").Scan(&products).Error
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"response": products})
}

// func CreateProduct(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	model := r.FormValue("model")
// 	company := r.FormValue("company")
// 	price := r.FormValue("price")
// 	salesprice := r.FormValue("salesprice")
// 	qty := r.FormValue("qty")
// 	products := []Product{}
// 	err = database.Raw("INSERT INTO products(model,company,price,salesprice,qty) VALUES(? , ? , ? , ? , ?)", model, company, price, salesprice, qty).Scan(&products).Error
// 	if err != nil {
// 		log.Println(err.Error())
// 		return
// 	}
// 	http.Redirect(w, r, "/", 301)
// }

func main() {
	dsn := "host=localhost user=app port=5432 dbname=db password=pass sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Println(err)
	}
	database = db

	router := gin.Default() // создаем роутер
	// r := mux.NewRouter()
	// r.HandleFunc("product", CreateProduct).Methods("POST")
	router.GET("/GetByCompany", getByCompany)          // создаем обработчик роутов
	router.GET("/GetByPrice", getByPrice)              // Получение продукт через Price
	router.GET("/GetById", GetById)                    // получение продукт через ID
	router.GET("/GetByModel", GetByModel)              // Получение продукт через Model
	router.POST("/add", CreateProduct)                 // Добавление продукт
	router.GET("/PriceFromTo", priceFromTo)            // Найти все продукты между двух ценах
	router.GET("/pagination", paginationPage)          // П а г и н а ц и я
	router.GET("/sortByID", sortByID)                  // ASC Sort products by ID
	router.GET("sortByModel", sortByModel)             // ASC Sort products by Model
	router.GET("sortByCompany", sortByCompany)         // ASC Sort products by Company
	router.DELETE("/delete", DeleteByID)               // Delete product by ID
	router.GET("/sortByPrice", sortByPrice)            // ASC Sort products by Price
	router.GET("/sortByQty", sortByQty)                // ASC Sort products by Qty
	router.GET("/sortByIdDesc", sortByIdDesc)          // DESC Sort products by ID
	router.GET("sortByModelDesc", sortByModelDesc)     // DESC Sort products by Model
	router.GET("/edit", editProduct)                   // Получение продукт по ID)
	router.PATCH("/edit", updateAndAddProduct)         // Изменим и Добавим
	router.GET("sortByCompanyDesc", sortByCompanyDesc) // DESC Sort products by Company
	router.GET("/sortByPriceDesc", sortByPriceDesc)    // DESC Sort products by Price
	router.GET("/sortByQtyDesc", sortByQtyDesc)        // DESC Sort products by Qty
	router.GET("/login", logginAndPassword)
	//http.ListenAndServe(":8080", nil)
	router.Run("localhost:8080")

}
