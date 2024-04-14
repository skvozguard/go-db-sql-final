package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	result, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = ?", number)

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	// менять статус можно только в порядке registered -> sent -> delivered
	var currentStatus string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&currentStatus)
	if err != nil {
		return err
	}

	switch status {
	case ParcelStatusSent:
		if currentStatus != ParcelStatusRegistered {
			return errors.New("можно менять статус только с registered на sent")
		}
	case ParcelStatusDelivered:
		if currentStatus != ParcelStatusSent {
			return errors.New("можно менять статус только с sent на delivered")
		}
	default:
		return errors.New("неправильный статус")
	}

	_, err = s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	var currentStatus string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != ParcelStatusRegistered {
		return errors.New("можно менят адрес только у посылки со статусом registered")
	}

	_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var currentStatus string

	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != ParcelStatusRegistered {
		return errors.New("можно удалять только посылки со статусом registered")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	if err != nil {
		return err
	}

	return nil
}
