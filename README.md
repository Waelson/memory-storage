# M-Storage 
## Um Redis-like para propósitos educacionais

M-Storage é um servidor Redis-like desenvolvido com **Golang**. Este projeto foi criado com o objetivo de aprendizado e exploração das funcionalidades básicas do Redis, incluindo comandos como `PING`, `SET`, `GET`, `DEL`, entre outros. Embora não seja uma implementação completa do Redis, ele suporta muitas de suas operações fundamentais e segue o protocolo RESP (Redis Serialization Protocol).

⚠️ **Nota:** Este projeto é para fins educacionais e não deve ser usado em produção.

---

## 🚀 Funcionalidades

- **Protocolo RESP**: Compatível com clientes Redis para comandos básicos.
- **Suporte a conexões seguras e não seguras**: Habilite TLS para maior segurança.
- **Persistência**: Dados armazenados em disco usando memória mapeada.
- **TTL e Expiração**: Controle de tempo de vida das chaves.
- **Comandos Redis-like**: Implementação de comandos como `PING`, `SET`, `GET`, `DEL`, entre outros.

---

## 📖 Comandos Suportados

Aqui está a lista dos comandos atualmente suportados, com detalhes sobre cada um:

### `PING`
- **Descrição**: Testa a conectividade com o servidor.
- **Uso**:
  ```plaintext
  PING

## 📂 Estrutura de Diretórios
```
m-storage/
├── cmd/
│   └── main.go                # Ponto de entrada da aplicação
├── internal/
│   ├── resp/
│   │   ├── handler.go         # Gerencia as requisições RESP
│   │   ├── ping.go            # Implementação do comando PING
│   │   ├── set.go             # Implementação do comando SET
│   │   ├── get.go             # Implementação do comando GET
│   │   ├── del.go             # Implementação do comando DEL
│   │   ├── exists.go          # Implementação do comando EXISTS
│   │   ├── expire.go          # Implementação do comando EXPIRE
│   │   ├── ttl.go             # Implementação do comando TTL
│   │   ├── flushall.go        # Implementação do comando FLUSHALL
│   │   ├── info.go            # Implementação do comando INFO
│   ├── server/
│   │   ├── server.go          # Lógica principal do servidor
│   │   ├── storage.go         # Persistência de dados
├── certs/                     # Diretório para certificados TLS
│   ├── ca.crt                 # Certificado da Autoridade Certificadora
│   ├── server.crt             # Certificado público do servidor
│   ├── server.key             # Chave privada do servidor
├── go.mod                     # Arquivo de configuração do módulo Go
├── go.sum                     # Dependências do projeto
├── README.md                  # Documentação do projeto

```

## 🛠️ Como Rodar o Projeto
1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/m-storage.git
cd m-storage
```

2. Instale as dependências:
```bash
go mod tidy
```

3. Execute o servidor:
```bash
go run cmd/main.go
```

4. Conecte-se ao servidor usando `redis-cli` ou ferramentas compatíveis:
```bash
redis-cli -p 6379
```

## ⚙️ Configuração TLS (Opcional)
Para habilitar conexões seguras com TLS:
1. Gere certificados TLS e coloque-os no diretório `certs/.
2. Execute o servidor com a flag --tls`.
```bash
go run cmd/main.go --tls --cert=certs/server.crt --key=certs/server.key --ca=certs/ca.crt 
```

## 📝 Licença
Este projeto é distribuído sob a licença MIT. Consulte o arquivo `LICENSE` para mais detalhes.