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
COPY cloud /app/cloud
RUN chmod +x /app/cloud

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./cloud", "-c", "configs/cloud.yml"]
# if you use the Configuration Center, cloud.yml is changed to the Configuration Center configuration.
#CMD ["./cloud", "-c", "configs/cloud.yml", "-enable-cc"]
