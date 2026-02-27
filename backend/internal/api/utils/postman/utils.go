package utils

import (
    "log"
    "fmt"
    "os"
    "github.com/Lexographics/go-postmangen"
)

func createDir(path string) error {
	err := os.MkdirAll(path, 0755)
    if err != nil {
        return err
    }

	return nil;
}

func generateCollection(path string, description string, baseUrl string, fileName string, endpoints []map[string]any) error {
	err := createDir(path)
	if err != nil {
		return fmt.Errorf("Ошибка создания директории: %w", err)
	}

    pg := postmangen.NewPostmanGen(fmt.Sprintf("%s API", fileName), description)
	
    pg.AddVariable("base_url", baseUrl)

    collectionPath := fmt.Sprintf("%s/%v.postman_collection.json", path, fileName)

    err = registerEndpoints(pg, endpoints)
    if err != nil {
        return err
    }
   
    err = pg.WriteToFile(collectionPath)
    if err != nil {
        return err
    }

    return nil
}

func registerEndpoints(pg *postmangen.PostmanGen, endpoints []map[string]any) error {
    for _, endpoint := range endpoints {
        err := pg.Register(endpoint)
        if err != nil {
            return fmt.Errorf("Ошибка в %s: %v", endpoint["path"], err)
        }
        log.Printf("Зарегистирован метод: %s", endpoint["path"])
    }

    return nil
}

func Generate(breakpointGroups map[string][]map[string]interface{}) error {
    log.Println("Генерируем коллекции")
    
    for group, breakpoints := range breakpointGroups { 
        err := generateCollection(
            os.Getenv("BASE_PATH"),
            "Сгенерированная коллекция API запросов для самовывоза",
            os.Getenv("BASE_API_URL"),
            group,
            breakpoints,
        )

        if err != nil {
            return fmt.Errorf("Ошибка генерации коллекции %s: %w", group, err)
        }
    }

    log.Println("Все коллекции успешно созданы")

    return nil
}