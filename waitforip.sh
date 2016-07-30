#!/bin/bash
while [ 1 ]; do
    ip=`/sbin/ifconfig |grep -v "127.0.0.1" |grep "inet "`
    [[ $? == 0 ]] && break;
    echo "wait for ip"; sleep 1;
done
echo "got ip: $ip"