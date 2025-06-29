#!/bin/bash
set -e
# wait for  module database
wait-for-it.sh $MODULE_USER_DB_HOST:$MODULE_USER_DB_PORT --strict --timeout=5 -- echo "User Module Database is up"
wait-for-it.sh $MODULE_BOOK_DB_HOST:$MODULE_BOOK_DB_PORT --strict --timeout=5 -- echo "Book Module Database is up"
# wait for service discovery
#wait-for-it.sh $SERVICE_DISCOVERY_IP:$SERVICE_DISCOVERY_PORT --strict --timeout=5 -- echo "Service Discovery is up"
# wait for service bus
#wait-for-it.sh $SERVICE_BUS_IP:$SERVICE_BUS_PORT --strict --timeout=5 -- echo "Service Bus is up"
exec "$@"