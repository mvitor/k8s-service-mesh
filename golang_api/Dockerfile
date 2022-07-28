FROM golang:1.19rc2-alpine3.16 as dev

WORKDIR /work

##
## Build
##

FROM golang:1.19rc2-alpine3.16 as build

WORKDIR /api
COPY ./api/* /api/
RUN go build -o api

##
## Deploy build to image
##

FROM alpine as runtime 
COPY --from=build /api/api /
CMD ./api

EXPOSE 8080