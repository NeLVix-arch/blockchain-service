package tests

import (
	"blockchain-service/src/database"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// testing connect to DB
func TestInitDB(t *testing.T) {
	db, gorm := database.InitDBS()
	defer db.Close()
	defer gorm.Close()

	err := db.Ping()
	if err != nil {
		t.Fatal(err)
	}

	err = gorm.DB().Ping()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Successfully connected to Postgres")
}

// Функция для создания таблицы wallets с данными id, coin, balance
func createWalletsTable(t *testing.T, db *sqlx.DB) {
	// SQL-запрос для создания таблицы
	query := `CREATE TABLE IF NOT EXISTS wallets (
        id SERIAL PRIMARY KEY,
        coin VARCHAR(20) NOT NULL,
        balance NUMERIC(38, 0) NOT NULL
    );`

	// Выполняем запрос
	_, err := db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

// Функция для создания 2 кошельков в таблице wallets с балансом равным 1000000000000000000
func createWallets(t *testing.T, db *sqlx.DB) {
	// SQL-запрос для вставки данных
	query := `INSERT INTO wallets (coin, balance) VALUES
    ('eth', 1000000000000000000), ('eth', 1000000000000000000);`

	// Выполняем запрос
	_, err := db.Exec(query)
	if err != nil {
		t.Fatal(err)
	}
}

// Тестовая функция для проверки работы кода
func TestCreateWallets(t *testing.T) {
	// Создаем инстанс бд постгрев в базе для теста
	db := InitTESTdb()
	defer db.Close() // Закрываем соединение в конце теста

	// Создаем таблицу wallets с данными id, coin, balance
	createWalletsTable(t, db)

	// Создаем 2 кошелька в таблице wallets с балансом равным 9000000000000000000
	createWallets(t, db)

	// SQL-запрос для получения данных из таблицы
	query := `SELECT id, coin, balance FROM wallets;`

	// Выполняем запрос
	rows, err := db.Query(query)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close() // Закрываем результирующий набор в конце теста

	// Перебираем строки и выводим данные
	for rows.Next() {
		var id int
		var coin string
		var balance int64
		err = rows.Scan(&id, &coin, &balance)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("id: %d, coin: %s, balance: %d\n", id, coin, balance)
	}

	// Проверяем, что нет ошибок при чтении данных
	err = rows.Err()
	if err != nil {
		t.Fatal(err)
	}
}
