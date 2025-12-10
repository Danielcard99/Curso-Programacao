# Dicas Sobre Programação

## DOCKER

Sempre que precisamos **compartilhar um arquivo ou uma pasta da máquina com o container**, usamos o **bind mount**. A função principal do bind é exatamente essa: fazer com que o container enxergue e trabalhe com arquivos que estão no nosso sistema local.

Além disso, o bind também é útil quando você **não quer instalar uma tecnologia ou linguagem direto na sua máquina**. Você pode instalar tudo dentro do container e, mesmo assim, manter os arquivos criados lá **salvos na sua máquina local**, já que o diretório fica sincronizado entre os dois ambientes.

### COMO USAR:

```bash
docker run --mount type=bind,source="$(pwd)",target=/app imagem
```

`Source` é o arquivo da máquina e target é o "destino" ou o arquivo do container.

### EXEMPLO REAL:

```bash
docker run -it --rm --mount type=bind,source="$(pwd)",target=/app -w /app \ golang:alpine sh
```

- **--mount** → espelha sua pasta local dentro do container.

- **target=/app** → sua pasta local aparece em /app lá dentro.

- **-w /app** → você já entra no container DENTRO dessa pasta.

- **sh** → abre o shell.

### Para nomear um container use o --name

## Portas

- usamos -p para escolher em qual porta exibir o conteúdo no computador

```bash
docker run -d -p 5000:8080 --name gabriele danielcard99/golang:latest
```

- A porta à direita é sempre a porta interna do container (onde sua app está escutando).

- A porta à esquerda você escolhe livremente, desde que esteja livre no seu computador.

- Você não precisa alterar o código para mudar a porta do host, só se quiser mudar a porta do container.

## TAG

para dar uma tag (nome + versão) na hora de buildar uma imagem usamos -t e o nome que você colocar vai ser o nome da imagem

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

É UMA IMAGEM VAZIA, LITERALEMNTE SEM NADA

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

No seu caso:

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

- Toda vez que você rodar o container, ele vai executar /app.
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

+## COPY

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

quando usamos o git flow hotfix usamos sempre a versão, fazemos isso para que a tag fique correta, e para mostar qual versão foi corrigida

EXEMPLO:

git flow hotfix start 0.1.0

git flow hotfix finish 0.1.0

