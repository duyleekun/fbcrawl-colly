FROM golang:1.14-alpine

RUN mkdir -p /app
WORKDIR /app
ADD ./go.mod /app
ADD ./go.sum /app
ADD ./ /app
RUN go get


ENV PORT 3000
RUN build -o server qnetwork.net/fbcrawl #gosetup
ENTRYPOINT ["./server"]
