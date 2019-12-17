FROM ubuntu:19.04
ENV DEBIAN_FRONTEND noninteractive
USER root 

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y gnupg git postgresql-11 postgresql-contrib
RUN apt-get install curl -y

ENV GOVERSION 1.13.1
RUN curl -s -O https://dl.google.com/go/go$GOVERSION.linux-amd64.tar.gz
RUN tar -xzf go$GOVERSION.linux-amd64.tar.gz -C /usr/local
RUN chown -R root:root /usr/local/go
ENV GOPATH $HOME/work
ENV PATH $PATH:/usr/local/go/bin
ENV GOBIN $GOPATH/bin
RUN mkdir -p "$GOPATH/bin" "$GOPATH/src"
RUN GO11MODULE=on


USER postgres
ENV PGVERSION 11
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

USER root
RUN git clone https://github.com/ansushina/tech-db-forum.git

WORKDIR tech-db-forum
ARG CACHE_DATE=2015-01-10
RUN git pull

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql docker -a -f  database/create.sql &&\
    /etc/init.d/postgresql stop
RUN echo "local all all md5" > /etc/postgresql/$PGVERSION/main/pg_hba.conf &&\
    echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVERSION/main/pg_hba.conf
RUN cat database/postgresql.conf >> /etc/postgresql/$PGVERSION/main/postgresql.conf
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
EXPOSE 5432
USER root

RUN go get
RUN go build main.go
CMD ["/tech-db-forum/main"]
EXPOSE 5000