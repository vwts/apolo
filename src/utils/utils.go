package utils

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-ini/ini"
)

// readanswer printa um form sim/não com uma string de `info`
// e retorna um valor booleano baseado no input do usuário (y/Y ou n/N) ou
// retorna `defaultAnswer` se o input for omitido
//
// se o input for nenhum deles, printar o form novamente
func ReadAnswer(info string, defaultAnswer bool) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(info)

	text, _ := reader.ReadString('\n')

	text = strings.Replace(text, "\r", "", 1)
	text = strings.Replace(text, "\n", "", 1)

	if len(text) == 0 {
		return defaultAnswer
	} else if text == "y" || text == "Y" {
		return true
	} else if text == "n" || text == "N" {
		return false
	}

	return ReadAnswer(info, defaultAnswer)
}

// checkexistandcreate checa a existência da pasta
// ou então cria essa pasta caso ela não exista
func CheckExistAndCreate(dir string) {
	_, err := os.Stat(dir)

	if err != nil {
		os.Mkdir(dir, 0700)
	}
}

// unzip descompacta arquivos zips
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)

	if err != nil {
		return err
	}

	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()

		if err != nil {
			return err
		}

		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, 0700)
		} else {
			var fdir string

			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, 0700)

			if err != nil {
				log.Fatal(err)

				return err
			}

			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				return err
			}

			defer f.Close()

			_, err = io.Copy(f, rc)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// copy
func Copy(src, dest string, recursive bool, filters []string) error {
	dir, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	os.MkdirAll(dest, 0700)

	for _, file := range dir {
		fSrcPath := filepath.Join(src, file.Name())
		fDestPath := filepath.Join(dest, file.Name())

		if file.IsDir() && recursive {
			os.MkdirAll(fDestPath, 0700)

			Copy(fSrcPath, fDestPath, true, filters)
		} else {
			if filters != nil && len(filters) > 0 {
				isMatch := false

				for _, filter := range filters {
					if strings.Contains(file.Name(), filter) {
						isMatch = true

						break
					}
				}

				if !isMatch {
					continue
				}
			}

			fSrc, err := os.Open(fSrcPath)

			if err != nil {
				return err
			}

			defer fSrc.Close()

			fDest, err := os.OpenFile(
				fDestPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				return err
			}

			defer fDest.Close()

			_, err = io.Copy(fDest, fSrc)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// replace usa regexp para encontrar qualquer coincidência do `input` com `regexpterm`
// e troca eles por `replaceterm` e então retorna uma nova string
func Replace(input string, regexpTerm string, replaceTerm string) string {
	re := regexp.MustCompile(regexpTerm)

	return re.ReplaceAllString(input, replaceTerm)
}

// modifyfile abre um arquivo, altera conteúdo de arquivo executando
// a função de callback `repl` e escreve novo conteúdo
func ModifyFile(path string, repl func(string) string) {
	raw, err := ioutil.ReadFile(path)

	if err != nil {
		log.Print(err)
		return
	}

	content := repl(string(raw))

	ioutil.WriteFile(path, []byte(content), 0700)
}

// getprefscfg encontra o arquivo `prefs` baseado do seu sistema operacional
// e retorna uma ref `ini.file`
func GetPrefsCfg(spotifyPath string) (*ini.File, string, error) {
	var path string

	if runtime.GOOS == "windows" {
		path = filepath.Join(spotifyPath, "prefs")
	} else if runtime.GOOS == "linux" {
		path = filepath.Join(os.Getenv("HOME"), ".config", "spotify", "prefs")
	} else if runtime.GOOS == "darwin" {
		path = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Spotify", "prefs")
	} else {
		return nil, "", errors.New("Unsupported OS")
	}

	cfg, err := ini.Load(path)

	if err != nil {
		cfg = ini.Empty()
	}

	return cfg, path, nil
}

// função getspotifyversion
func GetSpotifyVersion(spotifyPath string) string {
	pref, _, err := GetPrefsCfg(spotifyPath)

	if err != nil {
		log.Fatal(err)
	}

	rootSection, err := pref.GetSection("")

	if err != nil {
		log.Fatal(err)
	}

	version := rootSection.Key("app.last-launched-version")

	return version.MustString("")
}

// getexecutabledir retorna o diretório do processo atual
func GetExecutableDir() string {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		exe, err := os.Executable()

		if err != nil {
			log.Fatal(err)
		}

		return filepath.Dir(exe)
	}

	return ""
}

// getjshelperdir retorna o diretório do jshelper no diretório executável
func GetJsHelperDir() string {
	return filepath.Join(GetExecutableDir(), "jsHelper")
}

// restartspotify
func RestartSpotify(spotifyPath string) {
	if runtime.GOOS == "windows" {
		exec.Command("taskkill", "/F", "/IM", "spotify.exe").Run()
		exec.Command(filepath.Join(spotifyPath, "spotify.exe")).Start()
	} else if runtime.GOOS == "linux" {
		exec.Command("pkill", "spotify").Run()
		exec.Command(filepath.Join(spotifyPath, "spotify")).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("pkill", "Spotify").Run()
		exec.Command("open", "/Applications/Spotify.app").Start()
	}
}