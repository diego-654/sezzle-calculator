# ---- Stage 1: build the static site ----
FROM node:22-alpine AS build

WORKDIR /app

# Install dependencies first so Docker caches this layer
# until package.json or package-lock.json change.
COPY package.json package-lock.json ./
RUN npm ci

# Copy the source (filtered by .dockerignore) and build it into /app/dist.
COPY . .
RUN npm run build

# ---- Stage 2: serve the files with nginx ----
FROM nginxinc/nginx-unprivileged:1.27-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

# The base image already runs as a non-root user and starts nginx.
EXPOSE 8080
