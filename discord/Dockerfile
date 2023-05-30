FROM node:18-alpine as BUILDER

WORKDIR /build

COPY package-lock.json package.json ./

RUN npm ci

COPY . .

RUN npm run build

RUN rm -rf ./node_modules

RUN npm ci --omit=dev

FROM node:18-alpine
WORKDIR /app

COPY --from=BUILDER /build/dist ./dist
COPY --from=BUILDER /build/node_modules ./node_modules

CMD ["node", "./dist/main.js"]