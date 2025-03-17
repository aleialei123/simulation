#!/bin/bash
rm plot.dat

echo 0 0 >> plot.dat
writedata() {
    for i in {1..10}; do
        echo -e $i"\t"$((i*i)) >> plot.dat
        sleep 1
        # echo "done"
    done
}

writedata &
sleep 1
gnuplot liveplot.gnu