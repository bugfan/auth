FROM golang:alpine3.7 AS build-server
COPY . /usr/local/go/src/auth 

WORKDIR /usr/local/go/src/auth
RUN go get \
    && go build  # -ldflags "-s -w"

FROM alpine:3.7
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=build-server /usr/local/go/src/auth/auth . 
COPY --from=build-server /usr/local/go/src/auth/.env .  
COPY --from=build-server /usr/local/go/src/auth/conf ./conf

EXPOSE 5000

CMD ./auth