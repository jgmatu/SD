#!/bin/bash

go run barbery.go | grep 'Cliente' > clientes.out
go run barbery.go | grep 'Barbero' > barberos.out


