package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Employee struct {
	ID   int    `json:"id"`
	NAME string `json:"name"`
	AGE  int    `json:"age"`
	DEPT string `json:"dept"`
}

func main() {
	connectTODB()
	createTable()

	app := fiber.New()

	app.Get("/employees", func(c *fiber.Ctx) error {
		var res []Employee
		rows, err := DB.Query("select * from employee")
		if err != nil {
			log.Fatal(err)
		}

		var id int
		var name string
		var age int
		var dept string
		for rows.Next() {
			err = rows.Scan(&id, &name, &age, &dept)
			if err != nil {
				log.Fatal(err)
			}

			res = append(res, Employee{id, name, age, dept})
		}

		defer rows.Close()
		return c.JSON(res)
	})

	app.Get("/employee/:id", func(c *fiber.Ctx) error {
		eid := c.Params("id")

		query := "select * from employee where id = $1"

		row := DB.QueryRow(query, eid)
		var id int
		var name string
		var age int
		var dept string

		err := row.Scan(&id, &name, &age, &dept)

		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(Employee{id, name, age, dept})

	})

	app.Post("/employee", func(c *fiber.Ctx) error {
		var emp Employee
		if err := c.BodyParser(&emp); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		query := "Insert into employee (name, age, dept) values ($1, $2, $3)"

		_, err := DB.Query(query, emp.NAME, emp.AGE, emp.DEPT)

		if err != nil {
			log.Fatal(err)
		}

		// Return JSON response
		return c.JSON(fiber.Map{
			"msg": "User created"})
	})

	app.Patch("/employee/:id", func(c *fiber.Ctx) error {

		eid := c.Params("id")

		query := "Update employee set name=$1 where id = $2"

		DB.QueryRow(query, "updated", eid)

		return c.SendString("Updated")

	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		eid := c.Params("id")
		query := "delete from employee where id = $1"

		DB.QueryRow(query, eid)
		return c.SendString("Deleted")
	})

	app.Listen(":3000")
}

func connectTODB() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	URI := os.Getenv("POSTGRES_URI")

	DB, err = sql.Open("postgres", URI)

	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal(err)
	}
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS EMPLOYEE
	(ID SERIAL PRIMARY KEY,
	NAME TEXT NOT NULL,
	AGE INT,
	DEPT TEXT)
	`

	_, err := DB.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

}
