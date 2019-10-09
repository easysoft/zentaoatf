#!/usr/bin/env bash

rm -rf log/*
ztf/ztf run demo/lang demo/sample
ztf/ztf run demo/ztf/demo_test.sh