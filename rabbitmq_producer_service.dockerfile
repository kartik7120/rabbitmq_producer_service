FROM alpine

WORKDIR /app

COPY producerServiceApp .

RUN chmod +x producerServiceApp

CMD [ "./producerServiceApp" ]