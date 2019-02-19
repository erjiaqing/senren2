#!/bin/sh

_git_serv=http://localhost:22222
_git_user=pci2
_git_pass=pcipcipci
_listen_addr=:8079
_editor_serv=http://localhost:8078

rm ./pci2
go build
PCI_GIT_USER=$_git_user PCI_GIT_PASS=$_git_pass PCI_GIT_SERV=$_git_serv LISTEN_ADDR=$_listen_addr PCI_EDITOR_SERV=$_editor_serv ./pci2