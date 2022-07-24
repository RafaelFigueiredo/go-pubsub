FROM golang:alpine as build
WORKDIR /app
COPY . .

RUN go build -ldflags "-w" -o /app/bin/anymq

FROM golang:alpine as release
COPY --from=build /app/bin /


CMD ["/anymq"]