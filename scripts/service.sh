#!/bin/bash

# $1 == [ start | stop | restart | enable | disable ]
sudo systemctl $1 gopanel
