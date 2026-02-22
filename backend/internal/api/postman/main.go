package main

import (
    "log"
    "github.com/scrumno/scrumno-api/internal/api/postman/user"
    "github.com/scrumno/scrumno-api/internal/api/postman/health"
    "github.com/scrumno/scrumno-api/internal/api/postman/venues"
)

type CollectionGenerator interface {
    Generate() error
    Name() string
}

func main() {
    log.Println("Генерируем коллекции")

    generators := []CollectionGenerator{
        health.NewGenerator(),
        user.NewGenerator(),
        venues.NewGenerator(),
    }

    for _, gen := range generators {
        log.Printf("Генерация %s...", gen.Name())
        err := gen.Generate() 
        if err != nil {
            log.Fatalf("Ошибка в %s: %v", gen.Name(), err)
        }
    }
    
    log.Println("Все коллекции успешно созданы")
}