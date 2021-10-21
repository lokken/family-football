package gobbler

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/lokken/family-football/types"
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

func PutBonuses(bonusi map[string]*types.Bonus) {
	p := vault()
	filename := path.Join(p, "bonuses.gob")

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	var bonuses map[string]*types.Bonus
	enc := gob.NewDecoder(f)
	err = enc.Decode(&bonuses)
	if err != nil {
		log.Fatal("encode error: ", err)
	}

	for kbi, bi := range bonusi {
		if bonuses[kbi] == nil {
			bonuses[kbi] = bi
		} else {
			bonuses[kbi].Type = bi.Type
			bonuses[kbi].Qualifier = bi.Qualifier
			bonuses[kbi].Quantifier = bi.Quantifier
		}
	}
	SaveBonuses(bonuses)
}

func SaveBonuses(bonuses map[string]*types.Bonus) {
	p := vault()
	filename := path.Join(p, "bonuses.gob")

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(bonuses)
	if err != nil {
		log.Fatal("encode error: ", err)
	}
}

func LoadBonuses(bonuses *map[string]*types.Bonus) {
	p := vault()
	filename := path.Join(p, "bonuses.gob")

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	dec := gob.NewDecoder(file)
	err = dec.Decode(bonuses)
	if err != nil {
		log.Fatal("decode error: ", err)
	}
}
