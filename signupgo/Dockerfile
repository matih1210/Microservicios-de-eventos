FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/server .
COPY .env .env
EXPOSE 3003
USER nonroot:nonroot
ENTRYPOINT ["./server"]
