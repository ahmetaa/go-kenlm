#ifndef GO_KENLM_WRAP_H
#define GO_KENLM_WRAP_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct kenlm_model kenlm_model;

/* Load a KenLM model (ARPA or .binary). Returns NULL on failure and writes
 * an error message into *err_out (caller must free with kenlm_string_free). */
kenlm_model* kenlm_load(const char* path, char** err_out);

void kenlm_free(kenlm_model* m);

/* Order of the loaded model. */
unsigned int kenlm_order(const kenlm_model* m);

/* Score a sequence of whitespace-separated tokens. Adds <s> at the start
 * and </s> at the end. Returns total log10 probability. */
double kenlm_score_sentence(const kenlm_model* m, const char* tokens);

/* Free a string returned by this library. */
void kenlm_string_free(char* s);

#ifdef __cplusplus
}
#endif

#endif
