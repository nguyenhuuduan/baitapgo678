# Bước 1: Sử dụng hình ảnh Go chính thức để build ứng dụng
# Bước 1: Sử dụng hình ảnh Go chính thức để build ứng dụng
FROM golang:1.23-alpine AS builder

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép toàn bộ mã nguồn vào container
COPY . .

# Tải về các module phụ thuộc
RUN go mod download

# Build ứng dụng Go, tạo ra file thực thi `main`
RUN go build -o main .

# Bước 2: Tạo một hình ảnh nhẹ chỉ để chạy ứng dụng
FROM alpine:3.16

# Sao chép script wait-for.sh trước khi thay đổi quyền
COPY wait-for.sh /wait-for.sh

# Cài đặt các thư viện cần thiết và thiết lập quyền thực thi cho `wait-for.sh`
RUN apk --no-cache add ca-certificates bash && \
    chmod +x /wait-for.sh

# Thiết lập thư mục làm việc
WORKDIR /root/

# Sao chép file thực thi từ container builder sang container hiện tại
COPY --from=builder /app/main .

# Khởi động ứng dụng sau khi chờ dịch vụ db sẵn sàng
CMD ["/wait-for.sh", "db", "./main"]

# Mở cổng 8080
EXPOSE 8080



