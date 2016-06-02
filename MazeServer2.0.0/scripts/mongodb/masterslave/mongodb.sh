#master
/home/mongodb/mongodb/bin/mongod --dbpath /home/mongodb/data/db --logpath /home/mongodb/data/log/mdb.log --port 20143 --rest -master
#slave
/home/mongodb/mongodb/bin/mongod --dbpath /home/mongodb/data2/db --logpath /home/mongodb/data2/log/mdb.log --port 20144 --rest -slave -source 54.200.192.17:20145