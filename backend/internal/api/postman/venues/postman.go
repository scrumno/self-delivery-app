package venues

import (
    "log"
    "reflect"
    "fmt"
    "github.com/scrumno/scrumno-api/internal/api/postman/venues/params"
    "github.com/scrumno/scrumno-api/internal/api/postman/utils"
)

type Generator struct {
    BaseUrl string
    CollectionFileName string 
    CollectionOutputPath string
    CollectionDescription string
    CollectionEndpoints []map[string]any
}

func NewGenerator() *Generator {
    return &Generator{
        BaseUrl: "http://localhost:8081/api/v1",
        CollectionFileName: "venues",
        CollectionOutputPath: "collections",
        CollectionDescription: "Сгенерированная коллекция API запросов для самовывоза",
        CollectionEndpoints: []map[string]any{
            {
                "method":    "GET",
                "path":      "/users",
                "inputType": reflect.TypeOf(params.GetUserRequest{}),
            },
            {
                "method":    "POST",
                "path":      "/users",
                "inputType": reflect.TypeOf(params.CreateUserRequest{}),
            },
        },
    }
}

func (g *Generator) Name() string {
    return g.CollectionFileName
}

func (g *Generator) Generate() error {
    log.Printf("Начинаем генерацию коллекции %s", g.CollectionFileName)
    
    utils := &utils.PostmanUtils{}
    
    err := utils.Generate(
        g.CollectionOutputPath,
        g.CollectionDescription,
        g.BaseUrl,
        g.CollectionFileName,
        g.CollectionEndpoints,
    )
    
    if err != nil {
        return fmt.Errorf("ошибка генерации коллекции %s: %w", g.CollectionFileName, err)
    }
    
    log.Printf("Коллекция %s успешно сгенерирована", g.CollectionFileName)
    return nil
}