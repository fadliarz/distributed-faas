#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e
# Print each command to the console before it is executed for debugging.
set -x

# Start a temporary mongod instance in the background without authentication.
mongod --replSet rs0 --bind_ip_all &

# Wait until the MongoDB server is up and ready to accept connections.
until mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
  echo 'Waiting for MongoDB to start...'
  sleep 1
done

# Connect and run the setup script using a 'here document' (<<EOF).
mongosh <<EOF
try {
  // Try to get replica set status. If it fails, the set is not initiated.
  rs.status();
  print('Replica set already initiated.');
} catch (e) {
  // If getting status fails, initiate the replica set with explicit host config.
  print('Initiating replica set...');
  rs.initiate({
    _id: 'rs0',
    members: [
      { _id: 0, host: 'distributed-faas-mongo:27017' }
    ]
  });
}

// Wait for the replica set to have a primary member.
while (!rs.isMaster().ismaster) {
  print('Waiting for primary...');
  sleep(1000);
}

// Now that we have a primary, check if the user needs to be created.
// The shell correctly expands these environment variables before passing them to mongosh.
if (db.getSiblingDB('admin').getUsers().users.length === 0) {
  print('Creating admin user...');
  db.getSiblingDB('admin').createUser({
    user: '$MONGO_INITDB_ROOT_USERNAME',
    pwd: '$MONGO_INITDB_ROOT_PASSWORD',
    roles: [ { role: 'root', db: 'admin' } ]
  });
  print('Admin user created.');
} else {
  print('Admin user already exists.');
}
EOF

# Gracefully shut down the temporary mongod instance.
mongod --shutdown

# Wait for the background mongod process to fully terminate.
wait

echo 'Restarting with authentication...'
# Start the final mongod process in the foreground with authentication enabled.
# Use 'exec' to replace the script process with the mongod process,
# making it the main process of the container.
exec mongod --replSet rs0 --bind_ip_all --auth --keyFile /etc/mongo/auth.key