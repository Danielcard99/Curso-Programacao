# Dicas Sobre Programação

## DOCKER

### Bind Mount (Compartilhar arquivos entre máquina e container)

Sempre que precisamos **compartilhar um arquivo ou uma pasta da máquina com o container**, usamos o **bind mount**. A função principal do bind é exatamente essa: fazer com que o container enxergue e trabalhe com arquivos que estão no nosso sistema local.

Além disso, o bind também é útil quando você **não quer instalar uma tecnologia ou linguagem direto na sua máquina**. Você pode instalar tudo dentro do container e, mesmo assim, manter os arquivos criados lá **salvos na sua máquina local**, já que o diretório fica sincronizado entre os dois ambientes.

### COMO USAR:

```bash
docker run --mount type=bind,source="$(pwd)",target=/app imagem
```

- `source` → sua pasta local

- `target` → onde a pasta aparece dentro do container

### EXEMPLO REAL:

```bash
docker run -it --rm --mount type=bind,source="$(pwd)",target=/app -w /app \ golang:alpine sh
```

- **--mount** → espelha sua pasta local dentro do container.

- **target=/app** → sua pasta local aparece em /app lá dentro.

- **-w /app** → você já entra no container `DENTRO` dessa pasta.

- **sh** → abre o shell.

- **--name** → dá nome ao container

## Portas

- usamos -p para escolher em qual porta exibir o conteúdo no computador

```bash
docker run -d -p 5000:8080 --name gabriele danielcard99/golang:latest
```

- Esquerda (5000) = porta da sua máquina

- Direita (8080) = porta interna do container

- Você pode mudar a porta da máquina sem alterar o código

- Só altere a porta interna se quiser mudar onde a app escuta de fato

## TAG

para dar uma tag (nome + versão) na hora de "buildar" uma imagem usamos -t e o nome que você colocar vai ser o nome da imagem

EXEMPLO:

```bash
docker build -t danielcard99/go:1.0 .
```

Agora a imagem fica exatamente:

`danielcard99/go:1.0`

Podemos usar mais de uma tag na build.

#### 🔍 E o que é cada parte?

No comando:

```bash
docker build -t nome:tag .
```

`nome` → nome da imagem

`tag` → versão ou identificador (ex: v1, prod, alpine, etc.)

`.` → contexto do build (a pasta atual)

## IMAGEM SCRATCH

É UMA IMAGEM VAZIA, LITERALMENTE SEM NADA

- sem Linux
- sem shell
- sem apk, apt, nada
- sem libc
- sem /bin/sh
- sem arquivos do sistema
- sem certificados
- sem timezone
  É um container que começa absolutamente do zero.
  👉 É o menor ambiente possível que o Docker oferece.

A imagem scratch serve para rodar um único binário estático.
Ou seja:

- sem dependências externas
- sem bibliotecas dinâmicas
- apenas 1 arquivo executável

Perfeito para linguagens que geram binários estáticos:

- ✔ Go (com CGO_ENABLED=0)
- ✔ Rust
- ✔ C/C++ static
- ✘ Node, Python, PHP, Java — não funcionam em scratch

#### 📌 Quando usar scratch?

Use scratch quando:

- Você quer a menor imagem possível
- O binário é 100% estático (caso do Go sem CGO)
- Você quer uma imagem segura e minimalista
- Você não precisa rodar comandos dentro do container
- Logs, arquivos e dependências são gerenciados no binário

## UPX(ULTIMATE PACKER FOR EXECUTABLES)

Usamos quando precisamos comprimir um binário.
Ele pega o seu executável (o app do Go, por exemplo) e aplica compressão sem perda de funcionalidade.
Ele não adiciona nada ao container.
Só reduz o tamanho do executável, fazendo com que a imagem final seja menor.

- Seu binário continua 100% funcional
- Ele fica muito menor
- Não exige recompilação do Go
- Funciona para Windows, Linux, macOS, etc.

Por que usar UPX no Dockerfile?

EXEMPLO:

`RUN upx --best --lzma /usr/src/app/app`

Objetivo:

### Reduzir tamanho da imagem final

- Seu binário passa de ~3MB → ~500KB
- Isso ajuda a ficar abaixo do limite de 2MB do desafio Full Cycle

### Complementa flags do Go

- ldflags="-s -w" já remove debug e símbolos → binário menor
- UPX faz a compressão final do executável

### Mantém binário estático

- Funciona junto com CGO_ENABLED=0

### Não afeta performance de execução

- O binário é descomprimido em memória na hora de rodar
- Para apps pequenos, a diferença é desprezível

