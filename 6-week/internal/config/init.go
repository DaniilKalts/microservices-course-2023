package config

import "github.com/joho/godotenv"

func Load(path string) {
	godotenv.Load(path)
}
