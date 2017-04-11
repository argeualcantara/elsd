#!/bin/bash
#
# (C) Copyright 2016 HP Development Company, L.P.
# Confidential computer software. Valid license from HP required for possession, use or copying.
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items are licensed
# to the U.S. Government under vendor's standard commercial license.
#
echo "=> Loading container environment variables..."
declare -x ELS_ADDRESS=${ELS_ADDRESS}
declare -x ELS_PORT=${ELS_PORT}
declare -x ELS_DEBUG=${ELS_DEBUG}


echo "=> Provisioning Dynamodb Table..."

echo "=> Starting ELS..."
/els