## ENTRYPOINT

- Define o programa que sempre será executado quando o container iniciar.(Executável Principal)
- Um executável é um programa que pode ser rodado.
  - pode ser:
    - Um binário compilado (como o seu app em Go)
    - Um script shell (sh, bash)
    - Qualquer arquivo que o sistema consiga executar
- É como dizer: “isso é o meu container, ele sempre faz isso”.
- Normalmente usado para binários ou scripts essenciais.
- Uma vez definido não é facilmente sobrescrito via terminal, (somente usando --entrypoint).

Como Usar:
ENTRYPOINT ["/caminho onde fica o executável"]

Exemplo:

`ENTRYPOINT ["/app"]`

- Toda vez que você rodar o container **docker run**, ele vai executar /app.
- Você não precisa digitar nada ao rodar docker run — o container já sabe o que fazer.

## CMD

- Serve para fornecer argumentos padrão ou comando padrão.
- Pode ser sobrescrito ao rodar o container com outros argumentos.

Exemplos:

- CMD ["8080"] `como argumento padrão`
- CMD echo "Olá mundo" `comando shell padrão`

- Se você rodar:
  docker run imagem
- Ele vai usar o CMD

- Mas se você rodar:
  docker run imagem outro-comando
- O CMD será substituído por outro-comando

## COPY

Serve para copiar arquivos ou diretórios do contexto de build local ou de uma imagem anterior(via alias) para dentro da imagem final.

Sintaxe básica:

COPY `origem` `destino`

- `origem` → arquivo ou pasta no seu computador (ou no builder stage)
- `destino` → caminho dentro da imagem

- Se `destino` não existir, o Docker cria automaticamente a pasta ou arquivo no container.
- Se `destino` for um diretório que existe, ele copia para dentro dele.
- Se `destino` for um nome de arquivo, ele copia com esse nome.

Copiando de outro stage (multi-stage)
COPY --from=builder /usr/src/app/app /app

- Copia o binário do stage builder para a imagem final em /app
- Se /app não existir, ele será criado

## WORKDIR

- WORKDIR define o diretório de trabalho dentro do container.
- É como se fosse o comando cd /caminho no Linux.
- Todos os comandos seguintes no Dockerfile que usam caminhos relativos vão usar esse diretório como base.

  9.1.Regras importantes

- Você pode ter vários WORKDIR no Dockerfile.
- Cada novo WORKDIR muda o diretório para os comandos seguintes.

EXEMPLO:

FROM golang:1.21-alpine

### Primeiro WORKDIR

`WORKDIR /usr/src/app`  
`COPY . .`

### Rodar comando dentro de /usr/src/app

`RUN go build -o app .`

### Mudar WORKDIR

`WORKDIR /usr/src/config`  
`COPY config.yaml .`

### Rodar comando dentro de /usr/src/config

`RUN cat config.yaml`

- O primeiro WORKDIR /usr/src/app muda o contexto para /usr/src/app
- Tudo que vem depois (COPY, RUN) usa /usr/src/app como base
- Depois, você muda para /usr/src/config com outro WORKDIR
- Agora, todos os comandos seguintes usam /usr/src/config como base
- Se o diretório não existir, Docker cria automaticamente.
- Combinado com COPY ou RUN, WORKDIR deixa o Dockerfile mais limpo e legível.

Usar vários workdir ajuda a organizar arquivos em pastas diferentes, evitar escrever caminhos longos com COPY ou RUN e melhorar a legibilidade de Dockerfile.

## MYSQL

o mysql tem um mecanismo automático que executa qualquer arquivo .sql ou .sh que estiver dentro da pasta:

`/docker-entrypoint-initdb.d`

Essa pasta NÃO foi você quem escolheu.
Ela é definida pela imagem oficial do MySQL.

`- ./mysql:/docker-entrypoint-initdb.d`

Isso significa:

- ./mysql (sua pasta local)

- será montada DENTRO do container

- no caminho /docker-entrypoint-initdb.d

Tudo que estiver na pasta ./mysql (na sua máquina) será visto pelo MySQL como arquivos de inicialização do banco.

## NGINX

`- ./nginx/default.conf:/etc/nginx/conf.d/default.conf`

Ela está dizendo ao Docker, pegue o arquivo local e monte dentro do container do nginx no caminho: `/etc/nginx/conf.d/default.` Ou seja, o nginx vai usar seu arquivo como configuração principal.

NO container o arquivo tem que se chamar default.conf, o **nginx** por padrão, lê qualquer arquivo .conf dentrod a pasta `/etc/nginx/conf.d/`

