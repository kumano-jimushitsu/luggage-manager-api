#start SQL Server, start the script to create/setup the DB
 #wait for the SQL Server to come up
sleep 15s

echo "running set up script"
#run the setup script to create the DB and the schema in the DB
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P My-Secret-Pwd-1 -d master -i db-init.sql