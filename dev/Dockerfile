FROM golang:1.24-alpine
WORKDIR /app

RUN apk add --no-cache make gcc g++ curl nodejs npm

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest

COPY package.json package-lock.json ./
RUN npm install
RUN npm i tailwindcss @tailwindcss/cli @tailwindcss/forms @tailwindcss/typography @iconify/json @iconify/tailwind4

COPY dev/scripts/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
