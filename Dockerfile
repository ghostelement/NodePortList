FROM  debian:buster
ENV CLUSTER_NAME=""
RUN mkdir -p /app/  && touch /app/k8s.yaml
ADD NodePortList /app
ADD templates /app/templates
WORKDIR /app
RUN chmod +x /app/NodePortList
ENTRYPOINT ["./NodePortList"]