# go-mega-sena
CLI simples para validar jogos da mega sena localmente

## Como usar
`go run cmd/main.go adicionar [números do jogo]`

Exemplo: `go run cmd/main.go adicionar 1 2 3 4 5 6` irá salvar no arquivo local `jogos.json` que você fez o jogo 1 2 3 4 5 6.

Após adicionar todos os jogos que você quer conferir execute o comando

`go run cmd/main.go validar 1 2 3 4 5 6`, onde 1 2 3 4 5 6 foram os números sorteados.

Você irá ter uma saída similar a abaixo

```shell
go run cmd/main.go validar 1 17 19 29 50 57
Resultado: [1 17 19 29 50 57]
+----+-----------------------------+-------+--------+---------+---------+-------------------+
| #  |       NÚMEROS JOGADOS       | SENAS | QUINAS | QUADRAS | ACERTOS | NÚMEROS ACERTADOS |
+----+-----------------------------+-------+--------+---------+---------+-------------------+
|  1 | [1 6 26 41 49 50]           |     0 |      0 |       0 |       2 | [1 50]            |
|  2 | [11 14 15 17 21 44]         |     0 |      0 |       0 |       1 | [17]              |
+----+-----------------------------+-------+--------+---------+---------+-------------------+
```
Infelizmente eu não ganhei na mega da virada, pelo menos tem esse código merda gerado em chatgpt pra validar.

Se achar algum bug faz um fork aí na moral