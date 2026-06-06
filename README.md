# todotxt

> Aplicação CLI minimalista em Go para gestão de tarefas, 100% compatível com o formato [todo.txt](https://github.com/todotxt/todo.txt).

[![Go Version](https://img.shields.io/badge/go-1.22%2B-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)
[![Standard Library](https://img.shields.io/badge/dependencies-none-blue?style=flat-square)](go.mod)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macOS%20%7C%20windows-lightgrey?style=flat-square)](#instalação)

`todotxt` lê, escreve e manipula diretamente um arquivo `todo.txt` em texto puro, sem base de dados, sem servidor, sem dependências externas. Cada comando altera o arquivo no local e respeita integralmente a especificação do formato, pelo que os seus arquivos permanecem interoperáveis com qualquer outra ferramenta que implemente `todo.txt`.

---

## Índice

- [Porquê?](#porquê)
- [Funcionalidades](#funcionalidades)
- [Instalação](#instalação)
- [Início rápido](#início-rápido)
- [Comandos](#comandos)
- [Filtros do `list`](#filtros-do-list)
- [Formato `todo.txt`](#formato-todotxt)
- [Variáveis de ambiente](#variáveis-de-ambiente)
- [Arquivos](#arquivos)
- [Compatibilidade com a especificação](#compatibilidade-com-a-especificação)
- [Contribuir](#contribuir)
- [Licença](#licença)
- [Créditos](#créditos)

---

## Porquê?

- **Zero dependências.** Apenas a biblioteca padrão de Go. Compila em segundos, binário pequeno (~3 MB), sem cadeia de dependências para auditar.
- **Arquivos legíveis por humanos.** Os seus dados ficam num arquivo de texto que pode abrir em qualquer editor, sincronizar com `git`, `rsync`, `syncthing` ou Dropbox.
- **Compatível com `todo.txt`.** O formato é uma especificação aberta — o seu arquivo não fica preso a esta ferramenta.
- **CLI antes que GUI.** Construído para o terminal: pipes, scripts, aliases, autocompletar, integração com `fzf`, `tmux`, etc.
- **Internacionalizado em português.** Mensagens, ajuda e datas em pt-PT/pt-BR.

---

## Funcionalidades

- ✍️ Adicionar tarefas com prioridade, data de criação, projeto, contexto e data de vencimento
- 📋 Listar tarefas com filtros combináveis (`+Projeto`, `@Contexto`, `pri:A`, `overdue`, `today`, `due`, pesquisa textual)
- ✅ Concluir e reabrir tarefas (mantém o histórico de conclusão)
- 🔥 Definir, alterar e remover prioridades (`A` a `Z`)
- 🗑️ Apagar uma ou várias tarefas
- 📦 Arquivar concluídas para `done.txt`
- 🎨 Saída colorida com indicadores de vencimento (vermelho = vencido, amarelo = até 7 dias)
- 🚫 Suporte da variável `NO_COLOR` e deteção automática de terminais não coloridos
- 🔌 Variável `TODO_DIR` para apontar para qualquer localização
- 🧪 Cobertura de testes do parser com casos do formato real

---

## Instalação

### Via `go install` (recomendado)

```bash
go install github.com/evertonandrade/todotxt@latest
```

Isto coloca o binário em `$(go env GOPATH)/bin/todotxt`. Garanta que esse diretório está no seu `PATH`.

### A partir do código-fonte

```bash
git clone https://github.com/evertonandrade/todotxt.git
cd todotxt
go build -o todotxt .
./todotxt help
```

### Binários pré-compilados

Consulte a [página de releases](https://github.com/evertonandrade/todotxt/releases) — fornecemos binários para `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64` e `windows/amd64`.

### Verificar a instalação

```bash
todotxt version
# todotxt 1.0.0
```

---

## Início rápido

```bash
# 1. Adicionar uma tarefa com prioridade e data de vencimento
todotxt add "Reunião com a equipa +Trabalho @escritório pri:A due:2026-06-15"

# 2. Adicionar mais algumas
todotxt add "Comprar leite +Pessoal @supermercado"
todotxt add "Estudar Go +Estudo @casa"
todotxt add "Pagar contas +Pessoal @casa due:2026-06-30 pri:B"

# 3. Ver a lista (pendentes por defeito)
todotxt list

# 4. Filtrar por projeto
todotxt list +Trabalho

# 5. Ver o que vence hoje ou está vencido
todotxt list overdue
todotxt list today

# 6. Concluir a tarefa 2
todotxt do 2

# 7. Reabrir (se enganou)
todotxt undo 2

# 8. Arquivar concluídas
todotxt archive
```

### Exemplo de saída

```
Tarefas
────────────────────────────────────────────────────────────
  1 [ ] (A) Reunião com a equipa +Trabalho @escritório due:2026-06-15 [vence em 9d: 2026-06-15]
  2 [ ]      Pagar contas +Pessoal @casa due:2026-06-30 [vence em 24d: 2026-06-30]
  3 [ ]      Estudar Go +Estudo @casa
  4 [x] (x)  Comprar leite +Pessoal @supermercado [2026-06-06]

Total: 4 tarefa(s)
```

---

## Comandos

| Comando | Alias | Descrição |
|---|---|---|
| `add <descrição>` | `a` | Adiciona uma nova tarefa |
| `list [filtros]` | `ls` | Lista tarefas (pendentes por defeito) |
| `do <n>` | `x` | Marca a tarefa `n` como concluída |
| `undo <n>` | `unx` | Reabre a tarefa `n` |
| `pri <n> <A-Z>` | `p` | Define a prioridade da tarefa `n` |
| `depri <n>` | `dp` | Remove a prioridade da tarefa `n` |
| `del <n>...` | `rm` | Remove uma ou mais tarefas |
| `archive` | — | Move concluídas para `done.txt` |
| `help` | `-h`, `--help` | Mostra a ajuda |
| `version` | `-v`, `--version` | Mostra a versão |

### `add`

```bash
todotxt add "Descrição da tarefa +Projeto @Contexto pri:A due:YYYY-MM-DD"
```

A data de criação é inserida automaticamente. Os tokens `pri:` e `due:` são extraídos da descrição e formatados na linha final.

### `list`

Sem argumentos, mostra as tarefas pendentes (ocultando concluídas). Veja [Filtros do `list`](#filtros-do-list) para o conjunto completo.

### `do` / `undo`

A numeração corresponde à posição **na listagem atual**. Como a listagem é ordenada (concluídas em último, depois prioridade, depois data de vencimento), os números podem mudar após edição. Use o `list` antes de operações destrutivas.

### `archive`

Move todas as tarefas concluídas do `todo.txt` para o `done.txt` (criado se não existir), preservando a ordem cronológica.

---

## Filtros do `list`

Podem ser combinados livremente. Os filtros `+proj`, `@ctx` e `pri:X` aplicam lógica **E** entre si.

| Filtro | Significado |
|---|---|
| `+Projeto` | Mantém apenas tarefas que contenham `+Projeto` |
| `@Contexto` | Mantém apenas tarefas que contenham `@Contexto` |
| `+A +B` | Tarefas que tenham **ambos** os projetos |
| `@a @b` | Tarefas que tenham **ambos** os contextos |
| `pri:A` | Mantém tarefas com prioridade `A` (ou `A`, `B`, ..., `Z`) |
| `all` | Mostra pendentes **e** concluídas |
| `done` | Mostra apenas concluídas |
| `due` | Apenas tarefas com data de vencimento definida |
| `overdue` | Apenas tarefas vencidas (`due` no passado e não concluídas) |
| `today` | Apenas tarefas que vencem hoje |
| `texto livre` | Pesquisa textual na descrição (case-insensitive) |

Exemplos:

```bash
todotxt list +Trabalho
todotxt list @casa +Pessoal
todotxt list pri:A due
todotxt list overdue +Trabalho
todotxt list comprar
```

---

## Formato `todo.txt`

Cada linha do arquivo é uma tarefa. O parser e o serializador implementam a [especificação oficial](https://github.com/todotxt/todo.txt#specification).

### Estrutura

```
x 2024-01-15 2024-01-10 (B) Call Mom +Family @phone due:2024-01-20
^ ^---------^ ^---------^ ^^      ^----------------------^----------^
| |           |           ||      |                         |
| |           |           ||      |                         +-- campos chave:valor
| |           |           ||      +-- descrição (com +proj e @ctx embebidos)
| |           |           |+-- prioridade
| |           |           +-- (B) ou ausente
| |           +-- data de criação (se não concluída, segue a prioridade; se concluída, segue a data de conclusão)
| +-- data de conclusão
+-- marcador de concluída
```

### Regras implementadas

- **Prioridade** — `(A)` a `(Z)`, sempre em maiúsculas; ausente = sem prioridade. Em tarefas concluídas a prioridade é ignorada.
- **Conclusão** — `x` no início, seguido da data de conclusão, depois da data de criação.
- **Datas** — sempre `YYYY-MM-DD`. Datas inválidas são rejeitadas pelo `add`.
- **Projetos** — qualquer palavra começada por `+` na descrição (ex.: `+Trabalho`).
- **Contextos** — qualquer palavra começada por `@` na descrição (ex.: `@casa`).
- **Campos personalizados** — `chave:valor` para qualquer `chave` alfanumérica. Reconhecidos especialmente: `due:YYYY-MM-DD` (data de vencimento) e `t:YYYY-MM-DD` (data limite).
- **URLs** — não são interpretadas como campos personalizados, mesmo contendo `:`.

Consulte a [lista de issues](https://github.com/evertonandrade/todotxt/issues) para o roadmap completo.

---

## Variáveis de ambiente

| Variável | Efeito | Predefinição |
|---|---|---|
| `TODO_DIR` | Diretório que contém `todo.txt` e `done.txt` | Diretório atual |
| `TODO_FILE` | Caminho completo para um `todo.txt` alternativo | `todo.txt` (ou `$TODO_DIR/todo.txt`) |
| `NO_COLOR` | Qualquer valor desativa a saída colorida | Cores ativas |
| `TERM=dumb` | Desativa a saída colorida | Cores ativas |

---

## Arquivos

| Arquivo | Conteúdo |
|---|---|
| `todo.txt` | Tarefas ativas. Criado automaticamente no primeiro `add`. |
| `done.txt` | Tarefas concluídas e arquivadas. Criado no primeiro `archive`. |

A localização predefinida é `./todo.txt` e `./done.txt` (ou `done.txt` no `TODO_DIR`, se definido).

### Exemplo de `todo.txt`

```
(A) 2026-06-06 Reunião com equipa +Trabalho @escritório due:2026-06-15
(B) 2026-06-05 Pagar contas +Pessoal @casa due:2026-06-30
2026-06-06 Estudar Go +Estudo @casa
x 2026-06-04 2026-06-03 Comprar leite +Pessoal @supermercado
```

### Exemplo de `done.txt`

```
x 2026-06-04 2026-06-03 Comprar leite +Pessoal @supermercado
x 2026-05-30 2026-05-28 Renovar carta de condução +Pessoal
```

---

## Compatibilidade com a especificação

`todotxt` foi testado contra os exemplos canónicos do projeto [`todo.txt`](https://github.com/todotxt/todo.txt) e implementa o seguinte subconjunto:

- ✅ Marcação de conclusão (`x` + data)
- ✅ Prioridades `(A)`–`(Z)`
- ✅ Datas de criação e conclusão
- ✅ Projetos `+x` e contextos `@y`
- ✅ Campos personalizados `chave:valor` (incluindo `due:` e `t:`)
- ✅ Roundtrip preservado: `Parse → Format → Parse` é idempotente

O que **não** segue a especificação:

- ❌ Ordenação alfabética ou por prioridade no arquivo — usamos a ordem introduzida pelo utilizador. Use `list` (que ordena na saída) ou um `sort` se precisar.

---

## Contribuir

Contribuições são bem-vindas! Por favor abra uma [issue](https://github.com/evertonandrade/todotxt/issues) antes de submeter uma pull request significativa.

### Fluxo de desenvolvimento

```bash
# 1. Bifurcar e clonar
git clone https://github.com/evertonandrade/todotxt.git
cd todotxt

# 2. Criar um ramo
git checkout -b minha-feature

# 3. Compilar e testar
go build -o bin/todotxt .
go test ./...
go vet ./...
gofmt -l .    # deve sair vazio

# 4. Submeter PR
git push origin minha-feature
```

### Diretrizes

- **Sem dependências externas.** Se a sua contribuição precisar de uma, abra primeiro uma issue para discussão.
- **Estilo.** `gofmt` + `go vet` limpos. Mensagens de commit em português ou inglês, no imperativo.
- **Testes.** Adicione testes para novas funcionalidades e casos do parser.
- **Compatibilidade.** Não quebre o formato `todo.txt` nem introduza mudanças incompatíveis sem aviso.

### Ideias para contribuir

- Suporte para datas recorrentes (`rec:1w`).
- Comando `edit` para modificar uma tarefa existente.
- Importação/exportação de outros formatos (CSV, JSON, Markdown).
- Autocompletar para `bash`, `zsh` e `fish`.
- Suporte multi-arquivo (`TODO_DIR` com vários projetos).
- Internacionalização (i18n) das mensagens.
- Integração com `fzf` para seleção interativa.

---

## Licença

Este projeto é distribuído sob a licença **MIT**. Veja [`LICENSE`](LICENSE) para o texto completo.

```
MIT License

Copyright (c) 2026 todotxt contributors

É concedida permissão, sem custos, a qualquer pessoa que obtenha uma cópia
deste software e arquivos de documentação associados...
```

---

## Créditos

- O formato e a especificação [`todo.txt`](https://github.com/todotxt/todo.txt) são obra da comunidade em torno do projeto original de [Gina Trapani](https://github.com/ginatrapani) e contributors.
- Inspirado pela [`todo.txt-cli`](https://github.com/todotxt/todo.txt-cli) (shell script) e pelas várias implementações listadas em [github.com/todotxt](https://github.com/todotxt).
- Construído com a [biblioteca padrão de Go](https://pkg.go.dev/std).

---

<div align="center">
  <sub>Se esta ferramenta lhe é útil, considere dar uma ⭐.</sub>
</div>
