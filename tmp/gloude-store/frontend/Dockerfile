# Этап сборки
FROM node:20-alpine AS builder

WORKDIR /app

# Аргументы для переменных окружения Vite
ARG VITE_API_URL=/api/v1
ENV VITE_API_URL=$VITE_API_URL

# Копируем package.json и устанавливаем зависимости
COPY package.json ./
RUN npm install --legacy-peer-deps

# Копируем исходный код и собираем
COPY . .
RUN npm run build

# Финальный образ — nginx для раздачи статики
FROM nginx:alpine

# Копируем собранный фронтенд
COPY --from=builder /app/dist /usr/share/nginx/html

# Конфиг nginx для SPA (все пути ведут к index.html)
RUN cat > /etc/nginx/conf.d/default.conf << 'EOF'
server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip сжатие
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # Кеширование статики
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, no-transform";
    }

    # SPA fallback — все несуществующие пути возвращают index.html
    location / {
        try_files $uri $uri/ /index.html;
    }
}
EOF

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
