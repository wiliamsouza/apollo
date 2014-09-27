# Apollo API image
#
# Version 0.1.0

FROM ubuntu:14.04

MAINTAINER Wiliam Souza <wiliamsouza83@gmail.com>

# base
ENV LANG en_US.UTF-8
ENV DEBIAN_FRONTEND noninteractive

RUN locale-gen en_US en_US.UTF-8
RUN dpkg-reconfigure locales
RUN apt-get update

RUN apt-get install -y python-software-properties

# supervisor
RUN apt-get install supervisor -y
RUN update-rc.d -f supervisor disable

ADD etc/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# start script
ADD bin/startup /usr/local/bin/startup
RUN chmod +x /usr/local/bin/startup

CMD ["/usr/local/bin/startup"]

# environment

# dependencies
RUN apt-get install curl -y

# repos

# confd binary
RUN curl -L https://github.com/kelseyhightower/confd/releases/download/v0.5.0/confd-0.5.0-linux-amd64 -o /usr/local/bin/confd
RUN chmod +x /usr/local/bin/confd

# confd configuration
ADD confd /etc/confd

# Add apollod binary
ADD bin/apollod /usr/local/bin/apollod

EXPOSE 8000
