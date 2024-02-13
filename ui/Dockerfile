FROM node:17-alpine

# install dependencies
WORKDIR /app/ui
COPY ui/package.json /app/ui
COPY ui/package-lock.json /app/ui
RUN npm ci

# Copy all local files into the image.
COPY ui /app/ui
WORKDIR /app/ui

RUN npx svelte-kit build

###
# Only copy over the Node pieces we need
# ~> Saves 35MB
###
FROM node:17-alpine

WORKDIR /app
COPY --from=0 /app/ui .
COPY . .

EXPOSE 3000
CMD ["node", "./build"]
