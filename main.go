// Desafio 2: Crie um programa que receba um CEP como argumento e que retorne a rua, bairro, cidade, estado e o complemento.
// Utilize as APIs:
// - https://brasilapi.com.br/docs
// - https://viacep.com.br/
// O programa deve buscar o CEP nas duas APIs e retornar o resultado da API que responder mais rápido.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type BrasilAPIResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Complemento string `json:"complemento"`
}

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Complemento string `json:"complemento"`
}

func getBrasilAPI(cep string, wg *sync.WaitGroup, ch chan<- BrasilAPIResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), nil)
	if err != nil {
		fmt.Printf("Erro ao criar requisição BrasilAPI para CEP %s: %v\n", cep, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Erro ao obter resposta BrasilAPI para CEP %s: %v\n", cep, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro na resposta BrasilAPI para CEP %s: %v\n", cep, resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler corpo da resposta BrasilAPI para CEP %s: %v\n", cep, err)
		return
	}

	var brasilAPIResponse BrasilAPIResponse
	err = json.Unmarshal(body, &brasilAPIResponse)
	if err != nil {
		fmt.Printf("Erro ao decodificar JSON da resposta BrasilAPI para CEP %s: %v\n", cep, err)
		return
	}

	wg.Done()
	ch <- brasilAPIResponse
}

func getViaCEP(cep string, wg *sync.WaitGroup, ch chan<- ViaCEPResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep), nil)
	if err != nil {
		fmt.Printf("Erro ao criar requisição ViaCEP para CEP %s: %v\n", cep, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Erro ao obter resposta ViaCEP para CEP %s: %v\n", cep, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro na resposta ViaCEP para CEP %s: %v\n", cep, resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler corpo da resposta ViaCEP para CEP %s: %v\n", cep, err)
		return
	}

	var viaCEPResponse ViaCEPResponse
	err = json.Unmarshal(body, &viaCEPResponse)
	if err != nil {
		fmt.Printf("Erro ao decodificar JSON da resposta ViaCEP para CEP %s: %v\n", cep, err)
		return
	}

	wg.Done()
	ch <- viaCEPResponse
}

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scanln(&cep)

	brasilAPI := make(chan BrasilAPIResponse)
	viaCEP := make(chan ViaCEPResponse)

	var wg sync.WaitGroup

	wg.Add(2)

	go getBrasilAPI(cep, &wg, brasilAPI)
	go getViaCEP(cep, &wg, viaCEP)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case brasilAPIResponse := <-brasilAPI:
		fmt.Printf("\nResultado da BrasilAPI:\n")
		fmt.Printf("CEP: %s\n", brasilAPIResponse.Cep)
		fmt.Printf("Logradouro: %s\n", brasilAPIResponse.Logradouro)
		fmt.Printf("Bairro: %s\n", brasilAPIResponse.Bairro)
		fmt.Printf("Localidade: %s\n", brasilAPIResponse.Localidade)
		fmt.Printf("UF: %s\n", brasilAPIResponse.Uf)
		fmt.Printf("Complemento: %s\n", brasilAPIResponse.Complemento)
		fmt.Printf("Fonte: BrasilAPI\n")
		wg.Done()
		break
	case viaCEPResponse := <-viaCEP:
		fmt.Printf("\nResultado da ViaCEP:\n")
		fmt.Printf("CEP: %s\n", viaCEPResponse.Cep)
		fmt.Printf("Logradouro: %s\n", viaCEPResponse.Logradouro)
		fmt.Printf("Bairro: %s\n", viaCEPResponse.Bairro)
		fmt.Printf("Localidade: %s\n", viaCEPResponse.Localidade)
		fmt.Printf("UF: %s\n", viaCEPResponse.Uf)
		fmt.Printf("Complemento: %s\n", viaCEPResponse.Complemento)
		fmt.Printf("Fonte: ViaCEP\n")
		wg.Done()
		break
	case <-ctx.Done():
		fmt.Println("Tempo limite esgotado!")
		break
	}

	wg.Wait()
}
