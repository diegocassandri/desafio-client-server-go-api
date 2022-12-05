package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"-"`
		Codein     string `json:"-"`
		Name       string `json:"-"`
		High       string `json:"-"`
		Low        string `json:"-"`
		VarBid     string `json:"-"`
		PctChange  string `json:"-"`
		Bid        string `json:"bid"`
		Ask        string `json:"-"`
		Timestamp  string `json:"-"`
		CreateDate string `json:"-"`
	} `json:"USDBRL"`
}

var db *sql.DB

func main() {
	dbs, err := sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db = dbs

	http.HandleFunc("/cotacao", cotacaoHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	cotacao, err := buscaCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = salvaCotacao(db, cotacao)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(cotacao)
	log.Println("Request processada com sucesso!")

}

func buscaCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond) //timeout para buscar a cotação na API
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var c Cotacao
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func salvaCotacao(db *sql.DB, cotacao *Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond) // timeout para operações do banco de dados.
	defer cancel()

	stmt, err := db.Prepare("insert into cotacoes(id, valor) values (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, uuid.New().String(), cotacao.Usdbrl.Bid)
	if err != nil {
		return err
	}

	return nil
}
