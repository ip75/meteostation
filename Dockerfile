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
RUN apt-get update && \
    apt-get install yarn && \
    yarn && \
    yarn build

WORKDIR $PROJECT_PATH

FROM postgres:latest AS production

# copy result to target directory
RUN apk --no-cache add ca-certificates
COPY --from=develop $PROJECT_PATH/backend/meteostation /usr/bin/meteostation
COPY --from=develop $PROJECT_PATH/backend/.meteostation.json /etc/.meteostation.json
COPY --from=develop $PROJECT_PATH/ui/meteostation/dist /www/static

USER nobody:nogroup
ENTRYPOINT ["/usr/bin/meteostation"]





