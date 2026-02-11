# Linguagem Okay

Uma linguagem focada em microserviços e backend, construída em Go. A Okay transforma definições simples em serviços executáveis com suporte a lógica de negócio e processamento de dados.

## 📋 Status Atual

- **Lexer**: Suporte a símbolos matemáticos, comparadores, comentários (`//`), identificadores com `_` e palavras-chave booleanas (`true/false`).
- **Parser**: Árvore de Sintaxe Abstrata (AST) com suporte a:
  - Variáveis e Expressões Matemáticas.
  - Estruturas `if/else` com suporte a condições diretas (booleanas) ou comparativas.
  - Estrutura de repetição `repeat`.
- **Eval (Interpretador)**:
  - Gerenciamento de estados lógicos e API Dinâmica (Integração total com Query Params da URL).
  - Inicialização silenciosa de serviços.
  - Execução de loops e condicionais em tempo de execução via requisições HTTP.

## 🛠️ Exemplo de Poder da Okay

Este exemplo demonstra a Okay processando um microserviço de checkout com regras de negócio, lógica booleana e variáveis complexas:

<pre>
```ok
      service checkout port 8081 {
          // 1. Configurações (Podem ser mudadas via URL: ?tentativas=5)
          let cupom_ativo = true;
          let frete_gratis = false;
          let tentativas = 2;

          // 2. Valores
          let produto_preco = 150;
          let desconto = 20;
          let taxa_entrega = 15;

          // 3. Lógica Matemática
          let total_com_desconto = produto_preco - desconto;

          // 4. Teste de Loop (Processamento repetitivo)
          print("Iniciando verificacao de seguranca...");
          repeat tentativas {
              print("-> Checando integridade do pedido...");
          }

          // 5. Lógica de Condicional
          if (cupom_ativo) {
              print("Cupom aplicado! Subtotal:", total_com_desconto);
          }

          if (frete_gratis) {
              print("Frete: R$ 0");
          } else {
              let total_final = total_com_desconto + taxa_entrega;
              print("Valor final com frete:", total_final);
          }
      }
```
</pre>

## Como Executar

1. Certifique-se de ter o Go instalado.

2. Crie ou edite o arquivo teste.ok com seu código.

3. Execute o interpretador:

<pre>go run cmd/okay/main.go build teste.ok</pre>

**Acesse no Navegador:** O servidor estará disponível em <pre>```http://localhost:8081/?cupom_ativo=false&tentativas=10&produto_preco=500```</pre>.

## Estrutura do Projeto

**lexer/:** Faz a análise léxica, transformando texto bruto em tokens significativos.

**parser/:** Organiza os tokens em uma estrutura de árvore (AST) que a máquina entende.

**eval/:** O motor de execução onde a lógica é processada e o servidor HTTP é iniciado.

**cmd/:** Ponto de entrada (CLI) da linguagem.

## Próximos Desafios

[ ] Implementar concatenação de strings.

***Feito por Tiago LLopes*** Santos/SP


