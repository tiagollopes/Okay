# Linguagem Okay

Uma linguagem focada em microserviços e backend, construída em Go.

## 📋 Status Atual
- **Lexer**: Reconhece keywords (`service`, `port`, `let`, `print`), identificadores, strings e números.
- **Parser**: Constrói a árvore de sintaxe (AST) com suporte a blocos de serviço `{ }`.
- **Eval (Interpretador)**: Gerencia memória de variáveis e **executa um servidor HTTP real** baseado nas definições do código.

## 🛠️ Como Testar

1. Certifique-se de ter o Go instalado.
2. No arquivo `teste.ok`, defina seu serviço:

<pre>```ok
  service meuapp port 8081 {
       let mensagem = "Ola Mundo! Este dado vem da variavel da Okay.";
       print("Servidor configurado e pronto.");
   }```</pre>


3. Execute o compilador:

*go run cmd/okay/main.go build teste.ok*

4. **Acesse no Navegador**: Abra *http://localhost:8081* para ver a linguagem respondendo em tempo real.

***Estrutura do Projeto***

lexer/: *Transformação de texto em tokens.*

parser/: *Organização da lógica em árvores (AST).*

eval/: *Onde a mágica acontece (Execução e Servidor HTTP).*

cmd/: *Ponto de entrada da aplicação.*
