package config

import (
	"log"

	"github.com/go-playground/validator/v10"
)

type Tools struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	Validator   *validator.Validate
}
