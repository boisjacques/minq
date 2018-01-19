#!/bin/bash

activate_loss () {
	tc qdisc change dev enp3s0f0 root netem loss 5% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp3s0f0"
		exit 1
	fi
	tc qdisc change dev enp3s0f1 root netem loss 2% 25%
	if [ $? -ne 0 ]; then
		echo "Adding loss failed on enp3s0f1"
		exit 1
	fi
	echo "Added loss on all interfaces"
}

activate_reordering () {
	tc qdisc change dev enp3s0f0 root netem gap 5 delay 10ms
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp3s0f0"
		exit 1
	fi
	tc qdisc change dev enp3s0f1 root netem gap 2 delay 45ms
	if [ $? -ne 0 ]; then
		echo "Adding reordering failed on enp3s0f1"
		exit 1
	fi
	echo "Added reordering on all interfaces"

}

activate_delay () {
	tc qdisc add dev enp3s0f0 root netem delay 100ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp3s0f0"
		exit 1
	fi
	tc qdisc add dev enp3s0f1 root netem delay 150ms 10ms 25%
	if [ $? -ne 0 ]; then
		echo "Adding delay failed on enp3s0f1"
		exit 1
	fi
	echo "Activated delay on all interfaces"
}

deactivate_netem() {
	tc qdisc del dev enp3s0f0 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem failed on enp0s10"
		exit 1
	fi
	tc qdisc del dev enp3s0f1 root
	if [ $? -ne 0 ]; then
		echo "Deactivating netem failed on enp0s10"
		exit 1
	fi
	echo "Deactivated netem on all interfaces"
}

if [ "$(whoami)" != "root" ]; then
	echo "Sorry, you are not root."
	exit 1
fi

go build -o client main.go
modprobe sch_netem

if [ $? -ne 0 ]; then
	echo "Build failed, exiting"
	exit 1
fi
./client -addr=10.0.4.4:4433
wait
activate_delay
wait
cat alice.txt | ./client -addr=10.0.4.4:4433 > delay.result
wait
deactivate_netem
./flipper delay.result
wait
diff alice.txt flipped-delay.result > /dev/null
wait
if [ $? -eq 0 ]; then
	echo "Delay test passed without errors"
elif [$? -eq 1 ]; then
	echo "Delay test failed"
else
	echo "Diff exited with error code"
fi
activate_loss
cat alice.txt | ./client -addr=10.0.4.4:4433 > loss.result
deactivate_netem
./flipper loss.result
wait
diff alice.txt flipped-loss.result > /dev/null
wait
if [ $? -eq 0 ]; then
        echo "Loss test passed without errors"
elif [$? -eq 1 ]; then
        echo "Loss test failed"
else
        echo "Diff exited with error code"
fi

activate_reordering
wait
cat alice.txt | ./client -addr=10.0.4.4:4433 > reordering.result
wait
deactivate_netem
./flipper reordering.result
wait
diff alice.txt flipped-reordering.result > /dev/null
if [ $? -eq 0 ]; then
        echo "Delay test passed without errors"
elif [$? -eq 1 ]; then
        echo "Delay caused rordering"
else
        echo "Diff exited with error code"
fi
