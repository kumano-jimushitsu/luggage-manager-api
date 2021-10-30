FROM mcr.microsoft.com/mssql/server:latest

USER root
SHELL ["/bin/bash", "-c"]

WORKDIR /database
COPY ./database .
RUN chmod -R +x /database/db-init.sh

EXPOSE 1433

ENTRYPOINT ["/bin/bash", "/database/entrypoint.sh"]