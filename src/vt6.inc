# Call as `vt6_connect "$0" "$@"`. Will re-exec if not yet connected to the VT6 socket.
# KNOWN BUG: Will fail if there are whitespaces in $0 or any of the arguments.
vt6_connect() {
  if [ -z "${VT6:-}" ]; then
    echo "FATAL: no VT6 socket found" >&2
    exit 1
  fi

  # when this variable (with a UUID name in order to not pollute the namespace)
  # is set, it means that the VT6 socket is open on FD 3
  local VARNAME="var638bdbf63295449bb4df74775b01dcb0"
  eval "local VALUE=\${${VARNAME}:-0}"
  if [ "${VALUE}" = 1 ]; then
    unset "${VARNAME}"
    return 0
  fi

  local PROGNAME="$(shift)"
  exec socat EXEC:"env ${VARNAME}=1 ${PROGNAME} $@",fdin=3,fdout=3,sigint UNIX-CONNECT:"$VT6"
}
