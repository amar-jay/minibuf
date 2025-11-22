#include "./minibuf.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

int main() {
		// Example usage of the generated functions
		vector_t vec;
		config_t cfg;

		vec.x = 1.234;
		vec.y = 5.678;
		vec.z = 9.101;

		cfg.auto_restart = 1;
		cfg.id = 42;	
		strcpy(cfg.user_name, "Ted Balkjfa");
		cfg.score = 100.0;

		char vectorBuf[256];
		char configBuf[256];


		mb_vector_serialize(&vec, vectorBuf, sizeof(vectorBuf));
		mb_config_serialize(&cfg, configBuf, sizeof(configBuf));

		printf("Serialized Vector: %s\n", vectorBuf);
		printf("Serialized Config: %s\n", configBuf);

		 vector_t parsedVec;
		 config_t parsedCfg;
		mb_vector_parse(vectorBuf, &parsedVec);
		mb_config_parse(configBuf, &parsedCfg);
		printf("Parsed Vector: x=%f, y=%f, z=%f\n", parsedVec.x, parsedVec.y, parsedVec.z);
		printf("Parsed Config: auto_restart=%d, id=%d, user_name=%s, score=%f\n",
		       parsedCfg.auto_restart, parsedCfg.id, parsedCfg.user_name, parsedCfg.score);
		return 0;
}