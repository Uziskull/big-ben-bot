FROM golang:1.23 AS builder

WORKDIR /app-build

# install ffmpeg and dca, to convert audio files
RUN apt-get update && apt-get install ffmpeg -y

COPY *.mp3 .
RUN go install github.com/bwmarrin/dca/cmd/dca@latest
RUN ffmpeg -i initBell.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > initBell.dca
RUN ffmpeg -i mainBell.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > mainBell.dca

# build bot
COPY . .
RUN go mod tidy
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags timetzdata -o bot .

# pass it to new image
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app-build/bot /app-build/*.dca .

CMD [ "./bot" ]
