FROM jchorl/appengine

ADD . src/github.com/jchorl/financejc
WORKDIR src/github.com/jchorl/financejc
ENTRYPOINT dev_appserver.py --host=0.0.0.0 --admin_host=0.0.0.0 --skip_sdk_update_check=yes appengine