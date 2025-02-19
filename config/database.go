package config

import (
 "database/sql"
 _ "github.com/go-sql-driver/mysql"
 _ "github.com/joho/godotenv/autoload"
 "github.com/rs/zerolog/log"
 "os"
 "time"
)

func Database() *sql.DB {
 user, exist := os.LookupEnv("MYSQL_USER")

 if !exist {
  log.Fatal().Msg("MYSQL_USER not set in .env")
 }

 pass, exist := os.LookupEnv("MYSQL_PASSWORD")

 if !exist {
  log.Fatal().Msg("MYSQL_PASSWORD not set in .env")
 }

 host, exist := os.LookupEnv("MYSQL_HOST")

 if !exist {
  log.Fatal().Msg("MYSQL_HOST not set in .env")
 }

	port, exist := os.LookupEnv("MYSQL_PORT")

 if !exist {
  log.Fatal().Msg("MYSQL_PORT not set in .env")
 }

 credentials := user + ":" + pass + "@(" + host + ":" + port + ")/?charset=utf8&parseTime=True"

 start := time.Now()
 database, err := sql.Open("mysql", credentials)

 if err != nil {
  log.Fatal().Err(err).
   Str("host", host).
   Str("port", port).
   Str("user", user).
   Msg("Failed to connect to database")
 }

 log.Info().
  Dur("duration", time.Since(start)).
  Str("host", host).
  Str("port", port).
  Msg("Database connection successful")

 _, err = database.Exec(`CREATE DATABASE gotodo`)

 if err != nil {
  log.Warn().Err(err).Msg("Failed to create database gotodo")
 }

 _, err = database.Exec(`USE gotodo`)

 if err != nil {
  log.Error().Err(err).Msg("Failed to use database gotodo")
 }

 _, err = database.Exec(`
  CREATE TABLE todos (
      id INT AUTO_INCREMENT,
      item TEXT NOT NULL,
      completed BOOLEAN DEFAULT FALSE,
      PRIMARY KEY (id)
  );
 `)

 if err != nil {
  log.Warn().Err(err).Msg("Failed to create todos table")
 }

 return database
}
