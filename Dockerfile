FROM ubuntu:14.04
MAINTAINER Luke Segars <luke@lukesegars.com>

# Install every version control system known to man so that we can pull down the repo
# and install dependencies.
RUN apt-get update && apt-get install git golang mercurial bzr --yes
RUN mkdir server; cd server; git clone https://github.com/luke-segars/recipes.git
RUN cd server/recipes; ./build
