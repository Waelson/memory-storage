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

## 📖 Comandos suportados

### `PING`
- **Descrição**: Testa a conectividade com o servidor.
- **Uso**:
  ```plaintext
  PING
  ```
- **Resposta**:
    - `+PONG` se bem-sucedido.
    - Retorna um argumento opcional, se fornecido:
      ```plaintext
      PING "Hello"
      +Hello
      ```

---

### `SET`
- **Descrição**: Armazena uma chave com um valor e um tempo de expiração.
- **Uso**:
  ```plaintext
  SET key value ttl
  ```
- **Parâmetros**:
    - `key`: A chave para armazenar o valor.
    - `value`: O valor a ser armazenado.
    - `ttl`: Tempo de expiração em segundos.
- **Resposta**:
    - `+OK` se bem-sucedido.

---

### `GET`
- **Descrição**: Recupera o valor de uma chave.
- **Uso**:
  ```plaintext
  GET key
  ```
- **Resposta**:
    - O valor associado à chave ou `$-1` se não encontrado.

---

### `DEL`
- **Descrição**: Remove uma chave do armazenamento.
- **Uso**:
  ```plaintext
  DEL key
  ```
- **Resposta**:
    - `:1` se a chave foi removida.
    - `:0` se a chave não existe.

---

### `EXISTS`
- **Descrição**: Verifica se uma chave existe.
- **Uso**:
  ```plaintext
  EXISTS key
  ```
- **Resposta**:
    - `:1` se a chave existe.
    - `:0` se não existe.

---

### `EXPIRE`
- **Descrição**: Define um tempo de expiração (em segundos) para uma chave existente.
- **Uso**:
  ```plaintext
  EXPIRE key ttl
  ```
- **Resposta**:
    - `:1` se a operação foi bem-sucedida.
    - `:0` se a chave não existe.

---

### `TTL`
- **Descrição**: Retorna o tempo restante antes da expiração de uma chave.
- **Uso**:
  ```plaintext
  TTL key
  ```
- **Resposta**:
    - O tempo restante em segundos (`:N`).
    - `:-1` se a chave não tem expiração.
    - `:-2` se a chave não existe ou está expirada.

---

### `FLUSHALL`
- **Descrição**: Remove todas as chaves do armazenamento.
- **Uso**:
  ```plaintext
  FLUSHALL
  ```
- **Resposta**:
    - `+OK` se bem-sucedido.

---

### `INFO`
- **Descrição**: Exibe informações sobre o estado do servidor.
- **Uso**:
  ```plaintext
  INFO
  ```
- **Resposta**:
    - Informações detalhadas sobre uptime, número de chaves e uso de memória, no formato:
      ```plaintext
      # Server
      uptime_in_seconds:<tempo>
  
      # Stats
      number_of_keys:<número>
  
      # Memory
      used_memory:<uso>



## 📂 Estrutura de diretórios
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

## 🛠️ Como rodar o projeto
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