FROM scratch
LABEL maintainer="gideonhacer@gmail.com"
WORKDIR /
COPY app /
COPY certs /certs
COPY static /static
EXPOSE 80
ENTRYPOINT [ "/app" ]