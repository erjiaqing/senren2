#!/bin/sh

_git_serv=http://localhost:22222
_git_user=pci2
_git_pass=pcipcipci
_pci_serv=http://localhost:8079
_pci_session="5c6fa882rbhL3fWDnACIPK-vN--_0cwS"

rm ./judgedaemon
go build -o judgedaemon
PCI_GIT_USER=$_git_user PCI_GIT_PASS=$_git_pass PCI_GIT_SERV=$_git_serv PCI_SERV=$_pci_serv PCI_SESSION=$_pci_session ./judgedaemon