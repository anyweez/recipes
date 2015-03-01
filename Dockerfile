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
