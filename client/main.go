package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Usdbrl struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	cotacao, err := buscaCotacao()
	if err != nil {
		panic(err)
	}

	err = gravaArquivoCotacao(cotacao)
	if err != nil {
		panic(err)
	}
}

func buscaCotacao() (cotacao *Cotacao, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond) //timeout para buscar cotação do servidor
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return nil, error

	}
	var c Cotacao
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}
	return &c, nil

}

func gravaArquivoCotacao(Cotacao *Cotacao) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}

	tamanho, err := fmt.Fprintf(file, "Dólar: %s", Cotacao.Usdbrl.Bid)
	if err != nil {
		return err
	}
	fmt.Printf("Arquivo criado com sucesso! Tamanho: %d bytes\n", tamanho)
	return nil
}
