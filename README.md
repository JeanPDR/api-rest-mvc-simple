# 📚 Arquitetura em Camadas com Go (Estudo de Caso)

## 🎯 Objetivo do Projeto

Este repositório foi criado com um propósito principal: **estudar, praticar e reforçar os conceitos de Arquitetura de Software em Camadas (Layered Architecture)**, especificamente aplicados na linguagem Go (Golang). 

O projeto implementa uma API REST simples de criação de usuários, mas o foco não está na complexidade da regra de negócio, e sim em **como o código é estruturado e separado**. Ele serve como um guia prático de como aplicar o princípio de Responsabilidade Única (Separation of Concerns).

---

## 🏗️ Conceitos Reforçados

Neste repositório, você encontrará a implementação prática dos seguintes padrões:

*   **Entities (Entidades):** As estruturas de dados "puras" da aplicação. Representam o domínio do negócio e como os dados são (ou seriam) persistidos no banco de dados.
*   **DTOs (Data Transfer Objects):** Estruturas usadas **apenas** para transportar dados de fora para dentro da API (Requests) ou de dentro para fora (Responses). Eles garantem que não vamos expor nossas Entities (e dados sensíveis, como senhas) diretamente para o cliente.
*   **Repositories (Repositórios):** A camada responsável por conversar com o armazenamento de dados. Neste projeto, utilizamos **Interfaces**, o que nos permite ter um banco de dados em memória para testes, podendo ser facilmente substituído por um PostgreSQL ou MongoDB no futuro sem alterar as regras de negócio.
*   **Services (Serviços):** O "cérebro" da operação. Onde moram as regras de negócio (ex: validação de dados, verificação de e-mails duplicados). O Service não sabe nada sobre HTTP e não sabe qual é o banco de dados real.
*   **Handlers / Controllers:** A porta de entrada da API. Recebem a requisição HTTP, extraem o JSON para um DTO, repassam para o Service trabalhar e devolvem o *Status Code* correto para o cliente.

---

## 🚀 Como executar o projeto

1. Clone este repositório:

```bash
git clone https://github.com/seu-usuario/seu-repositorio.git
```

2. Acesse a pasta do projeto:

```bash
cd seu-repositorio
```

3. Execute o servidor:

```bash
go run main.go
```

O servidor estará rodando localmente na porta **8080**.

---

## 🧪 Como testar

Com a aplicação rodando, você pode usar o terminal para simular a criação de um usuário.

### Usando cURL:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ana Souza",
    "email": "ana.souza@exemplo.com",
    "password": "senhaSuperSegura123"
  }'
```

### ✅ Retorno esperado (Status 201 Created):

```json
{
  "id": "uuid-gerado-1234",
  "name": "Ana Souza",
  "email": "ana.souza@exemplo.com"
}
```

> ⚠️ Note que a senha não é retornada, demonstrando o uso correto de um Response DTO.

---
