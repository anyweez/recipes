FROM ubuntu:14.04
MAINTAINER Luke Segars <luke@lukesegars.com>

# Install git so that we can pull down the repo, which will have project-specific
# dependencies in it. Then we should execute those builds.
RUN apt-get update && apt-get install git --yes
RUN mkdir server; cd server; git clone https://github.com/luke-segars/recipes.git
RUN cd recipes; ls
RUN cd server/recipes; ./build