## GIT FLOW

quando usamos o git flow hotfix usamos sempre a versão, fazemos isso para que a tag fique correta, e para mostrar qual versão foi corrigida

EXEMPLO:

```bash
git flow hotfix start 0.1.0
```

```bash

git flow hotfix finish 0.1.0
```

Devemos usar o gpg para proteção!

Devemos configurar as Rulesets para que não possam fazer push direto no main, devemos usar PR(Pull REQUEST)

### No GitFlow, a branch principal do repositório é `develop`

E não a `main`.

A lógica é:

`develop` → onde todo desenvolvimento acontece, recebe features, hotfixes, etc.

`main` → só recebe versões estáveis, normalmente com tag (ex: v1.0.0).

### COMO FICA O FLUXO CORRETO

`develop` = branch padrão (default)

- Pull requests vão para ela

- Todo mundo trabalha nela

- Cada feature branch nasce de develop

`main` = branch de releases

- Só recebe merge quando você fecha uma versão

- Geralmente via PR “Release X.X.X”

- Recebe uma tag depois do merge

- Fica sempre estável

### ENTÃO O FLUXO DE MERGE É:

`feature/\* → develop → main (release)`

### Vantagem de deixar `develop` como default:

- Evita PR acidental indo para main

- Mantém a main sempre limpa e estável

- Facilita CI/CD (por exemplo:
  develop → homologação
  main → produção)

Criar uma nova Branch

```bash
git checkout -b nome-da-branch
```

EXEMPLO:

```bash
git checkout -b feature/login
```

#### caso já exista uma branch e você só queira trocar pra ela:

```bash
git checkout nome-da-branch
```

### Como limpar as branches remotas que já foram apagadas no GitHub

Execute:

```bash
git fetch --prune
```

Ou a forma mais moderna:

```bash
git remote prune origin
```

Isso remove do seu computador todas as referências remotas que não existem mais no GitHub.

### Como apagar as branches que só existem remotamente:

```bash
git push origin --delete (nomeDaBranch)
```

### PR TEMPLATES

Para criar Templates para os PRs devemos criar uma pasta .github/PULL_REQUEST_TEMPLATE.md assim podemos criar o template em markdown

## Git Tag

Uma tag é um marcador fixo em um commit.
Normalmente usada para:

- versões (v1.0.0)

- releases

- marcos importantes

Diferente de branch: tag não anda.

### Tipos de tag

1. **Lightweight tag (simples)**

Só um ponteiro para o commit.

2. **Annotated tag (recomendada)**

Tem:

- mensagem

- autor

- data

- pode ser assinada (GPG)

**Use sempre annotated.**

### Criar uma tag (forma correta)

- Tag anotada

```bash
git tag -a v1.0.0 -m "Primeira versão estável"
```

Isso marca o commit atual (HEAD).

- **Tag em um commit específico**

```bash
git tag -a v1.0.0 <hash-do-commit> -m "Primeira versão"
```

### Enviar tag para o GitHub

⚠️ Tags **não sobem automaticamente.**

#### Enviar uma tag

```bash
git push origin v1.0.0
```

#### Enviar todas as tags

```bash
git push origin --tags
```

## SemVer (Semantical Versioning)

### 2.1.4

## (2) - Major

- Usamos quando vamos fazer mudanças grandes, você aumenta esse número quando faz alterações incompatíveis com versões anteriores.

**EXEMPLOS:**

- Remove uma função

- Muda um comportamento de forma que quebra código existente

## (1) - Minor

- Adicionado funcionalidades, mas compatível com a API. aumenta quando adicionamos novas funcionalidades sem quebrar nada já existente.

**EXEMPLOS:**

- Criou um novo endpoint

- Adicionou um parâmetro opcional

- Expandiu recursos sem causar breaking change

## (4) - Patch

- Aumenta quando você faz correções de bugs, sem mudar comportamento e sem adicionar features.

**EXEMPLOS:**

- Corrigiu um erro de validação

- Ajustou um retorno

- Corrigiu um typo que quebrava algo

MAJOR = 0 - API Instável. Pode mudar a qualquer momento.

### Pré-release e builds

#### Pré-release:

Usamos quando temos versões instáveis.  
Usa um - depois da versão:

- 1.0.0-alpha
- 1.0.0-beta.2
- 2.3.1-rc.1 (release candidate 1)

#### Build metadata:

Usa +:

- 1.0.0+exp.sha.5114f85

## Conventional Commits

