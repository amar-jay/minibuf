#include "minibuf.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

int mb_float_precision = 3;

int mb_vector_parse(const char* buf, vector_t* out) {
    char* start = strchr(buf, '[');
    if (!start) return MB_ERR_INVALID_FORMAT;
    start++;
    char* end = strchr(start, ']');
    if (!end) return MB_ERR_INVALID_FORMAT;
    *end = '\0';
    int count = atoi(start);
    *end = ']';
    char* values = end + 1;
    char* vals = strdup(values);
    char* token = strtok(vals, ";");
    int i = 0;
    while (token && i < 3) {
        if (i == 0) {
            out->x = atof(token);
        }
        if (i == 1) {
            out->y = atof(token);
        }
        if (i == 2) {
            out->z = atof(token);
        }
        i++;
        token = strtok(NULL, ";");
    }
    free(vals);
    return MB_OK;
}

int mb_vector_serialize(const vector_t* in, char* buf, size_t buf_size) {
    int len = snprintf(buf, buf_size, "[3]", 3);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    char* pos = buf + len;
    len += snprintf(pos, buf_size - len, "%.*f", mb_float_precision, in->x);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%.*f", mb_float_precision, in->y);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%.*f", mb_float_precision, in->z);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    return MB_OK;
}

int mb_config_parse(const char* buf, config_t* out) {
    char* start = strchr(buf, '[');
    if (!start) return MB_ERR_INVALID_FORMAT;
    start++;
    char* end = strchr(start, ']');
    if (!end) return MB_ERR_INVALID_FORMAT;
    *end = '\0';
    int count = atoi(start);
    *end = ']';
    char* values = end + 1;
    char* vals = strdup(values);
    char* token = strtok(vals, ";");
    int i = 0;
    out->score = 6000.000;
    while (token && i < 4) {
        if (i == 0) {
            out->auto_restart = strcmp(token, "T") == 0;
        }
        if (i == 1) {
            out->id = atoi(token);
        }
        if (i == 2) {
            strcpy(out->user_name, token);
        }
        if (i == 3) {
            out->score = atof(token);
        }
        i++;
        token = strtok(NULL, ";");
    }
    free(vals);
    return MB_OK;
}

int mb_config_serialize(const config_t* in, char* buf, size_t buf_size) {
    int len = snprintf(buf, buf_size, "[4]", 4);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    char* pos = buf + len;
    len += snprintf(pos, buf_size - len, "%s", in->auto_restart ? "T" : "F");
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%d", in->id);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%s", in->user_name);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%.*f", mb_float_precision, in->score);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    return MB_OK;
}

int mb_govt_parse(const char* buf, govt_t* out) {
    char* start = strchr(buf, '[');
    if (!start) return MB_ERR_INVALID_FORMAT;
    start++;
    char* end = strchr(start, ']');
    if (!end) return MB_ERR_INVALID_FORMAT;
    *end = '\0';
    int count = atoi(start);
    *end = ']';
    char* values = end + 1;
    char* vals = strdup(values);
    char* token = strtok(vals, ";");
    int i = 0;
    while (token && i < 5) {
        if (i == 0) {
            out->minister_count = atoi(token);
        }
        if (i == 1) {
            strcpy(out->president_name, token);
        }
        if (i == 2) {
            out->president_term = atoi(token);
        }
        if (i == 3) {
            strcpy(out->prime_minister_name, token);
        }
        if (i == 4) {
            out->vote_count = atof(token);
        }
        i++;
        token = strtok(NULL, ";");
    }
    free(vals);
    return MB_OK;
}

int mb_govt_serialize(const govt_t* in, char* buf, size_t buf_size) {
    int len = snprintf(buf, buf_size, "[5]", 5);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    char* pos = buf + len;
    len += snprintf(pos, buf_size - len, "%d", in->minister_count);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%s", in->president_name);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%d", in->president_term);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%s", in->prime_minister_name);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    len += snprintf(pos, buf_size - len, ";%.*f", mb_float_precision, in->vote_count);
    if (len >= buf_size) return MB_ERR_BUFFER_TOO_SMALL;
    pos = buf + len;
    return MB_OK;
}

