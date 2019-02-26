#!/bin/sh

_pci_serv=http://localhost:8079
_listen_addr=:8080
_self_url="http://127.0.0.1:8080"

rm ./senren2
go build
LISTEN_ADDR=$_listen_addr PCI_SERV=$_pci_serv SENREN_SELF_URL=$_self_url ./senren2
