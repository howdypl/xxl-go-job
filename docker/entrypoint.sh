#!/bin/bash
set -e
if [[ ! -f /opt/jd-account/config/settings.yml ]]
then
    cp /opt/jd-account/default_config/* /opt/jd-account/config/
fi

/opt/jd-account/jd-account server -c=/opt/jd-account/config/settings.yml
