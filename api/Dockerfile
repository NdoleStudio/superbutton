FROM golang as builder

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT

WORKDIR /superbutton

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go get github.com/swaggo/swag/gen@latest
RUN go get github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/swaggo/swag/cmd/swag
RUN swag init
RUN go install github.com/gobuffalo/packr/v2/packr2@latest
RUN packr2
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$GIT_COMMIT" -o /bin/superbutton .

FROM alpine:latest

RUN addgroup -S superbutton && adduser -S superbutton -G superbutton

USER superbutton
WORKDIR /home/superbutton

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip
COPY --from=builder /bin/superbutton ./

ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8000

ENTRYPOINT ["./superbutton", "--dotenv=false"]
