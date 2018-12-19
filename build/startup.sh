#!/usr/bin/env bash
touch app.log
nohup ./input_tracker_app >> app.log 2>&1 & echo $!
tail -qF app.log