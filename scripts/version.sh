#!/bin/sh

if [ $# -eq 0 ]; then
    echo "No arguments supplied"
    exit 1
fi

IN=$1

branch=$(git rev-parse --abbrev-ref HEAD)
branch=${branch/\//\\\/}
short_ver=`git rev-parse --short HEAD`
sem_ver=`cat ${IN} | sed -n 's/\(v[0-9]*\.[0-9]*\.[0-9]*\).*/\1/p'`
patch=`echo ${sem_ver} | sed -n 's/.*[0-9][0-9]*\.[0-9][0-9]*\.\([0-9][0-9]*\)/\1/p'`
patch=`expr ${patch} + 1`
new_ver=`echo ${sem_ver} | sed -n "s/.*\([0-9][0-9]*\.[0-9][0-9]*\.\)\([0-9][0-9]*\)/\1$patch/p"`

function replace() {
    sed -i "" "s/[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/$new_ver/g" $1
    sed -i "" "s#-- .*([^\)]*)#-- ${branch}\($short_ver\)#g" $1
}

for var in "$@"; do
    echo "Update ${var} to $new_ver -- $branch($short_ver)"
    replace ${var}
done
