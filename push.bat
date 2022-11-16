@echo off
set commit=%1
git add --all
git commit -m %commit%
git push -f -u origin HEAD:main
@echo on