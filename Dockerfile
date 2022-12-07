FROM centos
ENV MY_SERVICE_PORT=80
LABEL go="1.14" version="1.0"
ADD bin/amd64/gosearch /gosearch
ADD configs/ /configs/
ADD html/ /html/

EXPOSE 80
ENTRYPOINT /gosearch
