#!/bin/bash

shopt -s expand_aliases

antlr_jar=antlr-4.13.0-complete.jar

if [ ! -f ${antlr_jar} ]; then
    wget "https://www.antlr.org/download/${antlr_jar}"
fi

alias antlr4="java -Xmx500M -cp './$antlr_jar:\$CLASSPATH' org.antlr.v4.Tool"
antlr4 -Dlanguage=Go -no-listener -visitor -package grammar *.g4
