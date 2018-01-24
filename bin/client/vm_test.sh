#!/bin/bash

activate_loss () {
	tc qdisc add dev enp0s3 root netem loss 5% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp0s3"
		exit 1
	fi
	tc qdisc add dev enp0s8 root netem loss 2% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp0s8"
		exit 1
	fi
	tc qdisc add dev enp0s9 root netem loss 1% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp0s9"
		exit 1
	fi

	tc qdisc add dev enp0s10 root netem loss 7% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp0s10"
		exit 1
	fi

	echo "Added loss on all interfaces"
}

activate_reordering () {
	tc qdisc add dev enp0s3 root netem reorder 15% 25%
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp0s3"
		exit 1
	fi
	tc qdisc add dev enp0s8 root netem reorder 20% 15%
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp0s8"
		exit 1
	fi
	tc qdisc add dev enp0s9 root netem reorder 35% 5%
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp0s9"
		exit 1
	fi
	tc qdisc add dev enp0s10 root netem reorder 5% 45%
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp0s10"
		exit 1
	fi

	echo "Added reordering on all interfaces"
}

activate_delay () {
	tc qdisc add dev enp0s3 root netem delay 100ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp0s3"
		exit 1
	fi
	tc qdisc add dev enp0s8 root netem delay 150ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp0s8"
		exit 1
	fi
	tc qdisc add dev enp0s9 root netem delay 175ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp0s9"
		exit 1
	fi
	tc qdisc add dev enp0s10 root netem delay 200ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp0s10"
		exit 1
	fi

	echo "Activated delay on all interfaces"
}

deactivate_netem() {
	tc qdisc del dev enp0s3 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem exited with non zero exit code for enp0s3"
	fi
	tc qdisc del dev enp0s8 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem exited with non zero exit code for enp0s8"
	fi
	tc qdisc del dev enp0s9 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem exited with non zero exit code for enp0s9"
	fi
	tc qdisc del dev enp0s10 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem exited with non zero exit code for enp0s10"
	fi
}

check_results() {
	sha256sum -c testfile2mb.sha
	if [ $? -ne 0]; then
		echo $1 " test failed with 2MB file" >> results
	else
		echo $1 " test passed with 2MB file" >> results
	fi
	sha256sum -c testfile10mb.sha
	if [ $? -ne 0]; then
		echo $1 " test failed with 10MB file" >> results
	else
		echo $1 " test passed with 10MB file" >> results
	fi
	sha256sum -c testfile100mb.sha
	if [ $? -ne 0]; then
		echo $1 " test failed with 100MB file" >> results
	else
		echo $1 " test passed with 100MB file" >> results
	fi
}


if [ "$(whoami)" != "root" ]; then
	echo "Sorry, you are not root."
	exit 1
fi

go build -o client main.go
if [ $? -ne 0 ]; then
	echo "Build failed, exiting"
	exit 1
fi

modprobe sch_netem
deactivate_netem

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


export "MINQ_LOG"=mp,mutex

./client -addr=10.0.4.4:4433
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

