FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

ARG service_name

RUN go build -o ./output/app_bin ./cmd/${service_name}

EXPOSE 8080:8080/tcp

CMD [ "./output/app_bin" ]