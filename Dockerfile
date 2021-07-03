# golang:alpine
FROM golang:latest AS develop
ENV PROJECT_PATH=/meteostation

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH

# build backend
WORKDIR $PROJECT_PATH/backend
RUN go build

# build frontend
WORKDIR $PROJECT_PATH/ui/meteostation

RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install -y yarn
RUN yarn
RUN yarn build

WORKDIR $PROJECT_PATH

FROM alpine:latest AS production
ENV PROJECT_PATH=/meteostation

# copy distrib to target directory
COPY --from=develop $PROJECT_PATH/backend/meteostation $PROJECT_PATH/meteostation
COPY --from=develop $PROJECT_PATH/backend/storage/migrations/* $PROJECT_PATH/storage/migrations/
COPY --from=develop $PROJECT_PATH/backend/.meteostation.json /etc/.meteostation.json
COPY --from=develop $PROJECT_PATH/ui/meteostation/dist $PROJECT_PATH/www/static

# if you set USER postgres:postgres you will get error on start container fixing permissions on existing directory /postgres/data ... initdb: error: could not change permissions of directory "/postgres/data": Operation not permitted
# chmod: changing permissions of '/postgres/data': Operation not permitted
# USER postgres:postgres

WORKDIR $PROJECT_PATH
ENTRYPOINT ["/meteostation/meteostation"]





