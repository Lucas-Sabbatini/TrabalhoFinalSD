<h1>Trabalho Final Sistemas Distribuídos

<sub>Servidor Chave - Valor em GO</sub></h1>
- Lucas Sabbatini Janot Procópio
- Pedro Martins Luiz Simoni
- João Caio Pereira Melo
- Gustavo Mascarenhas Amorim


## Visão Geral da Arquitetura

![Diagrama](./diagram.svg)

- **Comunicação Cliente-Nó (gRPC)**: Clientes interagem com qualquer nó no cluster usando gRPC para realizar operações de Put (escrita) e Get (leitura).
- **Comunicação Nó-Nó (MQTT)**: Os nós usam um broker MQTT para difundir atualizações (replicar dados) para todos os outros nós inscritos. Isso permite uma propagação de dados assíncrona e desacoplada.

Claro 🚀 Segue um resumo em **Markdown** de tudo que fizemos até aqui, com os pontos principais organizados:

---

## 1. O que é o arquivo `.proto`?

* O arquivo `.proto` é um **contrato de comunicação** entre cliente e servidor usando **Protocol Buffers (Protobuf)**.
* Nele definimos os "endpoints" da nossa aplicação porém não estamos lidando com o protocolo REST e sim o gRPC. Além disso temos as estruturas de classes essênciais neste projeto como **VectorClockEntry**, **VectorClock**, **Version**, **PutRequest** e **GetRequest**.
* Esse arquivo é **independente de linguagem**: a partir dele, o compilador `protoc` gera código Go, Python, Java ou Rust, conforme necessário.

## 2. O que são os arquivos na pasta `pkg/kvstore`

Depois de rodar o `protoc`, temos dois arquivos principais:

### 🔹 `kv_store.pb.go`

* Define as **estruturas de dados** (mensagens Protobuf).
* Exemplos:

  * `PutRequest`, `PutResponse`
  * `GetRequest`, `GetResponse`
  * `Version`, `VectorClock`, `VectorClockEntry`
* Contém getters, metadados e suporte de serialização para o Protobuf.

### 🔹 `kv_store_grpc.pb.go`

* Define as **interfaces do serviço gRPC**.
* Inclui:

  * `KvStoreClient` → usado pelo cliente para chamar `Put` e `Get`.
  * `KvStoreServer` → interface que o servidor precisa implementar.
* Em resumo: **cola o gRPC ao Go**, permitindo implementar servidor e cliente.

## 3. Etapas para Compilar e Executar

### Passo 1 — Baixar dependências

```bash
go mod tidy
```

### Passo 2 — Compilar servidor e cliente

```bash
go build -o server/bin/server ./server
go build -o client-test/bin/client ./client-test
```

### Passo 3 — Executar

Em um terminal, inicie o servidor:

```bash
./server
```

Em outro terminal, rode o cliente:

```bash
./client
```

Saída esperada:

* No servidor:

  ```
  [PUT] key=foo value=bar
  ```
* No cliente:

  ```
  Resposta Put: success:true
  Resposta Get:
    Valor=bar, Node=node_1
  ```

# Aspectos Futuros para o trabalho
- Implementar tudo, essa é apenas um esqueleto de uma estrutura cliente-servidor
- Foco não será no cliente porém ele foi criado para fins de TESTE.