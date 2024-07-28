FROM golang:latest
COPY . .
RUN ls
RUN go build -o main .
CMD ["./main"]