FROM mcr.microsoft.com/windows/servercore:ltsc2019

ADD wincat.exe /wincat.exe
CMD cmd /c ping -t localhost
