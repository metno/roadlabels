# FIRST STAGE:  build the app.
FROM registry.met.no/baseimg/ubuntu:22.04 AS build-app
#FROM registry.met.no/baseimg/ubuntu:22.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get -y update && \
    apt-get -y upgrade

RUN apt-get -y install apt-utils
RUN apt-get -y install libopencv-dev wget  make git gcc g++ 

RUN apt-get -y install  ca-certificates git-core ssh


WORKDIR /usr/local
#RUN wget https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
RUN wget https://go.dev/dl/go1.19.7.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.7.linux-amd64.tar.gz

ENV PATH="${PATH}:/usr/local/go/bin"



#ARG LOCALHOME
#ENV LOCALHOME $LOCALHOME
#RUN mkdir -p /root/.ssh

#ENV GOPRIVATE gitlab.met.no
#ADD /$LOCALHOME/.ssh/id_rsa /root/.ssh/id_rsa
#RUN chmod 700 /root/.ssh/id_rsa
#RUN echo "Host gitlab.met.no\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
#RUN git config --global url.ssh://git@gitlab.met.no/.insteadOf https://gitlab.met.no/
#ADD $LOCALHOME/.netrc /root/.netrc

WORKDIR /build/app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# Dependencies are downloaded only when go.mod or go.sum changes.
RUN go mod download

# Copy the rest of the source files.
COPY . .

RUN go mod tidy
ARG S3SecretKey
ARG S3AccessKey


ENV S3SecretKey $S3SecretKey
ENV S3AccessKey $S3AccessKey



RUN make

# Second stage. Install 
FROM registry.met.no/baseimg/ubuntu:22.04
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get -y update && \
    apt-get -y upgrade

RUN apt-get -y install apt-utils
RUN apt-get -y install curl
RUN apt-get -y install telnet
RUN apt-get -y install libopencv-dev
RUN apt-get -y install python3-netcdf4
RUN apt-get -y install python3-pyproj
RUN apt-get -y install sqlite3

WORKDIR /app
COPY --from=build-app /build/app/roadlabels /app/
COPY --from=build-app /build/app/exttools/ncvars.py /app/
ENV PYTHONPATH=/app/
ENV HDF5_USE_FILE_LOCKING=FALSE
RUN mkdir -p /lustre
ADD entry.sh /usr/local/bin/entry.sh
ADD var/lib/roadlabels/roadcams.db /roadlabels/roadcams.db
ADD var/lib/roadlabels/userdb-empty.db /roadlabels/userdb-empty.db
RUN mkdir -p /var/lib/roadlabels
RUN chown -R nobody.nogroup /var/lib/roadlabels
USER nobody:nogroup

EXPOSE 25260

ENTRYPOINT /usr/local/bin/entry.sh
#CMD /usr/local/bin/entry.sh
