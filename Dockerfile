FROM alpine

WORKDIR /opt/bin

ENV CONFIG_FILE_PATH=/etc/urbs-setting/config.yml
COPY config/default.yml /etc/urbs-setting/config.yml

COPY ./dist/urbs-setting .

CMD ["./urbs-setting"]
