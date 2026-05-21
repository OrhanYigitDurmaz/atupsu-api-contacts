package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/alexbrainman/odbc"
)

type Customer struct {
	AboneNo   int     `json:"abone_no"`
	Adi       string  `json:"adi"`
	Soyadi    string  `json:"soyadi"`
	FirmaAdi  string  `json:"firma_adi"`
	Telefon   string  `json:"telefon"`
	CepTel    string  `json:"cep_tel"`
	Telefon2  string  `json:"telefon2"`
	Telefon3  string  `json:"telefon3"`
	Telefon4  string  `json:"telefon4"`
	Telefon5  string  `json:"telefon5"`
	Telefon6  string  `json:"telefon6"`
	Telefon7  string  `json:"telefon7"`
	Telefon8  string  `json:"telefon8"`
	Telefon9  string  `json:"telefon9"`
	Telefon10 string  `json:"telefon10"`
	Mahalle   string  `json:"mahalle"`
	Cadde     string  `json:"cadde"`
	Sokak     string  `json:"sokak"`
	Apartman  string  `json:"apartman"`
	No        string  `json:"no"`
	Kat       string  `json:"kat"`
	Ilce      string  `json:"ilce"`
	Il        string  `json:"il"`
	Aktif     string  `json:"aktif"`
	Latitu    float64 `json:"latitu"`
	Longitu   float64 `json:"longitu"`
}

type DB struct {
	conn *sql.DB
}

func Open(connStr string) (*DB, error) {
	conn, err := sql.Open("odbc", connStr)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}
	return &DB{conn: conn}, nil
}

func ConnStringFromDSN(name string) string {
	return fmt.Sprintf("DSN=%s;", name)
}

func ConnStringFromFile(mdbPath string) string {
	return fmt.Sprintf("driver={Microsoft Access Driver (*.mdb, *.accdb)};DBQ=%s;", mdbPath)
}

func (d *DB) Close() error {
	return d.conn.Close()
}

const selectColumns = `abone_no, adi, soyadi, firma_adi, telefon, cep_tel, telefon2, telefon3, telefon4, telefon5, telefon6, telefon7, telefon8, telefon9, telefon10, mahalle, cadde, sokak, apartman, no, kat, ilce, il, aktif, latitu, longitu`

func scanCustomer(row interface{ Scan(...any) error }) (*Customer, error) {
	var c Customer
	var adi, soyadi, firmaAdi sql.NullString
	var telefon, cepTel, telefon2, telefon3 sql.NullString
	var telefon4, telefon5, telefon6, telefon7 sql.NullString
	var telefon8, telefon9, telefon10 sql.NullString
	var mahalle, cadde, sokak, apartman sql.NullString
	var no, kat, ilce, il, aktif sql.NullString
	var latitu, longitu sql.NullFloat64

	err := row.Scan(
		&c.AboneNo, &adi, &soyadi, &firmaAdi,
		&telefon, &cepTel, &telefon2, &telefon3,
		&telefon4, &telefon5, &telefon6, &telefon7,
		&telefon8, &telefon9, &telefon10,
		&mahalle, &cadde, &sokak, &apartman,
		&no, &kat, &ilce, &il, &aktif,
		&latitu, &longitu,
	)
	if err != nil {
		return nil, err
	}

	c.Adi = adi.String
	c.Soyadi = soyadi.String
	c.FirmaAdi = firmaAdi.String
	c.Telefon = telefon.String
	c.CepTel = cepTel.String
	c.Telefon2 = telefon2.String
	c.Telefon3 = telefon3.String
	c.Telefon4 = telefon4.String
	c.Telefon5 = telefon5.String
	c.Telefon6 = telefon6.String
	c.Telefon7 = telefon7.String
	c.Telefon8 = telefon8.String
	c.Telefon9 = telefon9.String
	c.Telefon10 = telefon10.String
	c.Mahalle = mahalle.String
	c.Cadde = cadde.String
	c.Sokak = sokak.String
	c.Apartman = apartman.String
	c.No = no.String
	c.Kat = kat.String
	c.Ilce = ilce.String
	c.Il = il.String
	c.Aktif = aktif.String
	c.Latitu = latitu.Float64
	c.Longitu = longitu.Float64

	return &c, nil
}

func (d *DB) GetByAboneNo(id int) (*Customer, error) {
	query := fmt.Sprintf("SELECT %s FROM abone WHERE abone_no = ?", selectColumns)
	row := d.conn.QueryRow(query, id)
	return scanCustomer(row)
}

func (d *DB) SearchByPhone(phone string) ([]Customer, error) {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return nil, nil
	}

	conditions := []string{
		"telefon LIKE ?", "cep_tel LIKE ?",
		"telefon2 LIKE ?", "telefon3 LIKE ?",
		"telefon4 LIKE ?", "telefon5 LIKE ?",
		"telefon6 LIKE ?", "telefon7 LIKE ?",
		"telefon8 LIKE ?", "telefon9 LIKE ?",
		"telefon10 LIKE ?",
	}
	where := strings.Join(conditions, " OR ")
	query := fmt.Sprintf("SELECT %s FROM abone WHERE abone_no <> 0 AND (%s)", selectColumns, where)

	pattern := "%" + phone + "%"
	args := make([]any, len(conditions))
	for i := range args {
		args[i] = pattern
	}

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectRows(rows)
}

func (d *DB) SearchByName(query string) ([]Customer, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	// Access string concat uses &
	sql := fmt.Sprintf(
		"SELECT %s FROM abone WHERE abone_no <> 0 AND ((adi & ' ' & soyadi) LIKE ? OR firma_adi LIKE ?)",
		selectColumns,
	)
	pattern := "%" + query + "%"

	rows, err := d.conn.Query(sql, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectRows(rows)
}

func (d *DB) ListAll(limit, lastAboneNo int) ([]Customer, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	q := fmt.Sprintf(
		"SELECT TOP %d %s FROM abone WHERE abone_no <> 0 AND abone_no > ? ORDER BY abone_no",
		limit, selectColumns,
	)

	rows, err := d.conn.Query(q, lastAboneNo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectRows(rows)
}

func collectRows(rows *sql.Rows) ([]Customer, error) {
	var customers []Customer
	for rows.Next() {
		c, err := scanCustomer(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, *c)
	}
	return customers, rows.Err()
}
