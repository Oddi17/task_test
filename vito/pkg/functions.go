package pkg

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func UserBalance(id int) (balance Balances, err error) {

	db, _ := sql.Open("mysql", "root:pass@tcp(172.17.0.1:3307)/MyFirstBD")
	row := db.QueryRow("SELECT ID,ACCOUNT FROM Users WHERE id=?", id)
	err = row.Scan(&balance.ID, &balance.ACCOUNT)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	return
}

func (User Balances) createBalance() (account float32, err error) {
	
	db, _ := sql.Open("mysql", "root:pass@tcp(172.17.0.1:3307)/MyFirstBD")
	var (
		check   float32
		checkID int
	)
	if User.ACCOUNT >= 0 {
		row := db.QueryRow("SELECT ACCOUNT,ID FROM Users WHERE ID = ?", User.ID)
		err = row.Scan(&check, &checkID)
		if err != nil {
			if err == sql.ErrNoRows {
				err = nil
				rs, err2 := db.Exec(
					"INSERT INTO Users(ID,ACCOUNT) VALUES (?,?)",
					User.ID, User.ACCOUNT,
				)
				if err2 != nil {
					log.Fatalln(err2)
				}
				result, err2 := rs.RowsAffected()
				if err2 != nil {
					log.Fatalln(err2)
				}
				account = User.ACCOUNT
				fmt.Println(result)
				defer db.Close()
				return
			}
			log.Fatalln(err)
			return
		}
		if checkID == User.ID {
			value_1 := check
			value_2 := User.ACCOUNT
			new_value := value_1 + value_2
			rs, err2 := db.Exec("UPDATE Users SET ACCOUNT = ? WHERE ID=?", new_value, User.ID)
			if err2 != nil {
				log.Fatalln(err2)
			}
			result, err2 := rs.RowsAffected()
			if err2 != nil {
				log.Fatalln(err2)
			}
			account = new_value
			fmt.Println(result)
			return
		}
	} else {
		err = errors.New("bad payment")
		return
	}
	defer db.Close()
	return
}


func (t Reservation) TransactionReserv() (newValue float32, err error) {
	
	db, _ := sql.Open("mysql", "root:pass@tcp(172.17.0.1:3307)/MyFirstBD")
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(t.ID)
	var savBalance float32
	var price float32
	
	row := tx.QueryRow("SELECT ACCOUNT FROM Users WHERE ID=?", t.ID)
	err = row.Scan(&savBalance)
	if err != nil {
		_ = tx.Rollback()
		log.Fatalln(err)
		return
	}
	fmt.Println(t.ID_SERVICE)
	row = tx.QueryRow("SELECT PRICE FROM Services WHERE ID_SERVICE=?", t.ID_SERVICE)
	err = row.Scan(&price)

	if err != nil {
		_ = tx.Rollback()
		log.Fatalln(err)
		return
	}
	newValue = savBalance - price
	if newValue < 0 {
		_ = tx.Rollback()
		err = errors.New("insufficient funds")
		return
	}
	stmt, err := db.Prepare("INSERT INTO Reservation(ID_USER,ID_SERVICE,AMOUNT) VALUES (?,?,?)")
	if err != nil {
		return
	}
	rs, execErr := stmt.Exec(t.ID, t.ID_SERVICE, price)
	rowsAffected, _ := rs.RowsAffected()
	fmt.Println("exec:", rs, "RowsAffected:", rowsAffected)
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		return
	}

	stmt, err = db.Prepare("UPDATE Users SET ACCOUNT = ? WHERE ID=?")
	if err != nil {
		return
	}
	rs, execErr = stmt.Exec(newValue, t.ID)
	rowsAffected, _ = rs.RowsAffected()
	fmt.Println("exec:", rs, "RowsAffected:", rowsAffected)
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		return
	}
	if err := tx.Commit(); err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	return
}

func (t Transaction) Confirm() (err error) {
	var id_service, amount int
	
	db, _ := sql.Open("mysql", "root:pass@tcp(172.17.0.1:3307)/MyFirstBD")
	row := db.QueryRow(
		"SELECT ID_SERVICE,AMOUNT FROM Reservation WHERE ID_USER=? AND ID_ORDER=?",
		t.ID, t.ID_ORDER,
	)
	err = row.Scan(&id_service, &amount)
	if err != nil {
		log.Fatalln(err)
	}
	rs, err := db.Exec(
		"INSERT INTO Report(ID_USER,ID_SERVICE,ID_ORDER,AMOUNT) VALUES(?,?,?,?)",
		t.ID, t.ID_ORDER, id_service, amount,
	)
	if err != nil {
		log.Fatalln(err)
	}
	rowsAffected, err := rs.RowsAffected()
	fmt.Println("exec:", rs, "RowsAffected:", rowsAffected)
	if err != nil || rowsAffected != 1 {
		log.Fatalln(err)
		return
	}
	rs, err = db.Exec(
		"DELETE FROM Reservation WHERE ID_USER=? AND ID_ORDER=?",
		t.ID, t.ID_ORDER,
	)
	fmt.Println("exec:", rs, "RowsAffected:", rowsAffected)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()
	return
}

func (t Transaction) Reject() (err error) {
	var amount, account int

	db, _ := sql.Open("mysql", "root:pass@tcp(172.17.0.1:3307)/MyFirstBD")
	fmt.Println(t.ID, t.ID_ORDER)
	row := db.QueryRow(
		"SELECT AMOUNT FROM Reservation WHERE ID_USER=? AND ID_ORDER=?",
		t.ID, t.ID_ORDER,
	)
	err = row.Scan(&amount)
	if err != nil {
		log.Fatalln(err)
	}
	row = db.QueryRow("SELECT ACCOUNT FROM Users WHERE ID=?", t.ID)
	err = row.Scan(&account)
	if err != nil {
		log.Fatalln(err)
	}
	newvalue := amount + account
	rs, execErr := db.Exec("UPDATE Users SET ACCOUNT = ? WHERE ID=?", newvalue, t.ID)
	rowsAffected, _ := rs.RowsAffected()
	fmt.Println("exec:", rs, "RowsAffected:", rowsAffected)
	if execErr != nil || rowsAffected != 1 {
		log.Fatalln(err)
		return
	}
	rs, execErr = db.Exec(
		"DELETE FROM Reservation WHERE ID_USER=? AND ID_ORDER=?",
		t.ID, t.ID_ORDER,
	)
	fmt.Println("exec:", rs, "RowsAffected:", rowsAffected)
	if execErr != nil || rowsAffected != 1 {
		log.Fatalln(err)
		return
	}
	defer db.Close()
	return
}
