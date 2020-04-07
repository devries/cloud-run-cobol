# Cloud Run with Cobol!

An simple web service in Cobol, because why not?

[![Run on Google
Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

This repository contains an invocation program written in Go which passes the
path as a command line argument into a cobol program. In the example program,
the path is appended after the word "HELLO". If the path is empty the service
returns "HELLO WORLD".

