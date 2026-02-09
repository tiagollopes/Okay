package main

import "fmt"

func main() {
	// O defer "agenda" a execução para o final da função main
	defer fmt.Println("3. Este será o ÚLTIMO (limpeza final)")
	defer fmt.Println("2. Este será o SEGUNDO (ordem inversa)")

	fmt.Println("1. Este será o PRIMEIRO (execução normal)")

	exemploPratico()
}

func exemploPratico() {
	fmt.Println("--- Início do exemplo prático ---")

	// Simulando a abertura de um recurso (arquivo ou banco)
	fmt.Println("Abrindo conexão com o banco de dados...")

	// Garante que a conexão será fechada, não importa o que aconteça
	defer fmt.Println("Conexão fechada com sucesso!")

	fmt.Println("Realizando consultas processando dados...")

	// Mesmo que a função termine aqui, o defer acima será executado
}
