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

# supervisor
RUN apt-get install supervisor -y
RUN update-rc.d -f supervisor disable

ADD etc/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# start script
ADD bin/startup /usr/local/bin/startup
RUN chmod +x /usr/local/bin/startup

CMD ["/usr/local/bin/startup"]

# install etcdctl binary
ADD bin/etcdctl /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

# Add apollod binary
ADD bin/apollod /usr/local/bin/apollod

EXPOSE 8000
