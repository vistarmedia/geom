#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <geos_c.h>


void notice(const char *fmt, ...) {
  va_list ap;

  fprintf(stderr, "[geos.notice] ");
  va_start(ap,fmt);
  vfprintf(stderr, fmt, ap);
  va_end(ap);
  fprintf(stderr, "\n");
}

void error(const char *fmt, ...) {
  va_list ap;

  fprintf(stderr, "[geos.error] ");
  va_start(ap,fmt);
  vfprintf(stderr, fmt, ap);
  va_end(ap);
  fprintf(stderr, "\n");
}

GEOSContextHandle_t createGEOSHandle() {
  return initGEOS_r(notice, error);
}
