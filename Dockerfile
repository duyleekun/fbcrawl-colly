FROM golang:1.17-alpine
RUN apk add --no-cache git build-base tzdata

RUN mkdir -p /app
WORKDIR /app
ADD ./go.mod /app
ADD ./go.sum /app
ADD ./ /app
RUN go get


ENV PORT 3000
RUN go build -o server qnetwork.net/fbcrawl
ENTRYPOINT ["./server"]
