#!/bin/bash

sub_dir() {
	local cur_dir=$(dirname "$BASH_SOURCE")
	echo -n $(basename $cur_dir) 
}


