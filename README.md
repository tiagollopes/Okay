# Linguagem Okay

Uma linguagem focada em microserviços e backend, construída em Go. A Okay transforma definições simples em serviços executáveis com suporte a lógica de negócio e processamento de dados.

## 📋 Status Atual

- **Lexer**: Reconhece símbolos matemáticos (`+`, `-`, `*`, `/`, `=`), comparadores (`>`, `<`, `==`), ignora comentários (`//`) e suporta identificadores complexos (ex: `total_saque`).
- **Parser**: Construção de Árvore de Sintaxe Abstrata (AST) com suporte a:
  - Declaração de variáveis dinâmicas.
  - Expressões matemáticas com as 4 operações básicas.
  - Estruturas de decisão completas (`if/else`).
- **Eval (Interpretador)**:
  - Execução de servidores HTTP nativos por serviço.
  - Gerenciamento de memória em tempo real (Ambiente de variáveis).
  - Resolução de lógica condicional para regras de negócio.

## 🛠️ Exemplo de Poder da Okay

Veja como a Okay resolve uma regra de negócio de saque bancário com taxa:

<pre>
```ok
service banco port 8081 {
    // Definição de valores iniciais
    let saldo = 500;
    let saque = 150;
    let taxa = 5;

    // A Okay resolve variáveis com underline e expressões matemáticas
    let total_saque = saque + taxa;

    // Lógica condicional para autorização
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

[ ] Adicionar suporte a tipos Booleanos (true/false).

[ ] Criar loops de repetição (repeat).

[ ] Implementar captura de parâmetros via URL (Query Params).

***Feito por Tiago LLopes*** Santos/SP


