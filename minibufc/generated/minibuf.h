#ifndef MINIBUF_H
#define MINIBUF_H

#include <stdbool.h>
#include <stdint.h>
#include <stddef.h>

#define MB_OK 0
#define MB_ERR_INVALID_FORMAT 1
#define MB_ERR_BUFFER_TOO_SMALL 2

extern int mb_float_precision;

typedef struct {
    float x;
    float y;
    float z;
} vector_t;

typedef struct {
    bool auto_restart;
    int32_t id;
    char user_name[256];
    float score;
} config_t;

int mb_vector_parse(const char* buf, vector_t* out);
int mb_vector_serialize(const vector_t* in, char* buf, size_t buf_size);

int mb_config_parse(const char* buf, config_t* out);
int mb_config_serialize(const config_t* in, char* buf, size_t buf_size);

#endif
