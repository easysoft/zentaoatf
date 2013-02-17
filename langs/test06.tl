#!/usr/bin/env tclsh
set case {
<<TC
title: tcl hello world.
expect:hello world.
TC
}

set hello "hello world."
puts [set hello]
