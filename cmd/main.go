package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type Jogo struct {
	Numeros []int `json:"numeros"`
}

const arquivoJogos = "jogos.json"

// Lê jogos do arquivo local.
func carregarJogos() ([]Jogo, error) {
	var jogos []Jogo
	file, err := os.Open(arquivoJogos)
	if err != nil {
		if os.IsNotExist(err) {
			return []Jogo{}, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jogos)
	if err != nil {
		return nil, err
	}
	return jogos, nil
}

// Salva os jogos no arquivo local.
func salvarJogos(jogos []Jogo) error {
	file, err := os.Create(arquivoJogos)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(jogos)
}

// Valida os números de entrada.
func validarJogo(numeros []int) error {
	if len(numeros) < 6 || len(numeros) > 15 {
		return fmt.Errorf("o jogo deve ter entre 6 e 15 números")
	}
	numSet := make(map[int]bool)
	for _, n := range numeros {
		if n < 1 || n > 60 {
			return fmt.Errorf("os números devem estar entre 1 e 60")
		}
		if numSet[n] {
			return fmt.Errorf("os números do jogo não podem se repetir")
		}
		numSet[n] = true
	}
	return nil
}

// Conta os acertos de um jogo.
func contarAcertos(jogo []int, resultado []int) (int, []int) {
	acertos := 0
	numerosAcertados := []int{}
	resultadoSet := make(map[int]bool)
	for _, n := range resultado {
		resultadoSet[n] = true
	}
	for _, n := range jogo {
		if resultadoSet[n] {
			acertos++
			numerosAcertados = append(numerosAcertados, n)
		}
	}
	return acertos, numerosAcertados
}

// Avalia o resultado de um jogo.
func avaliarJogo(jogo []int, resultado []int) (senas, quinas, quadras, totalAcertos int, numerosAcertados []int) {
	combinacoes := gerarCombinacoes(jogo, 6)
	numerosAcertados = []int{}
	for _, combinacao := range combinacoes {
		acertos, numeros := contarAcertos(combinacao, resultado)
		switch acertos {
		case 6:
			senas++
		case 5:
			quinas++
		case 4:
			quadras++
		}
		// Atualiza total de acertos e números acertados
		if acertos > totalAcertos {
			totalAcertos = acertos
			numerosAcertados = numeros
		}
	}
	return
}

// Gera combinações de números.
func gerarCombinacoes(numeros []int, k int) [][]int {
	var combinacoes [][]int
	var combinacao []int
	var aux func(int, int)
	aux = func(start, depth int) {
		if depth == k {
			combo := make([]int, k)
			copy(combo, combinacao)
			combinacoes = append(combinacoes, combo)
			return
		}
		for i := start; i < len(numeros); i++ {
			combinacao = append(combinacao, numeros[i])
			aux(i+1, depth+1)
			combinacao = combinacao[:len(combinacao)-1]
		}
	}
	aux(0, 0)
	return combinacoes
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "megasena",
		Short: "Gerencie jogos e resultados da Mega-Sena",
	}

	// Comando para adicionar um jogo
	adicionarCmd := &cobra.Command{
		Use:   "adicionar [números]",
		Short: "Adiciona um jogo à lista",
		Args:  cobra.MinimumNArgs(6),
		Run: func(cmd *cobra.Command, args []string) {
			var numeros []int
			for _, arg := range args {
				var n int
				fmt.Sscanf(arg, "%d", &n)
				numeros = append(numeros, n)
			}
			sort.Ints(numeros)

			if err := validarJogo(numeros); err != nil {
				fmt.Println("Erro ao adicionar jogo:", err)
				return
			}

			jogos, err := carregarJogos()
			if err != nil {
				fmt.Println("Erro ao carregar jogos:", err)
				return
			}

			jogos = append(jogos, Jogo{Numeros: numeros})
			if err := salvarJogos(jogos); err != nil {
				fmt.Println("Erro ao salvar jogos:", err)
				return
			}

			fmt.Println("Jogo adicionado:", numeros)
		},
	}

	// Comando para validar um resultado
	validarCmd := &cobra.Command{
		Use:   "validar [resultado]",
		Short: "Valida o resultado de um sorteio contra os jogos salvos",
		Args:  cobra.ExactArgs(6),
		Run: func(cmd *cobra.Command, args []string) {
			var numeros []int
			for _, arg := range args {
				var n int
				fmt.Sscanf(arg, "%d", &n)
				numeros = append(numeros, n)
			}
			sort.Ints(numeros)

			if err := validarJogo(numeros); err != nil {
				fmt.Println("Erro no resultado:", err)
				return
			}

			jogos, err := carregarJogos()
			if err != nil {
				fmt.Println("Erro ao carregar jogos:", err)
				return
			}

			if len(jogos) == 0 {
				fmt.Println("Nenhum jogo salvo para validar.")
				return
			}

			fmt.Println("Resultado:", numeros)

			// Configurar tabela
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"#", "Números Jogados", "Senas", "Quinas", "Quadras", "Acertos", "Números Acertados"})
			for i, jogo := range jogos {
				senas, quinas, quadras, totalAcertos, numerosAcertados := avaliarJogo(jogo.Numeros, numeros)
				table.Append([]string{
					strconv.Itoa(i + 1),
					fmt.Sprint(jogo.Numeros),
					strconv.Itoa(senas),
					strconv.Itoa(quinas),
					strconv.Itoa(quadras),
					strconv.Itoa(totalAcertos),
					fmt.Sprint(numerosAcertados),
				})
			}
			table.Render()
		},
	}

	rootCmd.AddCommand(adicionarCmd, validarCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Erro ao executar o comando:", err)
		os.Exit(1)
	}
}
