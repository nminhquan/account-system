# cleanup
echo "Cleaning up OLD log before starting new cluster"
rm -rf ./log/*
rm -rf ./test/test_db/*

#goreman
goreman -f ./Procfile2 start > ./log/tc.log &
P1=$!
sleep 1
goreman -f ./Procfile3 start > ./log/rm.log &
P2=$!
wait $P1 $P2
# goreman -f ./Procfile2 start
# goreman -f ./Procfile3 start