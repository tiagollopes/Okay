# Linguagem Okay

Uma linguagem focada em microserviços e backend, construída em Go. A Okay transforma definições simples em serviços executáveis com suporte a lógica de negócio e processamento de dados.

## 📋 Status Atual
- **Lexer**: Reconhece símbolos matemáticos (`+`, `-`, `*`, `/`, `=`) e comparadores (`>`, `<`, `==`).
- **Parser**: Constrói Árvores de Sintaxe Abstrata (AST) com suporte a:
  - Declaração de variáveis.
  - Expressões matemáticas binárias.
  - Estruturas de decisão (`if/else`).
- **Eval (Interpretador)**:
  - Executa servidores HTTP nativos.
  - Gerenciamento de memória (Ambiente de variáveis).
  - Resolução de lógica condicional em tempo de execução.

## 🛠️ Exemplo de Poder da Okay

Abaixo, um exemplo de um microserviço de validação financeira escrito em Okay:

<pre>
```ok
service banco port 8081 {
    let saldo = 500;
    let saque = 150;
    let taxa = 5;

    let total_saque = saque + taxa;

    if (total_saque < saldo) {
        print("Saque autorizado! Total com taxa:", total_saque);
    } else {
        print("Saldo insuficiente. Saldo atual:", saldo);
    }
}
```
</pre>

## Como Executar

1. Certifique-se de ter o Go instalado.

2. Crie ou edite o arquivo teste.ok com seu código.

3. Execute o interpretador:

<pre>
```ok
go run cmd/okay/main.go build teste.ok
```
</pre>

**Acesse no Navegador:** O servidor estará disponível em http://localhost:8081.

## Estrutura do Projeto

**lexer/:** Faz a análise léxica, transformando texto bruto em tokens significativos.

**parser/:** Organiza os tokens em uma estrutura de árvore (AST) que a máquina entende.

**eval/:** O motor de execução onde a lógica é processada e o servidor HTTP é iniciado.

**cmd/:** Ponto de entrada (CLI) da linguagem.

## Próximos Desafios

[ ] Implementar comentários (//) no Lexer.

[ ] Adicionar suporte a tipos Booleanos (true/false).

[ ] Criar loops de repetição (repeat).


