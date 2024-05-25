#!/bin/sh

( rabbitmqctl wait --timeout 60 $RABBITMQ_PID_FILE ; \
rabbitmqctl add_user $RABBIT_MQ_USER $RABBIT_MQ_PASSWORD 2>/dev/null ; \
rabbitmqctl set_user_tags $RABBIT_MQ_USER administrator ; \
rabbitmqctl set_permissions -p / $RABBIT_MQ_USER  ".*" ".*" ".*" ; \


rabbitmqctl add_user $RABBITMQ_DEFAULT_USER $RABBITMQ_DEFAULT_PASS 2>/dev/null ; \
rabbitmqctl set_permissions -p / $RABBITMQ_DEFAULT_USER  ".*" ".*" ".*" ; \
rabbitmqctl set_user_tags $RABBITMQ_DEFAULT_USER administrator ; ) &

# $@ is used to pass arguments to the rabbitmq-server command.
# For example if you use it like this: docker run -d rabbitmq arg1 arg2,
# it will be as you run in the container rabbitmq-server arg1 arg2
rabbitmq-server $@