package gobbler

import (
	"encoding/gob"
	"path/filepath"
	"path"
	"fmt"
	"github.com/lokken/family-football/types"
	"log"
	"os"
	"runtime"
	"io"
)

func vault() string {
	_, filePath, _, _ := runtime.Caller(0)
	p := filepath.Dir(filePath)
	vaultPath := path.Join(p, "vault")
	err := os.MkdirAll(vaultPath, os.ModePerm)
	if err != nil {
		log.Fatal("mkdir all error: ", err)
	}
	return vaultPath
}

func SaveSchedule(week int, data io.Reader) {
	p := vault()
	filename := path.Join(p, fmt.Sprintf("week%d-raw.html", week))

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(file, data)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadSchedule(week int) io.Reader {
	p := vault()
	filename := path.Join(p, fmt.Sprintf("week%d-raw.html", week))

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func SaveGames(games []types.Game) {
	p := vault()
	filename := path.Join(p, "games.gob")

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(games)
	if err != nil {
		log.Fatal("encode error: ", err)
	}
}

func LoadGames(games *[]types.Game) {
	p := vault()
	filename := path.Join(p, "games.gob")

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	enc := gob.NewDecoder(file)
	err = enc.Decode(games)
	if err != nil {
		log.Fatal("encode error: ", err)
	}
}
