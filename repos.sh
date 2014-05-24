#!/bin/bash

git remote add github git@github.com:wiliamsouza/apollo.git
git remote add bitbucket git@bitbucket.org:wiliamsouza/apollo.git
git remote add gitlab git@gitlab.com:wiliamsouza/apollo.git
git remote add gitorious git@gitorious.org:wiliamsouza/apollo.git
git remote set-url --push --add origin git@gitorious.org:wiliamsouza/apollo.git
git remote set-url --push --add origin git@bitbucket.org:wiliamsouza/apollo.git
git remote set-url --push --add origin git@gitlab.com:wiliamsouza/apollo.git
git remote set-url --push --add origin git@github.com:wiliamsouza/apollo.git
