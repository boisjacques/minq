#!/bin/bash

activate_loss () {
	sudo tc qdisc add dev enp3s0f0 root netem loss 5% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp3s0f0"
		exit 1
	fi
	sudo tc qdisc add dev enp3s0f1 root netem loss 2% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp3s0f1"
		exit 1
	fi
	echo "Added loss on all interfaces"
}

activate_reordering () {
	sudo tc qdisc add dev enp3s0f0 root netem delay 100ms 75ms
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp3s0f0"
		exit 1
	fi
	sudo tc qdisc add dev enp3s0f1 root netem delay 45ms 100ms
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp3s0f1"
		exit 1
	fi
	echo "Added reordering on all interfaces"

}

activate_delay () {
	sudo tc qdisc add dev enp3s0f0 root netem delay 100ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp3s0f0"
		exit 1
	fi
	sudo tc qdisc add dev enp3s0f1 root netem delay 150ms 10ms 25%
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

echo "Running test without wire errors"
cat testfile2mb | ./client -addr=10.0.1.10:4433 > testfile2mb.result
wait
cat testfile10mb | ./client -addr=10.0.1.10:4433 > testfile10mb.result
wait
cat testfile100mb | ./client -addr=10.0.1.10:4433 > testfile100mb.result
wait
check_results "Plain"

echo "Running test with delay"
activate_delay
wait
cat testfile2mb | ./client -addr=10.0.1.10:4433 > testfile2mb.result
wait
cat testfile10mb | ./client -addr=10.0.1.10:4433 > testfile10mb.result
wait
cat testfile100mb | ./client -addr=10.0.1.10:4433 > testfile100mb.result
wait
check_results "Delay"
deactivate_netem
wait

echo "Running test with loss"
activate_loss
cat testfile2mb | ./client -addr=10.0.1.10:4433 > testfile2mb.result
wait
cat testfile10mb | ./client -addr=10.0.1.10:4433 > testfile10mb.result
wait
cat testfile100mb | ./client -addr=10.0.1.10:4433 > testfile100mb.result
wait
check_results "Loss"
deactivate_netem
wait

cat results > `date`_results


rm results
rm *.result
rm client

