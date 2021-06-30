# golang:alpine
FROM golang:latest AS develop
ENV PROJECT_PATH=/meteostation

RUN mkdir -p $PROJECT_PATH
#COPY . $PROJECT_PATH

# build backend
WORKDIR $PROJECT_PATH/backend
#RUN go build

# build frontend
WORKDIR $PROJECT_PATH/ui/meteostation

#RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
#RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
#RUN apt-get update && apt-get install -y yarn
#RUN yarn
#RUN yarn build

WORKDIR $PROJECT_PATH

FROM postgres:latest AS production
ENV PROJECT_PATH=/meteostation

# copy result to target directory
#COPY --from=develop $PROJECT_PATH/backend/meteostation /usr/bin/meteostation
#COPY --from=develop $PROJECT_PATH/backend/.meteostation.json /etc/.meteostation.json
#COPY --from=develop $PROJECT_PATH/ui/meteostation/dist /www/static

#USER postgres:postgres
#ENTRYPOINT ["/usr/bin/meteostation"]





