FROM alpine:latest
WORKDIR /opt/restaurant
RUN mkdir /opt/restaurant/conf
COPY main /opt/restaurant
COPY start.sh /opt/restaurant
ENTRYPOINT ["sh","/opt/restaurant/start.sh"]