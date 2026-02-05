# Linguagem Okay

Uma linguagem focada em microserviços e backend, construída em Go.

## 📋 Status Atual
- **Lexer**: Reconhece símbolos matemáticos (+, - , *, /, =).
- **Parser**: Constrói expressões binárias (Soma).
- **Eval (Interpretador)**:
  - Executa servidores HTTP.
  - Resolve variáveis dinamicamente.
  - Realiza cálculos matemáticos em tempo de execução.

## 🛠️ Como Testar

1. Certifique-se de ter o Go instalado.
2. No arquivo `teste.ok`, defina seu serviço:

<pre>```ok
service calculadora port 8081 {
    let a = 100;
    let b = 20;
    let soma = a + b;
    let sub  = a - b;
    let mult = a * 2;
    let div  = a / b;
    print("Soma:", soma);
    print("Sub:", sub);
    print("Mult:", mult);
    print("Div:", div);
}
```</pre>


3. Execute o compilador:

*go run cmd/okay/main.go build teste.ok*

4. **Acesse no Navegador**: Abra *http://localhost:8081* para ver a linguagem respondendo em tempo real.

***Estrutura do Projeto***

lexer/: *Transformação de texto em tokens.*

parser/: *Organização da lógica em árvores (AST).*

eval/: *Onde a mágica acontece (Execução e Servidor HTTP).*

cmd/: *Ponto de entrada da aplicação.*
