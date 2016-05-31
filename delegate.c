#include "dtrace.h"
#include "_cgo_export.h"
#include "stdlib.h"

int chew(dtrace_probedata_t *data, void *arg) {
  return goChew(data, arg);
}

int chewrec(dtrace_probedata_t *data, dtrace_recdesc_t *rec, void *arg) {
  return goChewRec(data, rec, arg);
}

int dumpChewrec(dtrace_probedata_t *data, dtrace_recdesc_t *rec, void *arg) {
  return (DTRACE_CONSUME_THIS);
}

int bufhandler(dtrace_bufdata_t *bufdata, void *arg) {
  return goBufHandler(bufdata, arg);
}

void freeArray(int n, char** chs) {
  for (int i=0; i<n; i++) {
	free(chs[i]);
  }
  free(chs);
}

void freeString(char *ch) {
  free(ch);
}
