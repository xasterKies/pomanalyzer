#!/bin/bash

# Define the database file
DB_FILE="pomo.db"

# SQL command to delete all records from the "interval" table
SQL_COMMAND="DELETE FROM \"interval\";"

# Execute the SQL command using sqlite3
sqlite3 "$DB_FILE" "$SQL_COMMAND"

echo "All previous record from the 'interval' table have been deleted."