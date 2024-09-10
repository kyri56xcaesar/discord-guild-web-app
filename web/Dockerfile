FROM httpd:2.4
RUN rm -r /usr/local/apache2/htdocs/*
WORKDIR /usr/local/apache2/htdocs/
COPY /assets .
COPY /frontend .
EXPOSE 80