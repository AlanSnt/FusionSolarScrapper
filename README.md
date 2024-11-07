# Fusion Solar scrapper

This project, written in GO, scrapes solar production information from the fusion solar website and sends it to MQTT.

## How to use

1. Install dependencies
```bash
go mod init AlanSnt/FusionSolarScrapper
go mod tidy
```

2. Configure environment variables
```bash
cp .env.example .env
```

3. Run the scrapper
```bash
go run main.go
```

## Deploy

1. Build docker image
```bash
docker build -t fusion-solar-scrapper .
```

2. Run docker container
```bash
docker run -d --name fusion-solar-scrapper fusion-solar-scrapper
```
