# 2023 FMNX team.
# Use of this code is governed by GNU General Public License.
# Official web page: https://fmnx.su/core/pack
# Contact email: help@fmnx.su

FROM docker.io/golang:latest as build

WORKDIR /src

COPY go.mod /src
COPY go.sum /src
RUN go mod download

COPY . /src/
RUN go build -o registry .

FROM docker.io/library/alpine:3.18

LABEL description="Standalone registry for arch packages"
LABEL maintainer="dancheg97 <dancheg97@fmnx.su>"

WORKDIR /home

COPY --from=build /src/registry /usr/bin/registry

ENTRYPOINT ["registry"]
CMD ["serve"]