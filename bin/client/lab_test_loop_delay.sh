#!/bin/bash

activate_loss () {
	sudo tc qdisc add dev enp3s0f0 root netem loss 1% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp3s0f0"
		exit 1
	fi
	sudo tc qdisc add dev enp3s0f1 root netem loss 3% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp3s0f1"
		exit 1
	fi
	echo "Added loss on all interfaces"
}

activate_delay () {
	sudo tc qdisc add dev enp3s0f0 root netem delay 50ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp3s0f0"
		exit 1
	fi
	sudo tc qdisc add dev enp3s0f1 root netem delay 300ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp3s0f1"
		exit 1
	fi
	echo "Activated delay on all interfaces"
}

deactivate_netem() {
	sudo tc qdisc del dev enp3s0f0 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem exited with non zero exit code for enp3s0f0"
	fi
	sudo tc qdisc del dev enp3s0f1 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem exited with non zero exit code for enp0s10"
	fi
	echo "Deactivated netem on all interfaces"
}

check_results() {
	sha256sum -c testfile2mb.sha
	if [ $? -ne 0 ]; then
		echo $1 " test failed with 2MB file" >> results
	else
		echo $1 " test passed with 2MB file" >> results
	fi
	sha256sum -c testfile10mb.sha
	if [ $? -ne 0 ]; then
		echo $1 " test failed with 10MB file" >> results
	else
		echo $1 " test passed with 10MB file" >> results
	fi
	sha256sum -c testfile100mb.sha
	if [ $? -ne 0 ]; then
		echo $1 " test failed with 100MB file" >> results
	else
		echo $1 " test passed with 100MB file" >> results
	fi
}

run_tests() {
	START2MB=`date +%s%N | cut -b1-13`
	cat testfile2mb | ./client -addr=10.0.1.10:4433 > testfile2mb.result
	FINISH2MB=`date +%s%N | cut -b1-13`
	wait
	START10MB=`date +%s%N | cut -b1-13`
	cat testfile10mb | ./client -addr=10.0.1.10:4433 > testfile10mb.result
	FINISH10MB=`date +%s%N | cut -b1-13`
	wait
	START100MB=`date +%s%N | cut -b1-13`
	cat testfile100mb | ./client -addr=10.0.1.10:4433 > testfile100mb.result
	FINISH100MB=`date +%s%N | cut -b1-13`
	wait
	DURATION2MB=$(( FINISH2MB - START2MB ))
	echo $DURATION2MB
	DURATION10MB=$(( FINISH10MB - START10MB ))
	echo $DURATION10MB
	DURATION100MB=$(( FINISH100MB - START100MB ))
	echo $DURATION100MB
	echo "2MB transfer took " $(( DURATION2MB / 1000 )) " seconds" >> bandwidth
	echo "10MB transfer took " $(( DURATION10MB / 1000 )) " seconds" >> bandwidth
	echo "100MB transfer took " $(( DURATION100MB / 1000 )) " seconds" >> bandwidth 
}

if [ -f flipped-delay.result ]; then
    rm flipped-delay.result
fi

if [ -f flipped-loss.result ]; then
    rm flipped-loss.result
fi

if [ -f clientLog ]; then
    rm clientLog
fi
if [ -f clientlog ]; then
    rm clientLog
fi


go build -o client main.go
if [ $? -ne 0 ]; then
	echo "Build failed, exiting"
	exit 1
fi

./bootstrap.sh
deactivate_netem
wait
./client -addr=10.0.1.10:4433
wait



echo "Running test with delay"
activate_delay
wait	
counter=0
while [ $counter -lt 100 ]; do

	run_tests
	check_results "Delay"
	counter=$counter+1
done
deactivate_netem
wait

cat results > `date '+%Y_%m_%d__%H_%M_%S'`_delay_results
cat bandwidth >> `date '+%Y_%m_%d__%H_%M_%S'`_delay_results
rm results
rm bandwidth

rm bandwidth
rm results
rm *.result
rm client
