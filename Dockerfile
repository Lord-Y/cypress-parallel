ARG DOCKER_IMAGE_NODEJS="node:14-alpine"
FROM ${DOCKER_IMAGE_NODEJS} as builder

RUN mkdir /app
WORKDIR /app
COPY . .
RUN adduser -D -g '' appuser
RUN rm -rf .git .env*
RUN apk update && apk upgrade
RUN apk add --no-cache git curl && rm -rf /var/cache/apk/*
RUN npm install
RUN npm run build

ARG DOCKER_IMAGE_NODEJS="node:14-alpine"
FROM ${DOCKER_IMAGE_NODEJS}
RUN apk update && apk upgrade
RUN apk add --no-cache git curl && rm -rf /var/cache/apk/*
RUN mkdir /app /home/appuser
WORKDIR /app
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/dist ./dist
RUN npm install vite @vitejs/plugin-vue
ENV NODE_ENV=production
RUN npm install --only=production

USER appuser
EXPOSE 8080
CMD ./node_modules/.bin/vite preview --port 8080
