FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o city_population_app .

EXPOSE 8080

CMD ["./city_population_app"]
