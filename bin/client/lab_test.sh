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
./client -addr=10.0.1.10:4433
wait
deactivate_netem
wait

activate_delay
wait
cat alice.txt | ./client -addr=10.0.1.10:4433 > delay.result
wait
deactivate_netem
./flipper delay.result
wait

activate_loss
cat alice.txt | ./client -addr=10.0.1.10:4433 > loss.result
wait
deactivate_netem
./flipper loss.result
wait

# activate_reordering
# wait
# cat alice.txt | ./client -addr=10.0.1.10:4433 > reordering.result
# wait
# deactivate_netem
# ./flipper reordering.result
# wait

# Result Delay
diff alice.txt flipped-delay.result > /dev/null
if [ $? -eq 0 ]; then
	echo "Delay test passed without errors"
elif [ $? -eq 1 ]; then
	echo "Delay test failed"
else
	echo "Diff exited with error code"
fi

# Result Loss
diff alice.txt flipped-loss.result > /dev/null
if [ $? -eq 0 ]; then
        echo "Loss test passed without errors"
elif [ $? -eq 1 ]; then
        echo "Loss test failed"
else
        echo "Diff exited with error code"
fi

rm delay.result
rm loss.result
rm client

# Result Reordering
# diff alice.txt flipped-reordering.result > /dev/null
# if [ $? -eq 0 ]; then
#         echo "Delay test passed without errors"
# elif [ $? -eq 1 ]; then
#         echo "Reordering test failed"
# else
#         echo "Diff exited with error code"
# fi