**Template:**

```
<tipo>[escopo opcional]: <descrição>

[corpo opcional]

[rodapé(s) opcional(is)]
```

## Tipos

#### FIX

DEVE ser usado quando um commit representa a correção de um problema em seu aplicativo ou biblioteca.

**Exemplos:**

- fix: corrigir erro ao atualizar usuário

- fix: resolver problema de CORS no login

#### FEAT

DEVE ser usado quando um commit adiciona um novo recurso ao seu aplicativo ou biblioteca.

**Exemplos:**

- feat: adicionar endpoint de criação de usuário

- feat: implementar upload de arquivos

### DOCS

Deve ser usado somente para mudanças na documentação

**Exemplos:**

- docs: atualizar README

- docs: adicionar exemplos na API

### STYLE

Deve ser usado quando você altera formatação, mas não muda a lógica.

**Exemplos:**

- style: aplicar prettier no projeto

- style: ajustar indentação

### CHORE

Deve ser usado para mudanças que não afetam o código da aplicação

**Exemplos:**

- chore: atualizar dependências

- chore: configurar husky

- chore: mudar script do package.json

### REFACTOR

Quando você altera o código sem mudar comportamento e sem adicionar feature.

**Exemplos:**

- refactor: remover função duplicada

- refactor: otimizar consulta ao banco

### TEST

Deve ser usado quando adicionamos ou alteramos testes.

**Exemplos:**

- test: adicionar testes unitários para UserService

- test: corrigir mock

### PERF

deve ser usado quando a mudança melhora a performance.

**Exemplos:**

- perf: reduzir tempo de resposta no endpoint /posts

- perf: otimizar loop de validação

### BUILD:

Mudanças que afetam o processo de build ou dependências.

**Exemplos:**

- build: configurar webpack

- build: atualizar versão do Node

### CI:

Mudanças em pipelines ou integrações contínuas.

**Exemplos:**

- ci: corrigir workflow de deploy

- ci: adicionar lint no GitHub Actions

## Complementos úteis

### BREAKING CHANGE

Se sua mudança quebra compatibilidade, adicione isso no corpo do commit.

Caso eu não queira usar o breaking change, posso usar uma exclamação após o tipo.

Exemplo:

```bash
feat: alterar retorno do endpoint /users

BREAKING CHANGE: agora o endpoint retorna um array em vez de um objeto.
```

Usar COMMITLINT para não errar nos commits, de acordo com os conventional commits, ele verifica os commits e diz onde está os erros antes dos commits serem feitos!

## COMMITIZEN

Instalação:

```bash
npm install --save-dev commitizen
npm install --save-dev @commitlint/cz-commitlint
```

Para usar o commitizen devemos criar uma pasta na raiz do projeto chamada: .czrc (dizendo qual adaptador usar):

`{ "path": "@commitlint/cz-commitlint" }`

### Configurar o Commitlint

```bash
npm install --save-dev @commitlint/config-conventional @commitlint/cli
```

Crie commitlint.config.js na raiz com:

`export default { extends: ["@commitlint/config-conventional"] };`

### Configurar Husky para hooks de commit

```bash
npm install --save-dev husky
npx husky install
```

### Adicionar ao package.json (scripts):

```
"scripts": {
"prepare": "husky install"
}
```

### Criar hook para commit-msg:

```bash
npx husky add .husky/commit-msg 'npx --no -- commitlint --edit "$1"'
```

### Fluxo de Uso:

- Faça alterações nos arquivos.
- Adicione ao staged:

```bash
git add (nomeDoArquivo)
```

### Rode o Commitizen:

```bash
npx cz
```

O npx cz abrirá um prompt interativo para você escolher tipo, escopo, descrição curta, descrição longa, breaking changes e issues relacionadas.

## CI - Integração Continua

Integração Contínua (CI) é uma prática de desenvolvimento em que alterações de código são integradas frequentemente ao repositório principal e validadas automaticamente.

### Como funciona

- O código é alterado e enviado ao repositório

- Um pipeline automático é executado

- São rodados testes, lint e build

- Se algo falhar, a integração é bloqueada

### Benefícios da Integração Contínua

- Identificação rápida de erros

- Redução de conflitos entre códigos

- Maior qualidade e estabilidade do software

- Facilidade no trabalho em equipe

- Menor risco de falhas em produção

### Github Actions

No GitHub Actions, o CI é configurado por meio de arquivos .yml ou .yaml dentro da pasta `.github/workflows`, onde são definidas as etapas do pipeline, como instalação de dependências, testes e build.
