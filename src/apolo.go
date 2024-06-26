package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"./cmd"
	"./utils"

	"gopkg.in/mattn/go-colorable.v0"
)

const (
	version = "0.2.0"
)

func init() {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		utils.PrintError("sistema operacional não suportado.")

		os.Exit(1)
	}

	log.SetFlags(0)

	quiet := false

	for _, v := range os.Args {
		if v == "-q" || v == "--quiet" {
			quiet = true

			break
		}
	}

	cmd.Init(quiet)

	if len(os.Args) < 2 {
		help()

		os.Exit(0)
	}

	for k, v := range os.Args {
		switch v {
		case "-c", "--config":
			log.Print(cmd.GetConfigPath())
			os.Exit(0)
		case "-h", "--help":
			if len(os.Args) > k+1 && os.Args[k+1] == "config" {
				helpConfig()
			} else {
				help()
			}

			os.Exit(0)
		case "-v", "--version":
			log.Print(version)
			os.Exit(0)
		}
	}

	if quiet {
		log.SetOutput(ioutil.Discard)
	} else {
		// suporta saída de impressão em cores para Windows
		log.SetOutput(colorable.NewColorableStdout())
	}
}

func main() {
	utils.PrintBold("apolo v" + version)
	args := os.Args[1:]

	for _, argv := range args {
		switch argv {
		case "backup":
			cmd.Backup()

		case "clear":
			cmd.ClearBackup()

		case "apply":
			cmd.Apply()

		case "update":
			cmd.UpdateCSS()

		case "restore":
			cmd.Restore()

		case "enable-devtool":
			cmd.SetDevTool(true)

		case "disable-devtool":
			cmd.SetDevTool(false)

		case "watch":
			cmd.Watch()

		default:
			if argv[0] != '-' {
				utils.PrintError(`comando "` + argv + `" não encontrado.`)
				utils.PrintInfo(`rode "apolo -h" para lista de comandos válidos.`)
				
				os.Exit(1)
			}
		}
	}
}

func help() {
	fmt.Println("apolo v" + version)
	fmt.Print(`USO
apolo [<flag>] <comando>

DESCRIÇÃO
personaliza a interface e a funcionalidade do cliente spotify

COMANDOS
backup              inicia o backup e o pré-processamento dos arquivos do aplicativo
apply               aplica a customização
update              atualiza o css
restore             restaura o spotify ao estado original
clear               limpa os arquivos de backup atuais
enable-devtool      habilita as ferramentas para desenvolvedores do cliente spotify,
                    (pressione ctrl + shift + i no cliente para começar a usar)
disable-devtool     desativa as ferramentas de desenvolvedor do cliente spotify
watch               entra no modo watch. automaticamente atualiza o css quando o
					arquivo color.ini ou user.css for alterado

FLAGS
-q, --quiet         modo quieto (sem output). cuidado, operações perigosas como
					o backup limpo, procederão sem permissão de prompting
-c, --config        imprime caminho do arquivo de configuração
-h, --help          imprime este texto de ajuda
-v, --version       imprime o número da versão e saia

para informação de configuração, rode "apolo -h config".
`)
}

func.helpConfig() {
	fmt.Print(`SIGNIFICADO DE CONFIGURAÇÃO
[Setting]
spotify_path
    caminho para o diretório do spotify
	
current_theme
    nome da pasta do seu tema
	
inject_css
    se o css personalizado de user.css na pasta do tema é aplicado
	
replace_colors
	se as cores personalizadas são aplicadas

[pré-processos]
disable_sentry
	impede que o sentry envie log/erro/aviso do console aos desenvolvedores do spotify.
	ative se não quiser chamar a atenção deles ao desenvolver extensão ou aplicativo.

disable_ui_logging
	vários elementos registram cada clique e rolagem do usuário.
	ative para interromper o registro e melhorar a experiência do usuário.

remove_rtl_rule
	para oferecer suporte ao árabe e outros idiomas da direita para a esquerda, o spotify adicionou
	muitas regras css que estão obsoletos para usuários da esquerda para a direita.
	ative para remover todos eles e melhorar a velocidade de renderização.

expose_apis
	vaza algumas apis, funções e objetos do spotify para o objeto global do apolo que são
	útil para fazer extensões para estender a funcionalidade do spotify.

[additionaloptions]
experimental_features
	permite o acesso aos recursos experimentais do spotify. abra-os no menu do perfil (canto superior direito).

fastUser_switching
	permite a alteração de conta imediatamente. abra-a no menu do perfil.

home
	habilita a página home. acesse-a na barra lateral esquerda.

lyric_always_show
	força o botão letras para mostrar o tempo todo na barra do player.
	útil para quem deseja assistir a página de visualização.

lyric_force_no_sync
	força a exibição de todas as letras.

made_for_you_hub
	ativa a página feito para você. acesse-a na barra lateral esquerda.

radio
	habilita a página rádio. acesse-a na barra lateral esquerda.

song_page
	cliques no nome da música na barra do player acessarão a página da música
	(em vez da página do álbum) para descobrir listas de reprodução nas quais ele aparece.

visualization_high_framerate
	força a visualização no aplicativo de letras para renderizar em 60fps.
`)
}