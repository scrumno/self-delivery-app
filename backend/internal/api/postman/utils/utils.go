package utils

import (
    "log"
    "fmt"
    "os"
    "github.com/Lexographics/go-postmangen"
)

type PostmanUtils struct {}

func (p *PostmanUtils) CreateDir(path string) error {
	err := os.MkdirAll(path, 0755)
    if err != nil {
        return fmt.Errorf("не удалось создать директорию %s: %w", path, err)
    }

	return nil;
}

func (p *PostmanUtils) Generate(path string, description string, baseUrl string, fileName string, endpoints []map[string]any) error {
	err := p.CreateDir(path)
	if err != nil {
		return fmt.Errorf("ошибка создания директории: %w", err)
	}

    pg := postmangen.NewPostmanGen(fmt.Sprintf("%s API", fileName), description)
    pg.AddVariable("base_url", baseUrl)
    collectionPath := fmt.Sprintf("%s/%v.postman_collection.json", path, fileName)

    err = p.register(pg, endpoints)
    if err != nil {
        return err
    }
   
    err = pg.WriteToFile(collectionPath)
    if err != nil {
        return err
    }

    return nil
}

func (p *PostmanUtils) register(pg *postmangen.PostmanGen, endpoints []map[string]any) error {
    for _, endpoint := range endpoints {
        err := pg.Register(endpoint)
        if err != nil {
            return fmt.Errorf("Ошибка в %s: %v", endpoint["path"], err)
        }
        log.Printf("Зарегистирован метод: %s", endpoint["path"])
    }

    return nil
}