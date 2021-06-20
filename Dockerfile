# golang:alpine
FROM golang:latest AS develop
ENV PROJECT_PATH=/meteostation

# build backend
# build frontend
# copy result to target directory

RUN mkdir -p $PROJECT_PATH
COPY . $PROJECT_PATH
WORKDIR $PROJECT_PATH
#RUN make


FROM postgres:latest AS production 



#RUN apk --no-cache add ca-certificates
#COPY --from=develop /application-server/build/application-server /usr/bin/application-server
#COPY --from=develop /application-server/configuration/application-server/application-server.toml /etc/application-server/application-server.toml

#USER nobody:nogroup
#ENTRYPOINT ["/usr/bin/meteostation-service"]





