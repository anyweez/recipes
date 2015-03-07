# The command I use to run this is:
# 
# 	sudo docker run -d -i -t -p 13033:13033 --name frontend lukes/dinder:latest /bin/bash
#
# This runs the instance as a background daemon and maps the docker instance's port 13033 (frontend 
# server port) to the local machines port 13033. The docker instance is also explicitly named "frontend",
# so if more than one server need to be run simultaneously then you'd need to remove the explicitly mapping
# for both the port and instance name and then handle the dynamic port selection on the local machine.
FROM ubuntu:14.04
MAINTAINER Luke Segars <luke@lukesegars.com>

# Install every version control system known to man so that we can pull down the repo
# and install dependencies.
RUN apt-get update && apt-get install git golang mercurial bzr mongodb --yes
RUN mkdir /server; cd /server; git clone https://github.com/luke-segars/recipes.git
RUN cd /server/recipes; ./build
RUN echo "service mongodb start" >> /etc/bash.bashrc
# The mongo service start call returns before mongo is actually activated (the daemon
# is running but it won't allow you to connect to the database until it's built; which
# takes some time @ first launch).
RUN echo "until mongo; do sleep 1; done" >> /etc/bash.bashrc
RUN echo "cd /server/recipes; ./retrieve &" >> /etc/bash.bashrc
RUN echo "cd /server/recipes; ./frontend &" >> /etc/bash.bashrc
