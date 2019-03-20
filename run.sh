# cleanup
echo "Cleaning up OLD log before starting new cluster"
rm -rf ./log/*

#goreman
goreman -f ./Procfile start
# goreman -f ./Procfile2 start
# goreman -f ./Procfile3 start