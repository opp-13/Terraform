#!/bin/bash
compgen -c | grep -E '^.{4}$' > result.txt
