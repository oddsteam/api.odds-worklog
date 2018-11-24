FROM alpine

COPY run.sh /root
RUN apt-get update -y && apt-get install -y curl cron
RUN set -x && \
    chmod +x /root/run.sh && \
    apk update && \
    apk upgrade && \
    apk add --no-cache mongodb && \
    rm /usr/bin/mongoperf && \
    rm -rf /var/cache/apk/* 
    
VOLUME /data/db
EXPOSE 27017 28017

ENTRYPOINT [ "/root/run.sh" ]
CMD ["mongod"]