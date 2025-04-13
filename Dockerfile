FROM golang:1.24.0
COPY . .
RUN go build -o server .
CMD ["./server"]