FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# set the time zone to Shanghai
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY frame /app/frame
RUN chmod +x /app/frame

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./frame", "-c", "configs/frame.yml"]
# if you use the Configuration Center, frame.yml is changed to the Configuration Center configuration.
#CMD ["./frame", "-c", "configs/frame.yml", "-enable-cc"]
