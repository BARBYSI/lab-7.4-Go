package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)



func isExists(id int) bool {
	db, err := sql.Open("postgres", config())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Query("select * from flights where id= $1", id)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	if result.Next() {
		return true
	} else {
		return false
	}
}

func printAvailable() {
	type flight struct {
		id          string
		destination string
		time        string
		price       string
		seats       string
	}

	db, err := sql.Open("postgres", config())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from flights")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	flights := []flight{}

	for rows.Next() {
		res := flight{}
		err := rows.Scan(&res.id, &res.destination, &res.time, &res.price, &res.seats)
		if err != nil {
			fmt.Println(err)
			continue
		}
		flights = append(flights, res)
	}

	for index := range flights {
		fmt.Println(flights[index])
	}
	menu()
}

func buyTicket() {

	type flight struct {
		id          string
		destination string
		time        string
		price       string
		seats       string
	}

	db, err := sql.Open("postgres", config())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user_input string
	fmt.Println("Введите ID полёта для покупки билета")
	fmt.Scan(&user_input)

	int_ID, err := strconv.Atoi(user_input)
	if err != nil {
		panic(err)
	}

	if isExists(int_ID) {
		row := db.QueryRow("select * from flights where id = $1", user_input)
		res := flight{}
		err = row.Scan(&res.id, &res.destination, &res.time, &res.price, &res.seats)
		if err != nil {
			panic(err)
		}

		if res.seats == "0" {
			fmt.Println("Свободных мест нет")
		} else {
			int_seats, err := strconv.Atoi(res.seats)
			if err != nil {
				panic(err)
			}
			int_seats -= 1
			res, err := db.Exec("update flights set seats = $1 where id = $2", int_seats, user_input)
			if err != nil {
				panic(err)
			}
			log.Print(res)
			fmt.Println("Билет успешно куплен!")
			menu()
		}
	} else {
		fmt.Println("Введите существующий айди!")
		menu()
	}

}

func update() {

	type flight struct {
		id          string
		destination string
		time        string
		price       string
		seats       string
	}

	db, err := sql.Open("postgres", config())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var flight_to_edit string
	fmt.Println("Введите ID полета для редактирования:")
	fmt.Scanln(&flight_to_edit)

	int_ID, err := strconv.Atoi(flight_to_edit)
	if err != nil {
		panic(err)
	}
	//проверка существует ли айди
	if isExists(int_ID) {
		row := db.QueryRow("select * from flights where id = $1", flight_to_edit)
		res := flight{}
		err = row.Scan(&res.id, &res.destination, &res.time, &res.price, &res.seats)

		if err != nil {
			panic(err)
		}
		fmt.Println(res.id, res.destination, res.time, res.price, res.seats)

		var user_input [5]string
		str := [5]string{"Введите новый ID:", "Введите пункт назначения:", "Введите время:", "Введите цену:", "Введите количество оставшихся мест:"}
		for i := 0; i < len(str); i++ {
			fmt.Println(str[i])
			fmt.Scanln(&user_input[i])
		}
		int_ID, err = strconv.Atoi(user_input[0])
		if err != nil {
			panic(err)
		}
		// проверка не занят ли айди
		if !isExists(int_ID) {
			result, err := db.Exec("update flights set id = $1, destination = $2, time = $3, price = $4, seats = $5 where id = $6", user_input[0], user_input[1], user_input[2], user_input[3], user_input[4], flight_to_edit)
			if err != nil {
				panic(err)
			}
			fmt.Println(result)
			fmt.Printf("Запись ID: %s обновлена:\nID: %s\nПункт назначения %s\nВремя отправления: %s\nЦена: %s\nСвободных мест: %s", flight_to_edit, user_input[0], user_input[1], user_input[2], user_input[3], user_input[4])
		} else {
			fmt.Println("Введите уникальный ID!")
			update()
		}

	} else {
		fmt.Println("Введите существующий ID!")
		update()
	}
	menu()

}

func insert() {
	db, err := sql.Open("postgres", config())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user_input [5]string
	str := [5]string{"Введите ID:", "Введите пункт назначения:", "Введите время:", "Введите цену:", "Введите количество оставшихся мест:"}
	for i := 0; i < len(str); i++ {
		fmt.Println(str[i])
		fmt.Scanln(&user_input[i])
	}
	int_ID, err := strconv.Atoi(user_input[0])
	if err != nil {
		panic(err)
	}
	if !isExists(int_ID) {
		result, err := db.Exec("insert into flights(id, destination, time, price, seats) values ($1, $2, $3, $4, $5)", user_input[0], user_input[1], user_input[2], user_input[3], user_input[4])
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
		fmt.Printf("Добавлена запись:\nID: %s\nПункт назначения %s\nВремя отправления: %s\nЦена: %s\nСвободных мест: %s", user_input[0], user_input[1], user_input[2], user_input[3], user_input[4])
	} else {
		fmt.Println("Введите уникальный ID!")
		insert()
	}
	menu()
}

func delete() {

	type flight struct {
		id          string
		destination string
		time        string
		price       string
		seats       string
	}

	db, err := sql.Open("postgres", config())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var (
		user_input        string
		user_confirmation string
	)
	fmt.Println("Введите ID записи для удаления: ")
	fmt.Scan(&user_input)
	row := db.QueryRow("select * from flights where id = $1", user_input)
	res := flight{}
	err = row.Scan(&res.id, &res.destination, &res.time, &res.price, &res.seats)

	if err != nil {
		panic(err)
	}
	fmt.Println(res.id, res.destination, res.time, res.price, res.seats)
	fmt.Println("Вы действительно желаете удалить запись с ID $1?\nPrint yes or no(y/n):", user_input)
	fmt.Scan(&user_confirmation)
	int_ID, err := strconv.Atoi(user_input)
	if err != nil {
		panic(err)
	}

	if user_confirmation == "n" {
		menu()
	}

	if isExists(int_ID) {
		result, err := db.Exec("delete from flights where id = $1", user_input)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)
		fmt.Println("Запись с ID $1 успешно удалена!", user_input)
		menu()
	}
}

func menu() {
	counter := 0
	if counter == 0 {
		fmt.Println("Жмых Airlines menu ver 1.0")
	}
	var user_input string
	fmt.Println("Введите цифру от 1 до 5, чтобы активировать пункт меню")
	fmt.Println("1. Доступные полёты\n2. Купить билет\n3. Добавить полёт\n4. Редактировать полёт\n5. Удалить полёт\n6. Выход")
	fmt.Scanln(&user_input)
	switch user_input {
	case "1":
		printAvailable()
	case "2":
		buyTicket()
	case "3":
		insert()
	case "4":
		update()
	case "5":
		delete()
	case "6":
		os.Exit(0)
	default:
		menu()
	}

}

func main() {


	config()
	menu()
}
