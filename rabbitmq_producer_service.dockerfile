FROM alpine

WORKDIR /app

COPY producerServiceApp .

CMD [ "./producerServiceApp" ]