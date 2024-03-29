package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	cep := "45460000"
	viaCepChannel := make(chan ViaCepResponse)
	brasilApiChannel := make(chan BrasilApiResponse)

	go buscarEnderecoBrasilApi(cep, brasilApiChannel)
	go buscarEnderecoViaCep(cep, viaCepChannel)

	select {
	case viaCepResponse := <-viaCepChannel:
		fmt.Printf("Endereço retornado do ViaCEP: %+v\n", viaCepResponse)
	case brasilApiResponse := <-brasilApiChannel:
		fmt.Printf("Endereço retornado do BrasilAPI: %+v\n", brasilApiResponse)
	case <-time.After(time.Second * 1):
		fmt.Println("Timeout")
	}
}

func buscarEnderecoViaCep(cep string, result chan<- ViaCepResponse) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var data ViaCepResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	result <- data
}

func buscarEnderecoBrasilApi(cep string, result chan BrasilApiResponse) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var data BrasilApiResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}

	result <- data
}
