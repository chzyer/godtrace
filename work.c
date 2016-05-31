#include "dtrace.h"
#include "_cgo_export.h"

int chew(dtrace_probedata_t *data, void *arg) {
  return goChew(data, arg);
}

int chewrec(dtrace_probedata_t *data, dtrace_recdesc_t *rec, void *arg) {
  return goChewRec(data, rec, arg);
}

int dumpChewrec(dtrace_probedata_t *data, dtrace_recdesc_t *rec, void *arg) {
  return (DTRACE_CONSUME_THIS);
}


