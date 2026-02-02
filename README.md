# Linguagem Okay

Uma nova linguagem focada em microserviços e backend, construída em Go.

## 📋 Status Atual
- **Lexer**: Reconhece keywords (`service`, `port`, `let`, `print`), identificadores, strings e números.
- **Parser**: Constrói a árvore de sintaxe (AST) com suporte a blocos de serviço `{ }`.
- **Eval**: Interpretador funcional que gerencia memória de variáveis e execução de comandos.

## 🛠️ Como Testar

1. Certifique-se de ter o Go instalado.
2. Crie ou edite o arquivo `teste.ok`:
   ```ok
    service usuarios port 8080 {
        let versao = "1.0.5";
        let status = "Online";
        print("Servico:", status, "- Versao:", versao);
    }

3. Execute o compilador:

go run cmd/okay/main.go build teste.ok
