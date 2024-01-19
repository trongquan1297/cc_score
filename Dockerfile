# Sử dụng một hình ảnh cơ sở với Go runtime
FROM golang:1.16-alpine as builder

# Set thư mục làm việc trong container
WORKDIR /app

# Sao chép mã nguồn vào container
COPY . .

# Tạo binary của ứng dụng Go
RUN go build -o app

# Sử dụng hình ảnh nhỏ để triển khai ứng dụng
FROM alpine:latest

# Set thư mục làm việc trong container
WORKDIR /app

# Sao chép binary từ builder stage
COPY --from=builder /app/app /app/app
COPY config.yml .

# Expose cổng mà ứng dụng lắng nghe
EXPOSE 8080

# Chạy ứng dụng khi container được khởi động
CMD ["/app/app"]
