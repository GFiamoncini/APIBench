package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Estrutura Json BrasilCep
type BraCep struct {
	Cep        string `json:"cep"`
	Uf         string `json:"state"`
	Cidade     string `json:"city"`
	Bairro     string `json:"neighborhood"`
	Logradouro string `json:"street"`
}
type BraCepResult struct {
	BraCepSource   string
	BraCepAddress  BraCep
	BraCepDuration time.Duration
	BraCepError    error
}

// Estrutura Json ViaCep
type ViaCep struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Ibge       string `json:"ibge"`
}
type ViaCepResult struct {
	ViaCepSource   string
	ViaCepAddress  ViaCep
	ViaCepDuration time.Duration
	ViaCepError    error
}

func main() {
	cep := "89160222"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//TestCase01 - Simular context timeOut
	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	bracepresultChan := make(chan BraCepResult)
	viacepchan := make(chan ViaCepResult)

	// Iniciar requisições em goroutines
	go func() {
		bracepstart := time.Now()
		bracepaddress, err := BuscaBrasilAPI(ctx, cep)
		bracepduration := time.Since(bracepstart)
		bracepresultChan <- BraCepResult{BraCepSource: "BrasilAPI", BraCepAddress: bracepaddress, BraCepDuration: bracepduration, BraCepError: err}
	}()

	go func() {
		ViaCepstart := time.Now()
		viacepaddress, err := BuscaViaCEP(ctx, cep)
		ViaCepduration := time.Since(ViaCepstart)
		viacepchan <- ViaCepResult{ViaCepSource: "ViaCEP", ViaCepAddress: viacepaddress, ViaCepDuration: ViaCepduration, ViaCepError: err}
	}()

	select {
	case bracepresult := <-bracepresultChan:
		Impressao(bracepresult.BraCepSource, bracepresult.BraCepAddress.Cep, bracepresult.BraCepAddress.Logradouro, bracepresult.BraCepAddress.Bairro, bracepresult.BraCepAddress.Uf, bracepresult.BraCepDuration, bracepresult.BraCepError)
	case resultviacep := <-viacepchan:
		Impressao(resultviacep.ViaCepSource, resultviacep.ViaCepAddress.Cep, resultviacep.ViaCepAddress.Logradouro, resultviacep.ViaCepAddress.Bairro, resultviacep.ViaCepAddress.Uf, resultviacep.ViaCepDuration, resultviacep.ViaCepError)
	case <-ctx.Done():
		log.Println("Tempo limite de 1 segundo excedido.")
	}
}

func BuscaBrasilAPI(ctx context.Context, cep string) (BraCep, error) {
	//Test Case 01 - Simular um atraso de 400 milissegundos
	//time.Sleep(400 * time.Millisecond)

	url := "https://brasilapi.com.br/api/cep/v2/" + cep
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return BraCep{}, fmt.Errorf("erro ao criar requisição para BrasilAPI: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return BraCep{}, fmt.Errorf("erro ao realizar requisição para BrasilAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return BraCep{}, fmt.Errorf("BrasilAPI retornou status code não OK: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BraCep{}, fmt.Errorf("erro ao ler resposta da BrasilAPI: %w", err)
	}

	var address BraCep
	if err := json.Unmarshal(body, &address); err != nil {
		return BraCep{}, fmt.Errorf("erro ao decodificar JSON da BrasilAPI: %w", err)
	}

	log.Printf("BrasilAPI resposta: %+v", address)
	return address, nil
}

func BuscaViaCEP(ctx context.Context, cep string) (ViaCep, error) {

	//Test Case 01 - Simular um atraso de 400 milissegundos
	//time.Sleep(400 * time.Millisecond)

	url := "http://viacep.com.br/ws/" + cep + "/json/"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return ViaCep{}, fmt.Errorf("erro ao criar requisição para ViaCEP: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ViaCep{}, fmt.Errorf("erro ao realizar requisição para ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ViaCep{}, fmt.Errorf("ViaCEP retornou status code não OK: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ViaCep{}, fmt.Errorf("erro ao ler resposta da ViaCEP: %w", err)
	}

	var address ViaCep
	if err := json.Unmarshal(body, &address); err != nil {
		return ViaCep{}, fmt.Errorf("erro ao decodificar JSON da ViaCEP: %w", err)
	}

	log.Printf("ViaCEP resposta: %+v", address)
	return address, nil
}

// Extraido dados para impressão fora do channel
func Impressao(source, cep, logradouro, bairro, uf string, duration time.Duration, err error) {
	if err != nil {
		log.Printf("Erro ao obter dados do %s: %v", source, err)
	} else {
		fmt.Printf("Resposta mais rápida da API %s:\n", source)
		fmt.Printf("Tempo de resposta: %v\n", duration)
		fmt.Printf("CEP: %s\n", cep)
		fmt.Printf("Logradouro: %s\n", logradouro)
		fmt.Printf("Bairro: %s\n", bairro)
		fmt.Printf("UF: %s\n", uf)
	}
}
