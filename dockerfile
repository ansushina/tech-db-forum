FROM ubuntu:19.04
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y gnupg git postgresql-11 postgresql-contrib


USER root 

RUN git clone https://github.com/ansushina/tech-db-forum.git

WORKDIR tech-db-forum

RUN git pull
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql docker -a -f  database/create.sql &&\
    /etc/init.d/postgresql stop

RUN cat database/create.sql

USER root
ENV GOVERSION 1.13.1
RUN curl -s -O https://dl.google.com/go/go$GOVERSION.linux-amd64.tar.gz
RUN tar -xzf go$GOVERSION.linux-amd64.tar.gz -C /usr/local
RUN chown -R root:root /usr/local/go
ENV GOPATH $HOME/work
ENV PATH $PATH:/usr/local/go/bin
ENV GOBIN $GOPATH/bin
RUN mkdir -p "$GOPATH/bin" "$GOPATH/src"
RUN GO11MODULE=on

RUN go get
RUN go build main.go
CMD ["/tech-db-forum/main"]
EXPOSE 5000