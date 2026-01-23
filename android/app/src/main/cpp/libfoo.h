#ifndef LIBFOO_H
#define LIBFOO_H

#ifdef __cplusplus
extern "C" {
#endif

char* createKeyPair();
char* getPbKey();
char* getPvKey();
char* createVp(char* did, char* iss, char* pvKeyBase58, char* credTyp, char* vcJwt);
long claimsGetExp(char* vpJwt);
long claimsGetIat(char* vpJwt);
long claimsGetInt64(char* vpJwt, char* claim);

#ifdef __cplusplus
}
#endif

#endif
