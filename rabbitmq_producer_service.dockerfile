FROM alpine

WORKDIR /app

COPY producerServiceApp .

ARG RABBITMQ_URL
ENV RABBITMQ_URL=${RABBITMQ_URL}

RUN chmod +x producerServiceApp

CMD [ "./producerServiceApp" ]