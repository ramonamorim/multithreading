package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	processarCep("89201300")
}

func processarCep(cep string) {
	ch := make(chan string)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	go BrasilApi(ctx, ch, cep)
	go ViaCep(ctx, ch, cep)

	select {
	case resultado := <-ch:
		println(resultado)
	case <-ctx.Done():
		println("O tempo de execução da requisição foi excedido.")
	}
}

func ViaCep(ctx context.Context, ch chan string, cep string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := string(body)

	ch <- fmt.Sprintf("ViaCep - Resultado: %s", data)
	return &data, nil
}

func BrasilApi(ctx context.Context, ch chan string, cep string) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := string(body)

	ch <- fmt.Sprintf("BrasilAPI - Resultado: %s", body)
	return &data, nil
}
