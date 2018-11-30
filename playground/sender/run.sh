COUNT=0
while :
do
    let "COUNT++"
    curl "http://mock-server:5000/sink/$COUNT"
    sleep 1
done
