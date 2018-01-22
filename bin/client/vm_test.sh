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

./client -addr=10.0.4.4:4433
wait


activate_delay
wait
cat alice.txt | ./client -addr=10.0.4.4:4433 > delay.result
wait
deactivate_netem
./flipper delay.result
wait

activate_loss
wait
cat alice.txt | ./client -addr=10.0.4.4:4433 > loss.result
wait
deactivate_netem
./flipper loss.result
wait

# activate_reordering
# wait
# cat alice.txt | ./client -addr=10.0.4.4:4433 > reordering.result
# wait
# deactivate_netem
# ./flipper reordering.result
# wait

diff alice.txt flipped-delay.result > /dev/null
if [ $? -eq 0 ]; then
	echo "Delay test passed without errors"
elif [$? -eq 1 ]; then
	echo "Delay test failed"
else
	echo "Diff exited with error code"
fi

diff alice.txt flipped-loss.result > /dev/null
if [ $? -eq 0 ]; then
        echo "Loss test passed without errors"
elif [$? -eq 1 ]; then
        echo "Loss test failed"
else
        echo "Diff exited with error code"
fi

# diff alice.txt flipped-reordering.result > /dev/null
# if [ $? -eq 0 ]; then
#         echo "Delay test passed without errors"
# elif [$? -eq 1 ]; then
#         echo "Delay caused rordering"
# else
#         echo "Diff exited with error code"
# fi
